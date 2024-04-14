package pipeline

import (
	"encoding/csv"
	"fmt"
	"io"
	"sync"
)

func workerPoolPipelineStage[IN, OUT any](input <-chan IN, output chan<- OUT, process func(IN) OUT, numWorkers int) {
	defer close(output)

	wg := sync.WaitGroup{}
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for data := range input {
				output <- process(data)
			}
		}()
	}
	wg.Wait()
}

func asynchronousPipeline2Workers(input *csv.Reader) {
	fmt.Println("--Asynchronous pipeline with worker pool----")

	parseInputCh := make(chan []string)
	convertInputCh := make(chan Record)
	encodeInputCh := make(chan Record)

	outputCh := make(chan []byte)
	done := make(chan struct{})

	numWorkers := 2
	go workerPoolPipelineStage(parseInputCh, convertInputCh, parse, numWorkers)
	go workerPoolPipelineStage(convertInputCh, encodeInputCh, convert, numWorkers)
	go workerPoolPipelineStage(encodeInputCh, outputCh, encode, numWorkers)

	go func() {
		for data := range outputCh {
			fmt.Println(string(data))
		}
		close(done)
	}()

	_, _ = input.Read()
	for {
		rec, err := input.Read()
		if err == io.EOF {
			close(parseInputCh)
			break
		}
		if err != nil {
			panic(err)
		}
		parseInputCh <- rec
	}
	<-done
}
