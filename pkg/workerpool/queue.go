package workerpool

import "sync"

type node struct {
	task Task
	next *node
}

type queue struct {
	head   *node
	tail   *node
	mu     sync.Mutex //just for thread safe but consider again
	length int        //in future if we need a limit of queue
}

func newQueue() *queue {
	return &queue{
		head: nil,
		tail: nil,
		mu:   sync.Mutex{},
	}
}

func (q *queue) enqueue(t Task) {
	newNode := &node{
		task: t,
		next: nil,
	}

	q.mu.Lock()
	defer q.mu.Unlock()

	if q.head == nil {
		q.head = newNode
		q.tail = q.head
	} else {
		q.tail.next = newNode
		q.tail = newNode
	}

	q.length++
}

func (q *queue) dequeue() Task {
	if q.head == nil {
		return nil
	}

	q.mu.Lock()
	defer q.mu.Unlock()

	t := q.head.task
	q.head = q.head.next

	q.length--

	return t
}

func (q *queue) len() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	return q.length
}
