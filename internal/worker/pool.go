package worker

import (
	"log"
	"sync"
	"time"

	"github.com/Atharv-3105/File-Manager/internal/model"
)

type Pool struct {
	wg sync.WaitGroup
}

func New(workerCount int, jobs <- chan model.FileEvent) *Pool {
	p := &Pool{}

	for i := 1; i <= workerCount; i++{
		p.wg.Add(1)
		go worker(i, jobs, &p.wg)
	}

	return p
}


func (p *Pool) Wait() {
	p.wg.Wait()
}


func worker(id int, jobs <- chan model.FileEvent, wg *sync.WaitGroup) {
	defer wg.Done()

	for event := range jobs {
		log.Printf("[WORKER %d] Processing %s -> %s\n", id,event.EventType, event.Path)

		//Simulate the work by adding delay(sleep)
		time.Sleep(500 * time.Millisecond)
	}
}