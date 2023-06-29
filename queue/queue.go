package queue

import (
	"github.com/gin-gonic/gin"
	"sync"
)

type Iterm struct {
	C      *gin.Context
	Cmd    string
	Result string
}

type Queue struct {
	lock sync.Mutex
	cond *sync.Cond
	data []*Iterm
}

func NewQueue(size int) *Queue {
	q := &Queue{
		data: make([]*Iterm, 0, size),
	}
	q.cond = sync.NewCond(&q.lock)
	return q
}

func (q *Queue) Push(item *Iterm) {
	q.lock.Lock()
	defer q.lock.Unlock()

	q.data = append(q.data, item)
	q.cond.Signal()
}

func (q *Queue) Pop() *Iterm {
	q.lock.Lock()
	defer q.lock.Unlock()

	for len(q.data) == 0 {
		q.cond.Wait()
	}

	item := q.data[0]
	q.data = q.data[1:]

	return item
}
