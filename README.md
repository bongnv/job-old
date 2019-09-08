# task
[![codecov](https://codecov.io/gh/bongnv/task/branch/master/graph/badge.svg)](https://codecov.io/gh/bongnv/task)
[![Go](https://github.com/bongnv/task/workflows/Go/badge.svg)](https://github.com/bongnv/task/actions)
[![GoDoc](https://godoc.org/github.com/bongnv/task?status.svg)](https://godoc.org/github.com/bongnv/task)
[![Go Report Card](https://goreportcard.com/badge/github.com/bongnv/task)](https://goreportcard.com/report/github.com/bongnv/task)
[![GolangCI](https://golangci.com/badges/github.com/golangci/golangci-lint.svg)](https://golangci.com)

Package `task` helps compose multiple smaller tasks to achieve a complex logic in order to achieve loose coupling and SRP.

## Motivation

In any program, logic grows with time. Without any mechanism to break into smaller parts, the program will become harder to maintain day by day. `task` provides a functionality to allow to break a complex business logic into smaller `Doer` and then easily compose them together.

## How to use

```go
import "github.com/bongnv/task"
```

### Implement `Doer`

Your application logic needs to be broken into smaller pieces. Each task is executed by a `Doer`.

```go
func doSomething(data metadata) task.DoFunc {
    return func(ctx context.Context) error {
        // doSomething handles one piece of application logic
        return nil
    }
}

func doAnotherThing(data metadata) task.DoFunc {
    return func(ctx context.Context) error {
        // doAnotherThing handles another piece of application logic
        return nil
    }
}
```

### Compose `Doer` together

In order to compose all `Doer` together, a `Job` is needed. A `Job` groups `Task` together where each task executes a `Doer`. In this example, let's say `doAnotherThing` requires `doSomething` to be finished.

```go
func applicationLogic(ctx context.Context, data metadata) error {
    taskDoSomething := task.Run(ctx, doSomething(data))
    taskDoAnotherThing := task.Run(ctx, doAnotherThing(data), taskDoSomething)
    return task.Wait(ctx, taskDoAnotherThing)
}
```

## FAQs

### What happens if a `Doer` panics?

`task` does not use ```recover()``` so panics will kill the process like normal. 

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details

