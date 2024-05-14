package internal

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"time"
)

const APP_NAME = "GO ESTRESSE"

func RunStressTester(url string, requests int64, concurrency int64) {

	logMessage("Executando com parametros url=\"%s\" request=%d concurrency=%d", url, requests, concurrency)

	if concurrency > requests {
		concurrency = requests
		logMessage("Reduzindo o número de workers para %d", concurrency)
	}

	loadPerRunner, err := getLoadPerRunner(requests, concurrency)
	if err != nil {
		logMessage("Erro: %s", err)
		os.Exit(1)
	}

	logMessage("Usando runner load: %v", *loadPerRunner)

	ctx := context.Background()
	report := NewTestReport(url, concurrency)
	wg := &sync.WaitGroup{}
	start := time.Now()
	printReportOnExit(report, &start)

	for i, load := range *loadPerRunner {
		runner := NewRunner(strconv.FormatInt(int64(i+1), 10), url, load, report)
		wg.Add(1)
		go runner.Run(ctx, wg)
	}

	wg.Wait()

	duration := time.Since(start)
	report.Duration = &duration
	report.PrintReport(true)
}

func printReportOnExit(report *testReport, start *time.Time) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		fmt.Print("\n\n")
		duration := time.Since(*start)
		report.Duration = &duration
		report.PrintReport(false)
		os.Exit(1)
	}()
}
