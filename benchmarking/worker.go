package benchmarking

import (
	"database/sql"
	"fmt"
	"log"
	"sync"
	"time"
)

type Worker struct {
	instance int
	db       *sql.DB
	in       chan []string
	out      chan time.Duration
	errCh    chan error
	group    *sync.WaitGroup
}

func NewWorker(instance int, db *sql.DB, rows chan []string, out chan time.Duration, errCh chan error, group *sync.WaitGroup) *Worker {
	log.Println("instantiating worker ", instance)

	return &Worker{
		instance: instance,
		db:       db,
		in:       rows,
		out:      out,
		errCh:    errCh,
		group:    group,
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

func (w *Worker) BenchmarkQuery() {
	log.Println("running worker ", w.instance)

	for v := range w.in {
		t0 := time.Now()

		hostname := v[0]
		startTime := v[1]
		endTime := v[2]

		_, err := w.db.Query(queryCpuUsage, hostname, startTime, endTime)
		if err != nil {
			fmt.Println("error ", err)
			w.errCh <- err
			continue
		}

		fmt.Printf("running for host %s - instance %d\n", v[0], w.instance)
		t1 := time.Now()

		w.out <- t1.Sub(t0)
	}
	w.group.Done()
}
