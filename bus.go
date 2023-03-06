package eventbus

import (
	"fmt"
	"sync"
)

var (
	ErrInvalidChannelType   = fmt.Errorf("invalid channel type")
	ErrChannelAlreadyExists = fmt.Errorf("channel already exists")
)

type Bus interface {
	untypedChannelStore
}

type untypedChannelStore interface {
	Get(string) (any, bool)
	Set(string, any) error
	// Delete(string)
}

type simpleBus struct {
	mu       sync.RWMutex
	channels map[string]any // untyped map for typed Channels

}

// NewBus returns a new event bus
func NewBus() Bus {
	return &simpleBus{
		channels: make(map[string]any),
	}
}

// Get a channel of the event bus
// returns true if found, false otherwise
func (b *simpleBus) Get(evt string) (any, bool) {
	b.mu.RLock() // read lock
	defer b.mu.RUnlock()
	ch, exists := b.channels[evt]
	return ch, exists
}

// Set stores a channel in the event bus
// errors on duplicate channel names
func (b *simpleBus) Set(evt string, ch any) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	// we dont allow overrites
	_, exists := b.channels[evt]

	if exists {
		return ErrChannelAlreadyExists
	}

	b.channels[evt] = ch
	return nil
}

// Subscribe to the channel in the bus that matches the type parameter
func Subscribe[T any](bus Bus) (Subscription[T], error) {
	typedCh, err := getOrCreateTypedChannel[T](bus)
	if err != nil {
		return nil, err
	}

	return typedCh.Subscribe()
}

// Publish to the channel in the bus that matches the type parameter
func Publish[T any](bus Bus, evt T) error {
	typedCh, err := getOrCreateTypedChannel[T](bus)
	if err != nil {
		return err
	}

	typedCh.Publish(evt)
	return nil
}

// Get the name of a event from the type parameter
func generateEventName[T any]() string {
	var t T

	// struct
	name := fmt.Sprintf("%T", t)
	if name != "<nil>" {
		return name
	}

	// interface
	return fmt.Sprintf("%T", new(T))
}

// Get the typed channel that matches the given type parameter
func getOrCreateTypedChannel[T any](bus Bus) (Channel[T], error) {
	evtName := generateEventName[T]()
	ch, exists := bus.Get(evtName)

	var typedCh Channel[T]
	var ok bool
	if !exists {
		typedCh = New[T](0, 100)
		if err := bus.Set(evtName, typedCh); err != nil {
			return nil, err
		}
	} else if typedCh, ok = ch.(Channel[T]); !ok {
		return nil, ErrInvalidChannelType
	}

	return typedCh, nil
}
