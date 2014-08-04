package eventbus

import (
	"time"

	"github.com/citadel/citadel"
	"github.com/samalba/dockerclient"
)

type EventBus struct {
	engines  map[string]*citadel.Engine
	handlers map[string][]citadel.EventHandler
}

func New(engines ...*citadel.Engine) (*EventBus, error) {
	bus := &EventBus{
		engines:  make(map[string]*citadel.Engine),
		handlers: make(map[string][]citadel.EventHandler),
	}

	for _, e := range engines {
		bus.engines[e.ID] = e
	}

	return bus, nil
}

func (b *EventBus) AddHandler(eventType string, h citadel.EventHandler) error {
	b.handlers[eventType] = append(b.handlers[eventType], h)

	return nil
}

func (b *EventBus) handler(e *dockerclient.Event, args ...interface{}) {
	engine := args[0].(*citadel.Engine)

	event := &citadel.Event{
		Engine: engine,
		Type:   e.Status,
		Time:   time.Unix(int64(e.Time), 0),
	}

	container, err := citadel.FromDockerContainer(e.Id, e.From, engine)
	if err != nil {
		// TODO: un fuck this shit, fuckin handler
		return
	}

	event.Container = container

	for _, h := range b.handlers[event.Type] {
		h.Handle(event)
	}
}
