package benchmarking

import (
	"database/sql"
	"fmt"
	"sync"
	"time"
)

type Benchmarking struct {
	processed chan bool
	pool      map[int]*Worker
	out       chan time.Duration
	wg        *sync.WaitGroup
}

func NewBenchmarking(processed chan bool, numOfWorkers int, db *sql.DB, rows chan []string, errCh chan error) *Benchmarking {
	out := make(chan time.Duration)
	pool := make(map[int]*Worker)
	var wg sync.WaitGroup
	wg.Add(numOfWorkers)

	for i := 0; i < numOfWorkers; i++ {
		w := NewWorker(i+1, db, rows, out, errCh, &wg)
		pool[i+1] = w
	}

	return &Benchmarking{
		processed: processed,
		pool:      pool,
		out:       out,
		wg:        &wg,
	}
}

func (b *Benchmarking) Process() {
	for _, v := range b.pool {
		go v.BenchmarkQuery()
	}

	go func() {
		for r := range b.out {
			fmt.Println("aa", r)
		}
	}()

	b.wg.Wait()
	close(b.processed)
}
