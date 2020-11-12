package reader

// a slight modification of the repo : https://github.com/sourcegraph/lsif-protocol/blob/master/reader/reader.go

import (
	"encoding/csv"
	"io"
	"runtime"
	"sync"

	"github.com/ayushsherpa111/goPassd/schema"
)

var numMaxProcs = runtime.GOMAXPROCS(0)
var pool = sync.Pool{New: func() interface{} { return []string{} }}

const chanBuffer = 16

func ReadCSVFile(r io.Reader, marshaller func([]string, []string) schema.Item) chan schema.Item {
	csvReader := csv.NewReader(r)
	// Channel to read data from
	var fileChan = make(chan []string, chanBuffer)
	// go routine that reads stuff from the csvReader
	var header []string
	go func() {
		var err error
		defer close(fileChan)
		for {
			poolSlice := pool.Get().([]string)
			poolSlice, err = csvReader.Read()

			if err != nil {
				break
			}
			// Skip header from csv
			if header == nil {
				header = poolSlice
				continue
			}
			fileChan <- poolSlice
		}
	}()

	var readChan = make(chan schema.Item, chanBuffer)
	// the worker manager goRoutine
	go func() {
		// worker channels that take data from fileChannel and perform task
		var workerChan = make(chan int, numMaxProcs)
		defer close(workerChan)
		defer close(readChan)
		// The pool in which the data sent from the reading channel is stored in
		var workerDataPool = make([][]string, numMaxProcs)
		// Collection of results
		var resultCol = make([]schema.Item, numMaxProcs)

		signal := make(chan int, numMaxProcs)
		// workers
		for i := 0; i < numMaxProcs; i++ {
			go func() {
				for idx := range workerChan {
					resultCol[idx] = marshaller(header, workerDataPool[idx])
					signal <- idx
				}
			}()
		}
		done := false
		for !done {
			// Worker Routine counter
			c := 0
			for c < numMaxProcs {
				var ok bool
				workerDataPool[c], ok = <-fileChan
				if !ok {
					done = true
					break
				}
				workerChan <- c
				c++
			}

			// free up pool
			for j := 0; j < c; j++ {
				idx := <-signal
				pool.Put(workerDataPool[idx])
				workerDataPool[idx] = nil
			}

			// send data thru the reading channel for the user to use
			for j := 0; j < c; j++ {
				readChan <- resultCol[j]
			}
		}
	}()

	return readChan
}
