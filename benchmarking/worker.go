package benchmarking

import (
	"database/sql"
	"fmt"
	"sync"
	"time"
)

type Worker struct {
	instance int
	db       *sql.DB
	In       chan []string
	out      chan time.Duration
	errCh    chan error
	wg       *sync.WaitGroup
}

func NewWorker(instance int, db *sql.DB, in chan []string, out chan time.Duration, errCh chan error, wg *sync.WaitGroup) *Worker {
	fmt.Printf("instantiating worker %d\n", instance)

	return &Worker{
		instance: instance,
		db:       db,
		In:       in,
		out:      out,
		errCh:    errCh,
		wg:       wg,
	}
}

const queryCpuUsage = `select
	host,
	DATE_TRUNC('minute', ts),
	max(usage),
	min(usage)
from
	cpu_usage
where
	host = $1
	and ts between $2 and $3
group by
	host,
	DATE_TRUNC('minute', ts);`

func (w *Worker) Process() {
	fmt.Printf("running worker %d\n", w.instance)

	for v := range w.In {
		t0 := time.Now()

		hostname := v[0]
		startTime := v[1]
		endTime := v[2]

		var (
			h    string
			date time.Time
			max  string
			min  string
		)
		err := w.db.QueryRow(queryCpuUsage, hostname, startTime, endTime).Scan(&h, &date, &max, &min)
		if err != nil {
			w.errCh <- err
			continue
		}

		t1 := time.Now()
		w.out <- t1.Sub(t0)
	}

	w.wg.Done()
}
