package main

import (
	"context"
	"github.com/cebilon123/cqrs"
	"log"
)

const (
	commandKey = "commandKey"
	queryKey   = "queryKey"
)

type CommandTestPayload struct {
	FirstName string
}

type QueryTestPayload struct {
	FirstName string
}
type QueryTestResult struct {
	LastNameAsFirstName string
}

func main() {
	cqrs.RegisterCommandHandlerFuncs(commandKey, func(ctx context.Context, cmd cqrs.Command) error {
		// TryMapPayload is a helper function that tries to map the payload to the given type. It returns an error
		//if the payload is nil or if the payload is not of the given type.
		payload, err := cqrs.TryMapPayload[CommandTestPayload](cmd.Payload)
		if err != nil {
			return err
		}

		println(payload.FirstName)

		return nil
	})
	cqrs.RegisterQueryHandlerFuncs(queryKey, func(ctx context.Context, query cqrs.Query) (any, error) {
		payload, err := cqrs.TryMapPayload[QueryTestPayload](query.Payload)
		if err != nil {
			return nil, err
		}

		println(payload.FirstName)

		return QueryTestResult{LastNameAsFirstName: payload.FirstName}, nil
	})

	ctx := context.Background()
	err := cqrs.DispatchCommand(ctx, cqrs.Command{
		Key: commandKey,
		Payload: CommandTestPayload{
			FirstName: "John",
		},
	})
	if err != nil {
		log.Print(err)
	}

	res, err := cqrs.DispatchQuery[QueryTestResult](ctx, cqrs.Query{
		Key: queryKey,
		Payload: QueryTestPayload{
			FirstName: "Andrew",
		},
	})

	if err != nil {
		log.Print(err)
	}

	println(res.LastNameAsFirstName)
}
