package eventbus

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBusBasic(t *testing.T) {
	bus := NewBus()

	// subscribe
	subCh, err := Subscribe[int](bus)
	require.NoError(t, err)

	// publish
	err = Publish(bus, 1)
	require.NoError(t, err)

	item := <-subCh
	require.Equal(t, 1, item)
}
