// Package cqrs provides a simple implementation of the command query responsibility segregation pattern.
package cqrs

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

// Payload is the command payload. It can be any type.
// It can be then deserialized to the desired type using the TryMapPayload function.
type Payload any

// Command is the command that is dispatched to the command handler.
type Command struct {
	Key     Key     // Key is the unique identifier of the command.
	Payload Payload // Payload is the command payload.
}

// CommandHandlerFunc is the function that handles the command.
type CommandHandlerFunc func(ctx context.Context, cmd Command) error

var (
	// commandHandlers is a map of command keys to a slice of command handler functions.
	commandHandlers = make(map[Key][]CommandHandlerFunc)
	commandMutex    sync.Mutex
)

// RegisterCommandHandlerFuncs registers the given command handler functions for the given command key.
func RegisterCommandHandlerFuncs(cmdKey Key, h ...CommandHandlerFunc) {
	commandMutex.Lock()
	defer commandMutex.Unlock()

	commandHandlers[cmdKey] = append(commandHandlers[cmdKey], h...)
}

var ErrCommandHandlerFuncNotFound = errors.New("command handler func for given command not found")

// DispatchCommand dispatches the given command to the command handler.
func DispatchCommand(ctx context.Context, cmd Command) error {
	v, ok := commandHandlers[cmd.Key]
	if !ok {
		return fmt.Errorf("%w; %s", ErrCommandHandlerFuncNotFound, cmd.Key)
	}

	for _, h := range v {
		if err := h(ctx, cmd); err != nil {
			return err
		}
	}

	return nil
}
