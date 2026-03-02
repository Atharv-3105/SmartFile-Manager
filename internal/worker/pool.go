package worker

import (
	"database/sql"
	"log"
	"sync"

	"github.com/Atharv-3105/File-Manager/internal/client"
	"github.com/Atharv-3105/File-Manager/internal/model"
	"github.com/Atharv-3105/File-Manager/internal/storage"
)


type Pool struct {
	wg sync.WaitGroup
	db   *sql.DB
	extractor   *client.ExtractorClient
}

func New(workerCount int, 
		 jobs <- chan model.FileEvent,
		 db    *sql.DB,
		 extractor   *client.ExtractorClient) *Pool {

	p := &Pool{
		db: db,
		extractor: extractor,
	}

	for i := 1; i <= workerCount; i++{
		p.wg.Add(1)
		go p.worker(i, jobs)
	}

	return p
}


func (p *Pool) Wait() {
	p.wg.Wait()
}


func (p *Pool) worker(id int, jobs <- chan model.FileEvent) {
	defer p.wg.Done()


	for event := range jobs {

		log.Printf("[WORKER %d] Processing %s -> %s\n", id,event.EventType, event.Path)

		//Call Python Extraction service
		text, embedding, modelName, err := p.extractor.Extract(event.Path)
		if err != nil {
			log.Println("[WORKER] extraction failed:", err)
			continue	
		}

		//Begin Transaction
		tx, err := p.db.Begin()
		if err != nil {
			log.Println("[WORKER] tx execution failed: ",err)
			continue
		}

		//Upsert File
		fileID, err := storage.UpsertFile(tx, event.Path)
		if err != nil {
			log.Println("[WORKER] upsert failed:",err)
			tx.Rollback()
			continue
		}

		//Insert Extraction
		err = storage.InsertExtraction(tx, fileID, text, "ok")
		if err != nil {
			log.Println("[WORKER] insert extraction failed: ", err)
			tx.Rollback()
			continue
		}


		//Insert Extraction
		err = storage.InsertEmbedding(tx, fileID, modelName, embedding)
		if err != nil {
			log.Println("[WORKER] insert embedding failed:",err)
			tx.Rollback()
			continue
		}

		if err := tx.Commit(); err != nil {
			log.Println("[WORKER] commit failed:", err)
			continue
		}

		log.Println("[WORKER] successfully indexed:", event.Path)
	}
}