# CQRS

---

### What is CQRS?
CQRS stands for Command Query Responsibility Segregation. It is a pattern that separates read and write 
models in your application. It is a good practice to separate read and write models. It is also handy when
you are doing the DDD (Domain Driven Design) approach.


### About
This is a simple implementation of CQRS pattern in GO. It is made without any external dependencies.
It was made for one of my projects and I decided to share it with the community. Of course, there could be bugs here
and there, but I will try to fix them as soon as possible and I will be happy if you report them to me
or make a pull request.

### How to use

Usage of the library is very simple, and you don't need to instantiate any objects, you just need to register handlers
in the cqrs package and then dispatch commands and queries.
```go
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

```

Result:
```shell
John
Andrew
Andrew
```

It is hard to make it full generic-like as in C# or Java, but I tried to make it as generic as possible.