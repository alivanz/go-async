package async

func Go(r *Routine, fn func()) *Promise[bool] {
	var prom Promise[bool]
	r.GoFork(func(r *Routine) {
		fn()
		prom.Resolve(r, true)
	})
	return &prom
}

func GoParam[P any](r *Routine, fn func(p P), p P) *Promise[bool] {
	var prom Promise[bool]
	r.GoFork(func(r *Routine) {
		fn(p)
		prom.Resolve(r, true)
	})
	return &prom
}

func GoResult[R any](r *Routine, fn func() (R, error)) *Promise[R] {
	var prom Promise[R]
	r.GoFork(func(r *Routine) {
		result, err := fn()
		if err == nil {
			prom.Resolve(r, result)
		} else {
			prom.Reject(r, err)
		}
	})
	return &prom
}

func GoFunc[P, R any](r *Routine, fn func(p P) (R, error), p P) *Promise[R] {
	var prom Promise[R]
	r.GoFork(func(r *Routine) {
		result, err := fn(p)
		if err == nil {
			prom.Resolve(r, result)
		} else {
			prom.Reject(r, err)
		}
	})
	return &prom
}
