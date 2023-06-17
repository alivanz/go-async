package async

type Promise[T any] struct {
	state int32
	then  func(v T)
	catch func(e error)
}

// resolve promise
// will be run in MAIN thread
// can only be invoked ONCE
func (p *Promise[T]) Resolve(r *Routine, v T) {
	r.Fork(func(r *Routine) {
		if p.state != 0 {
			return
		}
		p.then(v)
		p.state = 1
	})
}

// reject promise
// will be run in MAIN thread
// can only be invoked ONCE
func (p *Promise[T]) Reject(r *Routine, err error) {
	r.Fork(func(r *Routine) {
		if p.state != 0 {
			return
		}
		p.catch(err)
		p.state = 2
	})
}

// will be run in MAIN thread
func (p *Promise[T]) Then(fn func(v T)) *Promise[T] {
	p.then = fn
	return p
}

// will be run in MAIN thread
func (p *Promise[T]) Catch(fn func(err error)) *Promise[T] {
	p.catch = fn
	return p
}

// block and get promise result
// must be called from MAIN thread
func (p *Promise[T]) Await(r *Routine) (*T, error) {
	var (
		out  T
		err  error
		done bool
	)
	p.then = func(v T) {
		out = v
		done = true
	}
	p.catch = func(e error) {
		err = e
		done = true
	}
	for !done {
		r.run()
	}
	if err != nil {
		return nil, err
	}
	return &out, nil
}
