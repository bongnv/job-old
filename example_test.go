package task_test

import (
	"context"
	"fmt"
	"sync/atomic"

	"github.com/bongnv/task"
)

func printNumber(number int) task.DoFunc {
	return func(_ context.Context) error {
		fmt.Println(number)
		return nil
	}
}

func ExampleJob_sequential() {
	ctx := context.Background()
	j := task.NewJob()
	taskPrint1 := j.Run(ctx, printNumber(1))
	taskPrint2 := j.Run(ctx, printNumber(2), taskPrint1)
	_ = j.Run(ctx, printNumber(3), taskPrint2)
	err := j.Wait(ctx)
	fmt.Println(err)
	// Output:
	// 1
	// 2
	// 3
	// <nil>
}

func counter(number *int64) task.DoFunc {
	return func(_ context.Context) error {
		atomic.AddInt64(number, 1)
		return nil
	}
}

func ExampleJob_concurrent() {
	ctx := context.Background()
	var total int64
	j := task.NewJob()
	j.Run(ctx, counter(&total))
	j.Run(ctx, counter(&total))
	j.Run(ctx, counter(&total))
	err := j.Wait(ctx)
	fmt.Println(err)
	fmt.Println(total)
	// Output:
	// <nil>
	// 3
}
