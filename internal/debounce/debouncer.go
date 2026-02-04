package debounce

import (
	"sync"
	"time"


	"github.com/Atharv-3105/File-Manager/internal/model"
)

type Debouncer struct {
	delay 	time.Duration
	input	<-chan model.FileEvent
	output	chan<- model.FileEvent

	mu 	sync.Mutex
	timers 	map[string]*time.Timer
}


func New(
	delay time.Duration,
	input   <-chan model.FileEvent,
	output  chan<- model.FileEvent,
) *Debouncer {
	return &Debouncer{
		delay: delay,
		input: input,
		output: output,
		timers: make(map[string]*time.Timer),
	}
}


func (d *Debouncer) Start() {
	go func() {
		for event := range d.input {
			d.handle(event)
		}
	}()
}


func (d *Debouncer) handle(event model.FileEvent) {
	d.mu.Lock()
	defer d.mu.Unlock()

	//If a timer already exists for curr file; Reset it
	if timer,exists := d.timers[event.Path]; exists {
		timer.Stop()
	}

	//Create a New timer
	d.timers[event.Path] = time.AfterFunc(d.delay, func() {
		d.output <- event

		//Cleanup the curr timer from the timers map
		d.mu.Lock()
		delete(d.timers, event.Path)
		d.mu.Unlock()
	})
}