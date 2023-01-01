package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"time"

	"github.com/williankazuo/timescale/benchmarking"
	"github.com/williankazuo/timescale/config"
)

var (
	filePath = flag.String("filepath", "./input/query_params.csv", "File path of csv file containing query params.")
	workers  = flag.Int("workers", 1, "Number of concurrent workers")
)

func main() {
	flag.Parse()

	t0 := time.Now()

	f, err := os.Open(*filePath)
	if err != nil {
		log.Fatalf("error opening file: %s", err)
	}

	var (
		db            = config.InitDB()
		errCh         = make(chan error)
		out           = make(chan time.Duration)
		wgWorkers     = sync.WaitGroup{}
		poolOfWorkers = []benchmarking.Worker{}
		done          = make(chan bool)
	)

	wgWorkers.Add(*workers)

	// instantiating workers and start processing them.
	for i := 0; i < *workers; i++ {
		in := make(chan []string)
		w := benchmarking.NewWorker(i+1, db, in, out, errCh, &wgWorkers)
		poolOfWorkers = append(poolOfWorkers, *w)
		go w.Process()
	}

	csvReader := csv.NewReader(f)
	bench := benchmarking.NewBenchmarking(out, errCh, &wgWorkers, done)
	go bench.CollectResults()
	go bench.CollectErrors()
	go readRows(csvReader, errCh, poolOfWorkers)

	<-done

	fmt.Printf("Time spent running the tool: %s\n", time.Now().Sub(t0))
}

func readRows(csvReader *csv.Reader, errCh chan<- error, poolOfWorkers []benchmarking.Worker) {
	// read header and move to next line
	_, err := csvReader.Read()
	if err != nil {
		log.Fatal("error reading first line of file")
	}

	mapWorkers := make(map[string]benchmarking.Worker)
	count := 0
	for {
		row, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			errCh <- fmt.Errorf("error reading row from csv: %s", err)
		}

		// chose the worker based on host
		host := row[0]
		chosenWorker, ok := mapWorkers[host]
		if !ok {
			if count > len(poolOfWorkers)-1 {
				count = 0
			}

			chosenWorker = poolOfWorkers[count]
			mapWorkers[host] = chosenWorker
			count++
		}
		chosenWorker.In <- row
	}

	// close all workers
	for _, w := range poolOfWorkers {
		close(w.In)
	}
	close(errCh)
}
