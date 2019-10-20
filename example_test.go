package job_test

import (
	"context"
	"fmt"
	"sync/atomic"

	"github.com/bongnv/job"
)

func printNumber(number int) job.TaskFunc {
	return func(_ context.Context) error {
		fmt.Println(number)
		return nil
	}
}

func Example_sequential() {
	ctx := context.Background()
	j := job.New()
	taskPrint1 := j.Start(ctx, printNumber(1))
	taskPrint2 := j.Start(ctx, printNumber(2), taskPrint1)
	_ = j.Start(ctx, printNumber(3), taskPrint2)
	err := j.Wait(ctx)
	fmt.Println(err)
	// Output:
	// 1
	// 2
	// 3
	// <nil>
}

func counter(number *int64) job.TaskFunc {
	return func(_ context.Context) error {
		atomic.AddInt64(number, 1)
		return nil
	}
}

func Example_concurrent() {
	j := job.New()
	ctx := context.Background()
	var total int64
	_ = j.Start(ctx, counter(&total))
	_ = j.Start(ctx, counter(&total))
	_ = j.Start(ctx, counter(&total))
	err := j.Wait(ctx)
	fmt.Println(err)
	fmt.Println(total)
	// Output:
	// <nil>
	// 3
}
