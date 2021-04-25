package queue

import (
	"sync"
	"time"

)

type Queue struct {
	mu             sync.Mutex
	workQueue      chan string
	inProcessing   sync.Map
	inQueue        sync.Map
	inCompensation sync.Map
	worker         int
	workFunc       WorkFunc
	onSuccess      CallBackFunc
	onFail         CallBackFunc
	stopCh         <-chan struct{}
}

type WorkFunc func(obj interface{}) error

type CallBackFunc func(obj interface{})

func NewQueue(capacity, worker int, workFunc WorkFunc, onSuccess, onFail CallBackFunc, stopCh <-chan struct{}) *Queue {
	return &Queue{
		workQueue:    make(chan string, capacity),
		worker:       worker,
		inQueue:      sync.Map{},
		inProcessing: sync.Map{},
		workFunc:     workFunc,
		onSuccess:    onSuccess,
		onFail:       onFail,
		stopCh:       stopCh,
	}
}

func (q *Queue) Run() {

	wg := &sync.WaitGroup{}

	for i := 0; i < q.worker; i++ {
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			q.run()
		}(wg)
	}

	wg.Wait()
}

func (q *Queue) run() {
	for {
		select {
		case <-q.stopCh:
			return
		default:
			key, obj := q.pop()
			_, loaded := q.inProcessing.LoadOrStore(key, obj)
			if loaded {
				q.Push(key, obj)
				continue
			}

			if err := q.workFunc(obj); err != nil {
				if q.onFail != nil {
					q.onFail(obj)
				}
			} else {
				if q.onSuccess != nil {
					q.onSuccess(obj)
				}
			}
			q.inProcessing.Delete(key)
		}
	}
}

func (q *Queue) Push(key string, value interface{}) {

	q.mu.Lock()
	_, ok := q.inQueue.Load(key)
	q.inQueue.Store(key, value)
	q.mu.Unlock()

	if !ok {
		q.workQueue <- key
	}

}

func (q *Queue) PushAfter(key string, value interface{}, after time.Duration) {

	_, loaded := q.inCompensation.LoadOrStore(key, value)
	if loaded {
		return
	}

	go func() {
		time.Sleep(after)
		q.Push(key, value)
		q.inCompensation.Delete(key)
	}()

}

func (q *Queue) pop() (string, interface{}) {
	q.mu.Lock()
	defer q.mu.Unlock()

	key := <-q.workQueue
	obj, _ := q.inQueue.LoadAndDelete(key)

	return key, obj
}
