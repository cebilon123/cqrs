package cqrs

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

// Query is the query that is dispatched to the query handler.
type Query struct {
	Key     Key     // Key is unique identifier of the query
	Payload Payload // Payload is the query payload
}

// QueryHandlerFunc is the function that handles the query.
type QueryHandlerFunc func(ctx context.Context, query Query) (any, error)

var (
	queryHandlers = make(map[Key]QueryHandlerFunc)
	queryMutex    sync.Mutex
)

var ErrQueryHandlerFuncNotFound = errors.New("query handler func for given query not found")

// RegisterQueryHandlerFuncs registers the given query handler function for the given query key.
func RegisterQueryHandlerFuncs(queryKey Key, h QueryHandlerFunc) {
	queryMutex.Lock()
	defer queryMutex.Unlock()

	queryHandlers[queryKey] = h
}

// DispatchQuery dispatches the given query to the query handler.
func DispatchQuery[T any](ctx context.Context, query Query) (*T, error) {
	h, ok := queryHandlers[query.Key]
	if !ok {
		return nil, fmt.Errorf("%w; %s", ErrQueryHandlerFuncNotFound, query.Key)
	}

	res, err := h(ctx, query)
	if err != nil {
		return nil, err
	}

	return tryMapResult[T](res)
}

func tryMapResult[T any](res any) (*T, error) {
	if res == nil {
		return nil, errors.New("result is nil")
	}

	v, ok := res.(T)
	if !ok {
		return nil, fmt.Errorf("result is not of type %T", v)
	}

	return &v, nil
}
