package async

import (
	"sync/atomic"

	"github.com/alivanz/go-lists"
)

type Routine struct {
	todo     lists.Queue[ThreadFunc]
	ready    int32
	parallel int32
	c        chan struct{}
}

func (r *Routine) RunForever(routine ThreadFunc) {
	r.c = make(chan struct{}, 1) // buffer is required
	r.Fork(routine)
	r.runForever()
}

func (r *Routine) runForever() {
	var fn ThreadFunc
	for {
		if atomic.LoadInt32(&r.ready) > 0 {
			for r.todo.Pop(&fn) {
				atomic.AddInt32(&r.ready, -1)
				fn(r)
			}
		}
		for atomic.LoadInt32(&r.parallel) == 0 {
			return
		}
		<-r.c
	}
}

func (r *Routine) run() {
	var fn ThreadFunc
	for r.todo.Pop(&fn) {
		atomic.AddInt32(&r.ready, -1)
		fn(r)
	}
}

func (r *Routine) wake() {
	select {
	default:
	case r.c <- struct{}{}:
	}
}

// fork WITHOUT actually creating a thread
// will be run in MAIN thread
func (r *Routine) Fork(fn ThreadFunc) {
	atomic.AddInt32(&r.ready, 1)
	r.todo.Push(fn)
}

// fork WITH actually creating a thread
func (r *Routine) GoFork(fn ThreadFunc) {
	atomic.AddInt32(&r.parallel, 1)
	go func() {
		fn(r)
		atomic.AddInt32(&r.parallel, -1)
		r.wake()
	}()
}
