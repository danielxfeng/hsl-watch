package service

import (
	"context"
	"log"
	"time"
)

const Interval = 10 * time.Second

type Task func(context.Context) error

// Schedule runs task immediately and then every Interval until ctx is cancelled.
// Fail open: task errors are logged and never stop the schedule.
func Schedule(ctx context.Context, task Task) {
	ticker := time.NewTicker(Interval)

	go func() {
		defer ticker.Stop()

		run(ctx, task)

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if ctx.Err() != nil {
					return
				}
				run(ctx, task)
			}
		}
	}()
}

func run(ctx context.Context, task Task) {
	if err := task(ctx); err != nil {
		log.Printf("scheduled task: %v", err)
	}
}
