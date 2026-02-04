package watcher

import (
	"log"
	"time"

	"github.com/Atharv-3105/File-Manager/internal/model"
	"github.com/fsnotify/fsnotify"
)

type FSWatcher struct {
	watcher  *fsnotify.Watcher
	out 	chan<- model.FileEvent
}

func New(path string, out chan<- model.FileEvent) (*FSWatcher, error) {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err 
	}


	if err := w.Add(path); err != nil{
		return nil, err
	}

	return &FSWatcher{
		watcher: w,
		out: out,
	}, nil
}


func (f *FSWatcher) Start() {
	go func() {
		for {
			select {
			case event := <-f.watcher.Events:
				if fe,ok := parseEvent(event); ok {
					f.out <- fe
				}
			
			case err := <-f.watcher.Errors:
				log.Println("[WATCHER] watcher error:", err)
			}
		}
	}()
}

func (f *FSWatcher) Close() error {
	return f.watcher.Close()
}

func parseEvent(event fsnotify.Event) (model.FileEvent, bool) {
	var etype model.EventType

	switch {
	case event.Op&fsnotify.Create == fsnotify.Create:
			etype = model.EventCreate
	case event.Op&fsnotify.Write == fsnotify.Write:
			etype = model.EventWrite
	case event.Op&fsnotify.Remove == fsnotify.Remove:
			etype = model.EventRemove
	case event.Op&fsnotify.Rename == fsnotify.Rename:
			etype = model.EventRename
	default:
			return model.FileEvent{}, false
	}

	return model.FileEvent{
		Path: event.Name,
		EventType: etype,
		Timestamp: time.Now(),
	}, true
}