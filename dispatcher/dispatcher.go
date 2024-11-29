package dispatcher

import (
	"reflect"
	"sync"

	"github.com/callevo/ari/arievent"
	"github.com/panjf2000/ants/v2"
)

type listenersCollection []Listener

type Dispatcher interface {
	Dispatch(e arievent.StasisEvent) arievent.StasisEvent

	AddListener(e arievent.EventType, l Listener)

	RemoveListener(e arievent.EventType, l Listener)

	RemoveAllListener(e arievent.EventType)

	HasListeners(e arievent.EventType) bool
}

func getEventName(e arievent.EventType) string {
	eventType := reflect.TypeOf(e)

	// Convert the type to a string
	return eventType.Name()
}

type EventDispatcher struct {
	sync.RWMutex
	listeners   map[arievent.EventType]listenersCollection
	workersPool *ants.Pool
}

func NewDispatcher() *EventDispatcher {
	pool, err := ants.NewPool(1000)

	if err != nil {
		return nil
	}

	d := &EventDispatcher{
		listeners:   make(map[arievent.EventType]listenersCollection),
		workersPool: pool,
	}

	return d
}

func (d *EventDispatcher) GetPool() *ants.Pool {
	return d.workersPool
}

func (d *EventDispatcher) AddListener(e arievent.EventType, l Listener) {
	d.RWMutex.Lock()
	defer d.RWMutex.Unlock()

	d.listeners[e] = append(d.listeners[e], l)
}

func (d *EventDispatcher) ExecuteOnce(e *arievent.StasisEvent, l Listener) Listener {
	var newListener func(e *arievent.StasisEvent)
	newListener = func(e *arievent.StasisEvent) {
		l(e)
		d.RWMutex.RUnlock() // The dispatcher is locked in the Dispatch method, need to unlock it
		d.RemoveListener(e.GetType(), newListener)
		d.RWMutex.RLock()
	}

	return newListener
}

func (d *EventDispatcher) RemoveListener(e arievent.EventType, l Listener) {
	d.RWMutex.Lock()
	defer d.RWMutex.Unlock()

	p := reflect.ValueOf(l).Pointer()

	listeners := d.listeners[e]
	for i, l := range listeners {
		lp := reflect.ValueOf(l).Pointer()
		if lp == p {
			d.listeners[e] = append(listeners[:i], listeners[i+1:]...)
		}
	}
}

func (d *EventDispatcher) RemoveAll(e arievent.EventType) {
	d.RWMutex.Lock()
	defer d.RWMutex.Unlock()

	_, ok := d.listeners[e]
	if ok != false {
		delete(d.listeners, e)
	}
}

func (d *EventDispatcher) HasListeners(e arievent.EventType) bool {

	listeners, ok := d.listeners[e]
	if ok == false {
		return false
	}

	return len(listeners) != 0
}

func (d *EventDispatcher) Dispatch(e *arievent.StasisEvent) *arievent.StasisEvent {
	d.RWMutex.RLock()
	defer d.RWMutex.RUnlock()

	for _, lst := range d.listeners[e.GetType()] {
		d.workersPool.Submit(func() {
			lst(e)
		})
	}

	return e
}
