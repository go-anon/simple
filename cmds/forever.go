package cmds

import (
	"context"
	"runtime/debug"

	"github.com/go-anon/simple/logs"
)

// RunFunc defines the type for functions that will be executed in a resilient manner.
// These functions are expected to perform their tasks continuously until the provided context is canceled.
// Any panics that occur within a RunFunc are recovered, and the function is retried, ensuring uninterrupted execution.
type RunFunc func(ctx context.Context)
type PanicFunc func()

// Run executes the given function in a resilient and continuous manner until the provided context is canceled.
// It ensures that the function, identified by the 'name' parameter for logging purposes, is executed repeatedly.
// In case of a panic within the function, Run recovers from the panic, logs the error, and retries the function execution.
// This mechanism allows for uninterrupted operation of critical functions that must run continuously in the background.
//
// Parameters:
//   - ctx: A context.Context used to control the cancellation of the function execution. When this context is canceled,
//     the function execution stops gracefully.
//   - name: A string that identifies the function being executed. This name is used in logging for easier identification
//     of log entries related to this function execution.
//   - f: The RunFunc to be executed. This function should contain the core logic that needs to run resiliently.
//
// Usage:
// Pass a context to control when the function should stop, a descriptive name for logging, and the function itself.
// The function will run continuously, handling any panics, until the context is canceled.
func Forever(ctx context.Context, name string, rf RunFunc, pf PanicFunc, done func()) {
	if rf == nil {
		logs.Error("Cannot run function ( %s ) because the provided function is nil", name)
		return
	}

	logs.Warn("Running function ( %s ) in a resilient way...", name)

	c := 0

	go func() {
		defer done()
		for {
			select {
			case <-ctx.Done():
				logs.Info("Context canceled, stopping execution of function %s", name)
				return

			default:
				func() {
					defer handlePanic(name, pf)
					c++
					logs.Info("Function ( %s ) is running, cicle %d  ...", name, c)

					// Execute the provided function
					rf(ctx)
				}()
			}
		}
	}()
}

// handlePanic is an internal function used by Run to recover from panics that occur within the resilient function execution.
// When a panic is caught, handlePanic logs the error along with a stack trace to aid in debugging, ensuring that the
// application remains stable and the function execution can be retried.
//
// Parameters:
//   - name: The name of the function where the panic occurred. This is used in the error log to identify the source
//     of the panic more easily.
//
// Note: This function is intended to be used with defer within the function execution loop in Run.
func handlePanic(name string, pf PanicFunc) {
	if r := recover(); r != nil {
		logs.Error("Recovered function ( %s ) from panic: %v \n", name, r)
		debug.PrintStack()
		println()

		if pf == nil {
			return
		}

		// Here you can perform any additional necessary actions after recovering from a panic
		pf()
	}
}
