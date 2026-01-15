package closer

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/delyke/gophermat_bonus_system/internal/logger"
)

func TestCloserCloseAllExecutesInReverseOrder(t *testing.T) {
	c := NewWithLogger(&logger.NoopLogger{})
	var mu sync.Mutex
	var order []string

	c.Add(
		func(context.Context) error {
			mu.Lock()
			order = append(order, "first")
			mu.Unlock()
			return nil
		},
		func(context.Context) error {
			mu.Lock()
			order = append(order, "second")
			mu.Unlock()
			return nil
		},
	)

	err := c.CloseAll(context.Background())

	require.NoError(t, err)
	require.Len(t, order, 2)
	require.Equal(t, "second", order[0])
	require.Equal(t, "first", order[1])
}

func TestCloserCloseAllReturnsFirstError(t *testing.T) {
	c := NewWithLogger(&logger.NoopLogger{})
	expectedErr := errors.New("close error")

	c.Add(
		func(context.Context) error { return nil },
		func(context.Context) error { return expectedErr },
	)

	err := c.CloseAll(context.Background())

	require.Error(t, err)
}

func TestCloserCloseAllHandlesPanic(t *testing.T) {
	c := NewWithLogger(&logger.NoopLogger{})
	c.Add(func(context.Context) error { panic("boom") })

	err := c.CloseAll(context.Background())

	require.Error(t, err)
}

func TestCloserCloseAllContextCanceled(t *testing.T) {
	c := NewWithLogger(&logger.NoopLogger{})
	c.Add(func(context.Context) error {
		time.Sleep(20 * time.Millisecond) //nolint:forbidigo
		return nil
	})

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := c.CloseAll(ctx)

	require.ErrorIs(t, err, context.Canceled)
}

func TestCloserCloseAllRunsOnce(t *testing.T) {
	c := NewWithLogger(&logger.NoopLogger{})
	var mu sync.Mutex
	count := 0

	c.Add(func(context.Context) error {
		mu.Lock()
		count++
		mu.Unlock()
		return nil
	})

	err := c.CloseAll(context.Background())
	require.NoError(t, err)

	err = c.CloseAll(context.Background())
	require.NoError(t, err)

	mu.Lock()
	defer mu.Unlock()
	require.Equal(t, 1, count)
}
