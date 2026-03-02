package main

import (
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/Atharv-3105/File-Manager/internal/client"
	"github.com/Atharv-3105/File-Manager/internal/model"
	"github.com/Atharv-3105/File-Manager/internal/storage"
	"github.com/Atharv-3105/File-Manager/internal/watcher"
	"github.com/Atharv-3105/File-Manager/internal/worker"

	// "github.com/fsnotify/fsnotify"
	"github.com/Atharv-3105/File-Manager/internal/debounce"
)


const (
	WorkerCount = 3
	QueueSize = 100
	DebounceDelay = 500 * time.Millisecond
)

func main() {
	log.Println("[WATCHER] starting ingestion pipeline")

	//Initialize DataBase
	_ = os.MkdirAll("data", os.ModePerm)

	db, err := storage.Open("data/index.db")
	if err != nil {
		log.Fatal("[WATCHER] db open failure: ",err)
	}
	defer db.Close()

	absDB, _ := filepath.Abs("data/index.db")
	log.Println("[WATCHER] DB absoulte path:", absDB)
	
	if err := storage.InitSchema(db); err != nil {
		log.Fatal("[WATCHER] schema init:", err)
	}

	//Set the WATCH Directory
	watchDir := "./watched"

	//Check Dir Exists or Not
	if err := os.MkdirAll(watchDir, os.ModePerm); err != nil {
		log.Fatalf("[WATCHER] failed to create watch directory: %v", err)
	}

	absPath, _ := filepath.Abs(watchDir)
	log.Printf("[WATCHER] watching directory: %s\n", absPath)

	// jobQueue := make(chan model.FileEvent, QueueSize)
	rawEvents := make(chan model.FileEvent, QueueSize)
	stableEvents := make(chan model.FileEvent, QueueSize)

	fsWatcher, err := watcher.New(absPath, rawEvents)
	if err != nil {
		log.Fatal("[WATCHER] fs watcher:", err)
	}
	defer fsWatcher.Close()

	fsWatcher.Start()

	//Add Debouncer
	debouncer := debounce.New(DebounceDelay, rawEvents, stableEvents)
	debouncer.Start()

	extractor := client.NewExtractorClient("http://127.0.0.1:8001")
	pool := worker.New(WorkerCount, stableEvents, db, extractor)

	//Handle shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	log.Println("[WATCHER] shutdown initiated")

	close(rawEvents)
	close(stableEvents)
	pool.Wait()

	log.Println("[WATCHER] graceful shutdown complete")

}
