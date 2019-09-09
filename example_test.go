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

func Example_sequential() {
	ctx := context.Background()
	taskPrint1 := task.Run(ctx, printNumber(1))
	taskPrint2 := task.Run(ctx, printNumber(2), taskPrint1)
	taskPrint3 := task.Run(ctx, printNumber(3), taskPrint2)
	err := task.Wait(ctx, taskPrint3)
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

func Example_concurrent() {
	ctx := context.Background()
	var total int64
	t1 := task.Run(ctx, counter(&total))
	t2 := task.Run(ctx, counter(&total))
	t3 := task.Run(ctx, counter(&total))
	err := task.Wait(ctx, t1, t2, t3)
	fmt.Println(err)
	fmt.Println(total)
	// Output:
	// <nil>
	// 3
}
