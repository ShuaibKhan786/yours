//ADVICE: always check the root context error before sending something
//        to a channel as root context might called done when an application
//        needs resource cleanup or in short when app closed
//        which interns might closed the channel
//        and if you send something on the closed channel it gets panic

package workerpool

import (
	"context"
	"time"
)

type Task interface {
	Execute() error
}

type dummyTask struct{}

func (t *dummyTask) Execute() error {
	return nil
}

type WorkerPool struct {
	task    chan Task
	err     chan error
	queue   *queue
	ctx     context.Context //root context
	wpCount int
}

func NewWorkerPool(ctx context.Context, wpCount int) *WorkerPool {
	wp := &WorkerPool{
		task:    make(chan Task),
		err:     make(chan error),
		queue:   newQueue(),
		ctx:     ctx,
		wpCount: wpCount,
	}

	wp.initWP()

	return wp
}

func (wp *WorkerPool) initWP() {
	notifyDone := make(chan int)
	avilableWorker := make(chan int, wp.wpCount) //as channel is queue
	wpChannels := make([]chan struct{}, wp.wpCount)

	//inshort each worker will gets its own channel
	//just to assigned a tasked in a RoundRobin fhasioned
	for id := range wpChannels {
		wpChannels[id] = make(chan struct{})
		go wp.worker(id, notifyDone, wpChannels[id])
		avilableWorker <- id
	}

	//workerpool manager go routine
	go func() {
		ticker := time.NewTicker(10 * time.Minute)
		defer ticker.Stop()
		for {
			select {
			case t := <-wp.task:
				wp.queue.enqueue(t)
				//if there is any avilable worker
				//it will assigned a task else
				//this manager will be in block for a worker to be free
				select {
				case workerID := <-avilableWorker:
					wpChannels[workerID] <- struct{}{}
				default:

				}
			case workerID := <-notifyDone:
				//there can be no task in the queue
				//so I need to first check
				if wp.queue.len() != 0 {
					wpChannels[workerID] <- struct{}{}
				} else {
					avilableWorker <- workerID
				}
			case <-ticker.C:
				if wp.queue.len() == 0 && wp.ctx.Err() == nil {
					for i := 0; i < wp.wpCount; i++ {
						select {
						case wp.task <- &dummyTask{}:
						default:
						}
					}
				}
			case <-wp.ctx.Done():
				for _, wpChannel := range wpChannels {
					close(wpChannel)
				}
				close(notifyDone)
				close(avilableWorker)
				return
			}
		}
	}()
}

func (wp *WorkerPool) worker(id int, notifyDone chan int, assignedChannel chan struct{}) {
	for {
		//priortizing ctx first
		select {
		case <-wp.ctx.Done():
			return
		default:
		}

		select {
		case <-wp.ctx.Done():
			return
		case <-assignedChannel:
			t := wp.queue.dequeue()
			if t != nil {
				err := t.Execute()
				if err != nil {
					wp.err <- err
				}
			}
			//checking wether the root context is being cancelled
			if wp.ctx.Err() != nil {
				return
			}
			notifyDone <- id
		}
	}
}

func (wp *WorkerPool) Add(task Task) {
	wp.task <- task
}

func (wp *WorkerPool) ErrorChannel() chan error {
	return wp.err
}

func (wp *WorkerPool) Close() {
	close(wp.task)
	close(wp.err)
}
