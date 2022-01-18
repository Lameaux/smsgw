package workers

import (
	"time"

	"euromoby.com/smsgw/internal/logger"
)

type Worker interface {
	Name() string
	Run() (bool, error)
	SleepTime() time.Duration
}

type Runner struct {
	w    Worker
	done chan bool
}

func NewRunner(w Worker) *Runner {
	return &Runner{
		w:    w,
		done: make(chan bool),
	}
}

func (r *Runner) Start() {
	logger.Infow("worker started", "worker", r.w.Name())
	for {
		logger.Infow("worker is running", "worker", r.w.Name())
		r.executeTask()

		logger.Infow("worker is sleeping", "worker", r.w.Name())

		select {
		case <-r.done:
			logger.Infow("worker stopped", "worker", r.w.Name())
			return
		case <-time.After(r.w.SleepTime()):
		}
	}
}

func (r *Runner) Stop() {
	close(r.done)
}

func (r *Runner) executeTask() {
	for {
		hasNext, err := r.w.Run()

		if err != nil {
			logger.Error(err)
			return
		}
		if !hasNext {
			logger.Infow("worker found nothing to process", "worker", r.w.Name())
			return
		}
	}
}
