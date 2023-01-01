package benchmarking

import (
	"fmt"
	"sort"
	"sync"
	"time"
)

type Benchmarking struct {
	out       chan time.Duration
	errCh     <-chan error
	wgWorkers *sync.WaitGroup
	done      chan<- bool
}

func NewBenchmarking(out chan time.Duration, errCh <-chan error, wgWorkers *sync.WaitGroup, done chan<- bool) *Benchmarking {
	return &Benchmarking{
		out:       out,
		errCh:     errCh,
		wgWorkers: wgWorkers,
		done:      done,
	}
}

func (b *Benchmarking) CollectResults() {
	var (
		totalQueries    = 0
		totalProcessing time.Duration
		minQueryTime    time.Duration
		medianQueryTime time.Duration
		avgQueryTime    time.Duration
		maxQueryTime    time.Duration
		results         []time.Duration
	)

	var wgCollected sync.WaitGroup
	go func() {
		wgCollected.Add(1)
		for r := range b.out {
			totalProcessing += r
			results = append(results, r)
		}
		wgCollected.Done()
	}()

	b.wgWorkers.Wait()
	close(b.out)
	wgCollected.Wait()

	sort.Slice(results, func(i, j int) bool {
		return results[i] < results[j]
	})

	totalQueries = len(results)
	minQueryTime = results[0]
	avgQueryTime = totalProcessing / time.Duration(totalQueries)
	maxQueryTime = results[totalQueries-1]
	if totalQueries%2 == 0 {
		medianQueryTime = (results[totalQueries/2-1] + results[totalQueries/2]) / 2
	} else {
		medianQueryTime = results[totalQueries/2]
	}

	fmt.Printf("--------------\n")
	fmt.Printf("Number of queries: %d\n", totalQueries)
	fmt.Printf("Total processing time all queries: %s\n", totalProcessing)
	fmt.Printf("Minimum query time: %s\n", minQueryTime)
	fmt.Printf("Median query time: %s\n", medianQueryTime)
	fmt.Printf("Average query time: %s\n", avgQueryTime)
	fmt.Printf("Maximum query time: %s\n", maxQueryTime)
	fmt.Printf("--------------\n")

	b.done <- true
}

func (b *Benchmarking) CollectErrors() {
	for e := range b.errCh {
		fmt.Println("error ", e.Error())
	}
}
