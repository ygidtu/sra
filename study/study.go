package study

import (
	"context"
	"encoding/csv"
	"fmt"
	"github.com/schollz/progressbar/v3"
	"github.com/ygidtu/sra/client"
	"go.uber.org/zap"
	"os"
	"sync"
)

const (
	URL = "https://trace.ncbi.nlm.nih.gov/Traces/study/"
)

var (
	sugar *zap.SugaredLogger
	ctx   context.Context
)

func write(path string, output chan [][]string, bar *progressbar.ProgressBar, wg *sync.WaitGroup) {
	defer wg.Done()
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

		_ = bar.Add(1)

		if len(data) > 0 {
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
				writer.Flush()
			}
		}

		if bar.IsFinished() {
			break
		}
	}
}

func Study(options *Params, sugar_ *zap.SugaredLogger) {
	sugar = sugar_

	sugar.Info(options.String())

	// create context
	ctx_, cancel := client.SetChromeClient(options.Open, options.Proxy, options.Exec, sugar)
	defer cancel()
	ctx = ctx_

	page, doc, err := getPage(map[string]string{"acc": options.StudyID})
	if err != nil {
		sugar.Fatal(err)
	}
	sugar.Infof("Current study has %d pages", page)

	if page < options.Threads {
		options.Threads = page
	}

	var wg sync.WaitGroup
	params := make(chan map[string]string)
	output := make(chan [][]string)

	for i := 0; i < options.Threads; i++ {
		wg.Add(1)

		go getResults(&wg, params, output)
	}

	bar := progressbar.Default(int64(page))
	wg.Add(1)
	go write(options.Output, output, bar, &wg)

	extractData(doc, output)
	for i := 2; i <= page; i++ {
		params <- map[string]string{
			"page": fmt.Sprintf("%d", i),
			"acc":  options.StudyID,
			"o":    "experiment_s:a;acc_s:a",
		}
	}

	close(params)
	wg.Wait()
	close(output)
}
