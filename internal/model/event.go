package model

import "time"

type EventType string 

const (
	EventCreate EventType = "CREATE"
	EventWrite EventType = "WRITE"
	EventRemove EventType = "REMOVE"
	EventRename EventType = "RENAME"
)

type FileEvent struct {
	Path     string
	EventType EventType
	Timestamp time.Time 

}