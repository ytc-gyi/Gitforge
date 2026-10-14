package core

import "sync"

// EventType identifies the kind of event.
type EventType string

const (
	EventRepoChanged    EventType = "repo.changed"
	EventStatusChanged  EventType = "status.changed"
	EventFetchComplete  EventType = "fetch.complete"
	EventCommitCreated  EventType = "commit.created"
	EventBranchChanged  EventType = "branch.changed"
	EventConfigChanged  EventType = "config.changed"
	EventPluginLoaded   EventType = "plugin.loaded"
	EventPluginUnloaded EventType = "plugin.unloaded"
	EventAIToken        EventType = "ai.token"
	EventAIComplete     EventType = "ai.complete"
)

// Event represents something that happened in the system.
type Event struct {
	Type    EventType
	Payload any
}

// EventHandler handles an event.
type EventHandler func(Event)

// EventBus provides publish/subscribe for system events.
type EventBus struct {
	mu       sync.RWMutex
	handlers map[EventType][]eventSub
	nextID   int
}

type eventSub struct {
	id      int
	handler EventHandler
}

// NewEventBus creates a new EventBus.
func NewEventBus() *EventBus {
	return &EventBus{
		handlers: make(map[EventType][]eventSub),
	}
}

// Subscribe registers a handler for the given event type and returns an unsubscribe ID.
func (b *EventBus) Subscribe(eventType EventType, handler EventHandler) int {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.nextID++
	id := b.nextID
	b.handlers[eventType] = append(b.handlers[eventType], eventSub{id: id, handler: handler})
	return id
}

// Unsubscribe removes a handler by its subscription ID.
func (b *EventBus) Unsubscribe(id int) {
	b.mu.Lock()
	defer b.mu.Unlock()
	for eventType, subs := range b.handlers {
		for i, sub := range subs {
			if sub.id == id {
				b.handlers[eventType] = append(subs[:i], subs[i+1:]...)
				return
			}
		}
	}
}

// Publish sends an event to all subscribers of that event type.
func (b *EventBus) Publish(event Event) {
	b.mu.RLock()
	subs := make([]eventSub, len(b.handlers[event.Type]))
	copy(subs, b.handlers[event.Type])
	b.mu.RUnlock()

	for _, sub := range subs {
		sub.handler(event)
	}
}
