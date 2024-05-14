package internal

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
)

var httpClient = &http.Client{}

type stressTestRunner struct {
	name         string
	url          string
	requestsToDo int64
	requestsDone int64
	report       *testReport
}

func NewRunner(name string, url string, requests int64, report *testReport) *stressTestRunner {
	runner := &stressTestRunner{}
	runner.name = name
	runner.url = url
	runner.requestsToDo = requests
	runner.report = report
	return runner
}

func (r *stressTestRunner) Run(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	logPrefix := fmt.Sprintf("[TEST RUNNER %s]", strings.ToUpper(r.name))
	log := func(format string, a ...any) {
		logMessageWithPrefix(logPrefix, format, a...)
	}

	log("Iniciado com %d requisições", r.requestsToDo)

	for i := int64(0); i < r.requestsToDo; i++ {
		statusCode, err := r.doRequest(ctx)
		if err != nil {
			log("Erro: %s", err.Error())
			r.report.AddError()
			continue
		}

		r.report.AddRequest(statusCode)
		r.requestsDone++
	}

	log("Finalizado (%d/%d)", r.requestsDone, r.requestsToDo)
}

func (r *stressTestRunner) doRequest(ctx context.Context) (int, error) {
	request, err := http.NewRequest("GET", r.url, nil)
	if err != nil {
		return 0, err
	}
	request.Close = true
	request = request.WithContext(ctx)

	response, err := httpClient.Do(request)
	if err != nil {
		return 0, err
	}
	io.Copy(io.Discard, response.Body)
	defer response.Body.Close()

	return response.StatusCode, nil
}
