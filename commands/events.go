package commands

import (
	"strconv"
	"sync"
	"time"

	"github.com/phonkee/patrol/context"
	"github.com/phonkee/patrol/models"

	"github.com/golang/glog"
	"github.com/phonkee/patrol/core"
	"github.com/phonkee/patrol/parser"
	"github.com/phonkee/patrol/settings"
)

/*
Event commands
*/

func NewEventWorkerCommand(context *context.Context, onevent func(*models.Event, *models.EventGroup)) core.Commander {
	return &EventWorkerCommand{
		context: context,
		onevent: onevent,
	}
}

type EventWorkerCommand struct {
	core.Command
	context *context.Context

	// cli args settings
	goroutinesCount int

	// signal handler on event
	onevent func(*models.Event, *models.EventGroup)
}

func (ew *EventWorkerCommand) ID() string { return "worker" }
func (ew *EventWorkerCommand) Description() string {
	return `Runs event process worker
patrol event:worker [goroutines=100]`
}
func (ew *EventWorkerCommand) Run() error {
	glog.Infof("event worker: running %d workers (goroutines).", ew.goroutinesCount)
	var wg sync.WaitGroup

	wg.Add(ew.goroutinesCount)
	for i := 0; i < ew.goroutinesCount; i++ {
		worker := NewEventWorker(ew.context, i, ew.onevent)
		go func() {
			defer wg.Done()
			worker.Run()
		}()
	}

	wg.Wait()

	return nil
}

func (e *EventWorkerCommand) ParseArgs(args []string) (err error) {
	// Parse goroutines count
	e.goroutinesCount = settings.EVENT_WORKER_DEFAULT_GOROUTINES_COUNT
	if len(args) > 0 {
		e.goroutinesCount, err = strconv.Atoi(args[0])
		if err != nil {
			return
		} else if e.goroutinesCount <= 0 {
			e.goroutinesCount = settings.EVENT_WORKER_DEFAULT_GOROUTINES_COUNT
		}
	}
	return nil
}

/*
Event background worker
*/

func NewEventWorker(context *context.Context, id int, onevent func(*models.Event, *models.EventGroup)) *EventWorker {
	return &EventWorker{context: context, id: id, onevent: onevent}
}

type EventWorker struct {
	context *context.Context
	id      int

	// signal handler
	onevent func(*models.Event, *models.EventGroup)
}

/* runs event worker
 */
func (e *EventWorker) Run() error {
	processing := false

	remanager := parser.NewRawEventManager(e.context)

	for {
		select {
		case _ = <-e.context.Quit:
			break
		case <-time.After(time.Second):
			if processing == false {
				processing = true
				for {
					re, message, err := remanager.PopRawEvent()
					if err != nil {
						glog.V(2).Infof("event worker-%d: %s.", e.id, err)
						break
					}

					if err := e.ProcessEvent(re); err != nil {
						glog.Errorf("event worker-%d: process message failed: %+v", e.id, re)
						continue
					}

					if err := message.Ack(); err != nil {
						glog.Error("message ack failed with %s.", err)
					}
				}
				processing = false
			}
		}
	}
	return nil
}

/*
	Process event
*/
func (e *EventWorker) ProcessEvent(re *parser.RawEvent) (err error) {
	eventManager := models.NewEventManager(e.context)
	eventgroupManager := models.NewEventGroupManager(e.context)

	var (
		event      *models.Event
		eventgroup *models.EventGroup
	)

	// some error occured
	if event, eventgroup, err = eventManager.NewEventFromRaw(re); err != nil {
		return
	}

	// increment counters
	err = eventgroupManager.IncrementCounters(eventgroup)
	if err != nil {
		glog.Errorf("Increment counters returned error %+v", err)
	}

	// send signal
	e.onevent(event, eventgroup)

	return
}
