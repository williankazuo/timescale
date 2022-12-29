package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/williankazuo/timescale/benchmarking"
	"github.com/williankazuo/timescale/config"
)

var (
	filePath = flag.String("filepath", "./input/query_params_sample.csv", "File path of csv file containing query params.")
	workers  = flag.Int("workers", 1, "Number of concurrent workers")
)

func main() {
	flag.Parse()

	f, err := os.Open(*filePath)
	if err != nil {
		log.Fatalf("error opening file: %s", err)
	}

	rows := make(chan []string)
	errCh := make(chan error)
	processed := make(chan bool)
	db := config.InitDB()

	csvReader := csv.NewReader(f)

	bench := benchmarking.NewBenchmarking(processed, *workers, db, rows, errCh)
	go bench.Process()

	go readRows(csvReader, rows, errCh)

	<-processed
}

func readRows(csvReader *csv.Reader, rows chan []string, errCh chan error) {
	// read header and move to next line
	_, err := csvReader.Read()
	if err != nil {
		log.Fatal("error reading first line of file")
	}

	for {
		row, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			errCh <- fmt.Errorf("error reading row from csv: %s", err)
		}

		rows <- row
	}
	close(rows)
	close(errCh)
}

// SUMMARY
// # of queries processed
// total processing across all queries
// the minimum query time (for a single query),
// the median query time,
// the average query time,
// and the maximum query time.

// should not wait to start processing queries until all input is consumed
