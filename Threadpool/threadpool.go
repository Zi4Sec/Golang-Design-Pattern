package main

import "sync"

// Task defines a job to be executed by a worker.
type Task func()

// Worker represents a worker that processes tasks.
type Worker struct {
	id          int
	wg          *sync.WaitGroup
	taskChannel chan Task
	stopChannel chan bool
}

// Pool represents a pool of workers.
type Pool struct {
	numWorkers int
	workers    []Worker
	taskQueue  chan Task
	wg         *sync.WaitGroup
}

func NewPool(numWorkers int) *Pool {
	taskQueue := make(chan Task)
	stopChannel := make(chan bool)
	wg := &sync.WaitGroup{}

	pool := &Pool{
		numWorkers: numWorkers,
		taskQueue:  taskQueue,
		wg:         wg,
	}

	for i := 0; i < numWorkers; i++ {
		worker := NewWorker(i, pool.taskQueue, pool.wg, &stopChannel)
		pool.workers = append(pool.workers, worker)
		worker.Start()
	}
	return pool
}

func NewWorker(id int, taskChannel chan Task, wg *sync.WaitGroup, stopChannel *chan bool) Worker {
	return Worker{
		id:          id,
		taskChannel: taskChannel,
		wg:          wg,
		stopChannel: *stopChannel,
	}
}

func (w Worker) Start() {
	go func() {
		for {
			select {
			case task := <-w.taskChannel:
				if task != nil {
					task()
					w.wg.Done()
				}

			case <-w.stopChannel:
				return
			}
		}
	}()
}

func (p *Pool) AddTask(task Task) {
	p.wg.Add(1)
	p.taskQueue <- task
}

func (p *Pool) Stop() {
	close(p.taskQueue)
	for _, worker := range p.workers {
		worker.stopChannel <- true
	}
	p.wg.Wait()
}
