package local

import (
	"github.com/satori/go.uuid"
	"github.com/spiral/jobs"
	"sync"
	"time"
	"sync/atomic"
)

// Broker run queue using local goroutines.
type Broker struct {
	mu          sync.Mutex
	wg          sync.WaitGroup
	queue       chan entry
	handlerPool chan jobs.Handler
	err         jobs.ErrorHandler
	stat        *jobs.PipelineStat
}

type entry struct {
	id      string
	attempt int
	job     *jobs.Job
}

// Listen configures broker with list of pipelines to listen and handler function. Broker broker groups all pipelines
// together.
func (b *Broker) Listen(pipelines []*jobs.Pipeline, pool chan jobs.Handler, err jobs.ErrorHandler) error {
	b.handlerPool = pool
	b.err = err
	return nil
}

// Init configures local job broker.
func (b *Broker) Init() (bool, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.queue = make(chan entry)
	b.stat = &jobs.PipelineStat{Name: "local", Details: "in-memory"}

	return true, nil
}

// Serve local broker.
func (b *Broker) Serve() error {
	b.mu.Lock()
	b.queue = make(chan entry)
	b.mu.Unlock()

	var handler jobs.Handler
	for q := range b.queue {
		id, job := q.id, q.job

		if job.Options.Delay != 0 {
			time.Sleep(job.Options.DelayDuration())
		}

		// wait for free handler
		handler = <-b.handlerPool

		go func() {
			err := handler(id, job)
			b.handlerPool <- handler

			if err == nil {
				atomic.AddInt64(&b.stat.Completed, 1)
				return
			}

			if !job.CanRetry() {
				b.err(id, job, err)
				atomic.AddInt64(&b.stat.Failed, 1)
				return
			}

			if job.Options.RetryDelay != 0 {
				time.Sleep(job.Options.RetryDuration())
			}

			b.queue <- entry{id: id, job: job}
		}()
	}

	return nil
}

// Stop local broker.
func (b *Broker) Stop() {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.queue != nil {
		close(b.queue)
	}
}

// Push new job to queue
func (b *Broker) Push(p *jobs.Pipeline, j *jobs.Job) (string, error) {
	id, err := uuid.NewV4()
	if err != nil {
		return "", err
	}

	go func() { b.queue <- entry{id: id.String(), job: j} }()

	atomic.AddInt64(&b.stat.Total, 1)
	if j.Options.Delay != 0 {
		// todo: must be interactive
		atomic.AddInt64(&b.stat.Delayed, 1)
	}

	return id.String(), nil
}

// Stat must fetch statistics about given pipeline or return error.
func (b *Broker) Stat(p *jobs.Pipeline) (stats *jobs.PipelineStat, err error) {
	return b.stat, nil
}
