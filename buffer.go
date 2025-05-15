package main

import (
	"context"
	"errors"
	"sync"
	"time"
)

var (
	timeout = 10 * time.Second
)

type TaskWithErrCh[I, O any] struct {
	task     I
	resultCh chan *O
	errCh    chan error
}

func (t *TaskWithErrCh[I, O]) GetResult() (*O, error) {
	select {
	case err := <-t.errCh:
		return nil, err
	case result := <-t.resultCh:
		return result, nil
	}
}

type Buffer[I, O any] struct {
	bufferCh chan *TaskWithErrCh[I, O]
	notifier chan struct{}

	semaphore chan struct{}

	write func(task []I) ([]*O, error)
}

func NewBuffer[I, O any](
	ctx context.Context,
	write func(task []I) ([]*O, error),
) *Buffer[I, O] {
	b := &Buffer[I, O]{
		bufferCh:  make(chan *TaskWithErrCh[I, O], batchSize*continuousWritersCount),
		notifier:  make(chan struct{}, 1), // buffered to notify even when a flush is in progress
		write:     write,
		semaphore: make(chan struct{}, continuousWritersCount),
	}

	b.startFlusher(ctx)

	return b
}

func (b *Buffer[I, O]) WriteNoWait(task I) (*TaskWithErrCh[I, O], error) {
	resultCh := make(chan *O)
	errCh := make(chan error)

	taskWithErrCh := &TaskWithErrCh[I, O]{
		task:     task,
		resultCh: resultCh,
		errCh:    errCh,
	}

	select {
	case b.bufferCh <- taskWithErrCh:
	case <-time.After(timeout):
		return nil, errors.New("timeout while writing to buffer")
	}

	// notify the flusher that there is a new message
	select {
	case b.notifier <- struct{}{}:
	default:
	}

	return taskWithErrCh, nil
}

func (b *Buffer[I, O]) Write(task I) (*O, error) {
	taskWithErrCh, err := b.WriteNoWait(task)

	if err != nil {
		return nil, err
	}

	return taskWithErrCh.GetResult()
}

func (b *Buffer[I, O]) startFlusher(ctx context.Context) {
	ticker := time.NewTicker(flushInterval)

	go func() {
		for {
			select {
			case <-ctx.Done():
				b.flush()
				return
			case <-ticker.C:
				go b.flush()
			case <-b.notifier:
				go b.flush()
			}
		}
	}()
}

func (b *Buffer[I, O]) flush() {
	wg := sync.WaitGroup{}

outer:
	for i := 0; i < continuousWritersCount; i++ {
		select {
		case b.semaphore <- struct{}{}:
		default:
			break outer
		}

		wg.Add(1)

		go func() {
			defer wg.Done()

			startedFlush := time.Now()

			defer func() {
				go func() {
					<-time.After(flushInterval - time.Since(startedFlush))
					<-b.semaphore
				}()
			}()

			msgsWithChs := make([]*TaskWithErrCh[I, O], 0)
			tasks := make([]I, 0)

			// read all messages currently in the buffer
			for i := 0; i < batchSize; i++ {
				select {
				case msg := <-b.bufferCh:
					msgsWithChs = append(msgsWithChs, msg)

					tasks = append(tasks, msg.task)
				default:
					i = batchSize
				}
			}

			if len(tasks) == 0 {
				return
			}

			tasksOut, err := b.write(tasks)

			if err != nil {
				for _, msgWithErrCh := range msgsWithChs {
					msgWithErrCh.errCh <- err
				}

				return
			}

			if len(tasksOut) != len(tasks) {
				err := errors.New("number of tasks out does not match number of tasks in")
				for _, msgWithErrCh := range msgsWithChs {
					msgWithErrCh.errCh <- err
				}

				return
			}

			for i, msgWithErrCh := range msgsWithChs {
				msgWithErrCh.resultCh <- tasksOut[i]
			}

		}()
	}

	wg.Wait()
}
