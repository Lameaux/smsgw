package workers

import (
	"context"
	"fmt"
	"runtime/debug"
	"time"

	"euromoby.com/smsgw/internal/logger"
)

type Worker interface {
	Name() string
	Run() (bool, error)
	SleepTime() time.Duration
	MaxAttempts() int
}

type Runner struct {
	w   Worker
	ctx context.Context
}

func NewRunner(ctx context.Context, w Worker) *Runner {
	return &Runner{
		w, ctx,
	}
}

func (r *Runner) Exec() error {
	logger.Infow("worker started", "worker", r.w.Name())

	for {
		logger.Infow("worker is running", "worker", r.w.Name())
		r.executeTask()

		logger.Infow("worker is sleeping", "worker", r.w.Name())

		select {
		case <-r.ctx.Done():
			logger.Infow("worker stopped",
				"worker", r.w.Name(),
				"error", r.ctx.Err(),
			)

			return nil
		case <-time.After(r.w.SleepTime()):
		}
	}
}

func (r *Runner) executeTask() {
	defer r.recoverPanic()

	for {
		hasNext, err := r.w.Run()
		if err != nil {
			logger.Errorw("worker execution failed",
				"error", err,
				"worker", r.w.Name(),
			)

			return
		}

		if !hasNext {
			logger.Infow("worker found nothing to process", "worker", r.w.Name())

			return
		}
	}
}

func (r *Runner) recoverPanic() {
	if e := recover(); e != nil {
		err, ok := e.(error)
		if !ok {
			err = fmt.Errorf("%v", e) //nolint: goerr113
		}

		logger.Errorw("Oops! panic!",
			"error", err,
			"worker", r.w.Name(),
			"stacktrace", string(debug.Stack()),
		)
	}
}
