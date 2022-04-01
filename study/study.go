package study

import (
	"context"
	"encoding/csv"
	"fmt"
	"github.com/ygidtu/sra/client"
	"go.uber.org/zap"
	"os"
	"sync"
)

const (
	URL = "https://trace.ncbi.nlm.nih.gov/Traces/study/?acc=SRP344246&o=experiment_s:a;acc_s:a"
)

var (
	sugar *zap.SugaredLogger
	ctx   context.Context
)

func write(path string, output chan [][]string) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.ModePerm)
	if err != nil {
		sugar.Fatal(err)
	}
	defer f.Close()

	writer := csv.NewWriter(f)
	defer writer.Flush()

	header := true
	for {
		data, ok := <-output
		if !ok {
			break
		}

		if !header {
			data = data[1:]
		} else {
			header = false
		}

		for _, value := range data {
			err := writer.Write(value)
			if err != nil {
				sugar.Fatal(err)
			}
		}
	}
}

func Study(options *Params, sugar_ *zap.SugaredLogger) {
	sugar = sugar_

	sugar.Info(options.String())

	// create context
	ctx_, cancel := client.SetChromeClient(options.Open, options.Proxy, sugar)
	defer cancel()
	ctx = ctx_

	page, err := getPage(map[string]string{"acc": options.StudyID})
	if err != nil {
		sugar.Fatal(err)
	}
	sugar.Debugf("Current study has %d pages", page)

	var wg sync.WaitGroup
	params := make(chan map[string]string)
	output := make(chan [][]string)

	for i := 0; i < options.Threads; i++ {
		wg.Add(1)

		go getResults(&wg, params, output)
	}

	go write(options.Output, output)

	for i := 1; i <= page; i++ {
		params <- map[string]string{
			"page": fmt.Sprintf("%d", i),
			"acc":  options.StudyID,
			"o":    "experiment_s:a;acc_s:a",
		}
	}
	close(params)
	wg.Wait()

}
