package common

import (
	"context"
)

// Runnable should be executed in a goroutine.
type Runnable interface {
	// Run is the method to be called in the goroutine.
	Run(context.Context) error
}
