package main

import (
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/Atharv-3105/File-Manager/internal/model"
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
	watchDir := "./watched"

	//Check Dir Exists or Not
	if err := os.MkdirAll(watchDir, os.ModePerm); err != nil {
		log.Fatalf("[MAIN] failed to create watch directory: %v", err)
	}

	absPath, _ := filepath.Abs(watchDir)
	log.Printf("[MAIN] Watching Directory: %s\n", absPath)

	// jobQueue := make(chan model.FileEvent, QueueSize)
	rawEvents := make(chan model.FileEvent, QueueSize)
	stableEvents := make(chan model.FileEvent, QueueSize)

	fsWatcher, err := watcher.New(absPath, rawEvents)
	if err != nil {
		log.Fatal(err)
	}
	defer fsWatcher.Close()

	fsWatcher.Start()

	//Add Debouncer
	debouncer := debounce.New(DebounceDelay, rawEvents, stableEvents)
	debouncer.Start()

	pool := worker.New(WorkerCount, stableEvents)

	//Handle shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	log.Println("[MAIN] Shutdown Initiated")

	close(rawEvents)
	close(stableEvents)
	pool.Wait()

	log.Println("[MAIN] Graceful shutdown complete")

}
