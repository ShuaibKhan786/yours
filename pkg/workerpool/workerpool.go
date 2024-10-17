package workerpool

import (
	"context"
)

type Task interface {
	Execute() error
}

type WorkerPool struct {
	task      chan Task
	err       chan error
	queue     *queue
	ctx       context.Context
	cancelCtx func()
	wpCount   int
}

func NewWorkerPool(wpCount int) *WorkerPool {
	ctx, cancel := context.WithCancel(context.Background())

	wp := &WorkerPool{
		task:      make(chan Task),
		err:       make(chan error),
		queue:     newQueue(),
		ctx:       ctx,
		cancelCtx: cancel,
		wpCount:   wpCount,
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
		select {
		case <-assignedChannel:
			t := wp.queue.dequeue()
			if t != nil {
				err := t.Execute()
				if err != nil {
					wp.err <- err
				}
			}
			notifyDone <- id
		case <-wp.ctx.Done():
			return
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
	wp.cancelCtx()
	close(wp.task)
	close(wp.err)
}
