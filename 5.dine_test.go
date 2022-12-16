package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// go test -timeout 30s -run ^TestDine$ go-concurrency -v -count=1 -race
func TestDine(t *testing.T) {
	eatTime = 10 * time.Millisecond
	thinkTime = 10 * time.Millisecond
	for i := 0; i < 10; i++ {
		orderSlice = []string{}
		dine()
		require.Equal(t, 5, len(orderSlice))
	}
}
