package watcher

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/Atharv-3105/File-Manager/internal/model"
)

func StartPolling(dir string, out chan<- model.FileEvent) {

	ticker := time.NewTicker(2 * time.Second)

	lastSeen := make(map[string]time.Time)

	go func() {

		for range ticker.C {
			filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {

				if err != nil || info.IsDir() {
					return nil
				}

				modTime := info.ModTime()

				prev, exists := lastSeen[path]

				if !exists || modTime.After(prev) {

					lastSeen[path] = modTime 

					event := model.FileEvent {
						Path: path,
						EventType: "POLL",
					}

					log.Println("[POLLER] detected changes: ", path)

					out <- event
				}

				return nil
			})
		}
	}()
}