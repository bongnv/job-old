# job
[![codecov](https://codecov.io/gh/bongnv/job/branch/master/graph/badge.svg)](https://codecov.io/gh/bongnv/job)
[![Go](https://github.com/bongnv/job/workflows/Go/badge.svg)](https://github.com/bongnv/job/actions)
[![GoDoc](https://godoc.org/github.com/bongnv/job?status.svg)](https://godoc.org/github.com/bongnv/job)
[![Go Report Card](https://goreportcard.com/badge/github.com/bongnv/task)](https://goreportcard.com/report/github.com/bongnv/job)
[![GolangCI](https://golangci.com/badges/github.com/golangci/golangci-lint.svg)](https://golangci.com)

Package `job` helps compose multiple smaller tasks to achieve a complex logic while archieving loose coupling and SRP between tasks.

## Motivation

In any program, logic grows with time. Without any mechanism to break into smaller parts, the program will become harder to maintain day by day. `job` provides a functionality to allow to break a complex business logic into smaller `Task` and then easily compose them together.

## How to use

```go
import "github.com/bongnv/job"
```

### Implement `Task`

Your application logic needs to be broken into smaller pieces. Each task will be executed separately.

```go
func doSomething(data metadata) job.TaskFunc {
    return func(ctx context.Context) error {
        // doSomething handles one piece of application logic
        return nil
    }
}

func doAnotherThing(data metadata) job.DoFunc {
    return func(ctx context.Context) error {
        // doAnotherThing handles another piece of application logic
        return nil
    }
}
```

### Compose `Task` together

In order to compose all `Task` together, a `Job` is needed. A `Job` groups `Task` together. In this example, let's say `doAnotherThing` requires `doSomething` to be finished.

```go
func applicationLogic(ctx context.Context, data metadata) error {
    j := job.New()
    taskDoSomething := j.Start(ctx, doSomething(data))
    taskDoAnotherThing := j.Start(ctx, doAnotherThing(data), taskDoSomething)
    return j.Wait(ctx)
}
```

## FAQs

TODO

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details

