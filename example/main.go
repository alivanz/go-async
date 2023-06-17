package main

import (
	"log"
	"time"

	"github.com/alivanz/go-async"
)

func main() {
	var r async.Routine
	r.RunForever(routine)
}

func routine(r *async.Routine) {
	// fork
	log.Printf("fork")
	r.Fork(func(r *async.Routine) {
		log.Printf("from different thread")
	})
	log.Printf("forked")
	// await
	log.Printf("await begin")
	async.GoParam(r, time.Sleep, 3*time.Second).Await(r)
	log.Printf("await done")
	// pending
	async.GoParam(r, time.Sleep, 3*time.Second).Then(func(v bool) {
		log.Printf("forked after 3 secs")
	})
}
