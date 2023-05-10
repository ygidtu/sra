package enaP

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/schollz/progressbar/v3"
	"github.com/ygidtu/sra/client"
	"go.uber.org/zap"
	"html"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
)

type summary struct {
	Accession          string   `json:"accession"`
	SecondaryAccession []string `json:"secondaryAccession"`
	Title              string   `json:"title"`
	Description        string   `json:"description"`
	Taxon              int      `json:"taxon"`
	DataType           string   `json:"dataType"`
	Status             int      `json:"status"`
	StatusDescription  string   `json:"statusDescription"`
}

func (s summary) GetHeaders() []string {
	res := make([]string, 0)
	val := reflect.ValueOf(s)
	for i := 0; i < val.Type().NumField(); i++ {
		res = append(res, val.Type().Field(i).Tag.Get("json"))
	}
	return res
}

func (s summary) GetValues() []string {
	res := make([]string, 0)
	val := reflect.ValueOf(s)
	for i := 0; i < val.Type().NumField(); i++ {
		if val.Field(i).CanInt() {
			res = append(res, fmt.Sprintf("%v", val.Field(i)))
		} else if val.Field(i).Type() == reflect.TypeOf(s.SecondaryAccession) {
			res = append(res, strings.Join(s.SecondaryAccession, ","))
		} else {
			res = append(res, val.Field(i).String())
		}
	}
	return res
}

type summaries struct {
	Summaries []summary `json:"summaries"`
	Total     string    `json:"total"`
}

func parseQuery(urlA *url.URL, key string) {
	// Use the Query() method to get the query string params as an url.Values map.
	values := urlA.Query()

	// Make the changes that you want using the Add(), Set() and Del() methods. If
	// you want to retrieve or check for a specific parameter you can use the Get()
	// and Has() methods respectively.
	urlA.Path = filepath.Join(urlA.Path, key)
	values.Add("offset", "0")
	values.Add("limit", "1000")
	values.Add("format", "json")
	urlA.RawQuery = values.Encode()
}

func parseKey(param *Params) ([]string, error) {
	keys := make([]string, 0, 0)
	if _, err := os.Stat(param.Key); os.IsNotExist(err) {
		keys = append(keys, strings.TrimSpace(param.Key))
	} else {
		f, err := os.Open(param.Key)
		if err != nil {
			return keys, err
		}
		defer f.Close()
		r := bufio.NewScanner(f)

		for r.Scan() {
			keys = append(keys, strings.TrimSpace(r.Text()))
		}
	}

	return keys, nil
}

func write(outfile string, oc chan *summaries, sugar *zap.SugaredLogger, wg *sync.WaitGroup, bar *progressbar.ProgressBar) {
	defer wg.Done()
	f, err := os.OpenFile(outfile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		sugar.Fatal(err)
	}
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()
	w.UseCRLF = true

	header := true
	for !bar.IsFinished() {
		data, ok := <-oc
		if !ok {
			break
		}

		if data != nil {
			for _, s := range data.Summaries {
				if header {
					_ = w.Write(s.GetHeaders())
					header = false
				}

				_ = w.Write(s.GetValues())
			}
		}

		_ = bar.Add(1)
	}
}

func EnaP(options *Params, sugar *zap.SugaredLogger) {
	sugar.Debug(options.String())

	if _, err := os.Stat(filepath.Dir(options.Output)); os.IsNotExist(err) {
		err := os.MkdirAll(filepath.Dir(options.Output), os.ModePerm)
		if err != nil {
			sugar.Fatal(err)
		}
	}

	urlA, err := url.Parse(options.ENA)
	if err != nil {
		sugar.Fatalf("failed to parse ENA url: %v", err)
	}

	keys, err := parseKey(options)
	if err != nil {
		sugar.Fatal(err)
	}

	// create and start new bar
	bar := client.Bar(len(keys))

	var wg sync.WaitGroup
	paramChan := make(chan string)
	outChan := make(chan *summaries)

	wg.Add(1)
	go write(options.Output, outChan, sugar, &wg, bar)

	for i := 0; i < options.Threads; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				key, ok := <-paramChan
				if !ok {
					break
				}
				if key == "" {
					continue
				}

				urlB := *urlA
				parseQuery(&urlB, strings.TrimSpace(key))
				sugar.Debug(urlB.String())

				if cli, err := client.SetSurfClient(options.Proxy); err == nil {
					if err := cli.Open(urlB.String()); err != nil {
						sugar.Warn(err)
						continue
					}
					var summ *summaries
					err = json.Unmarshal([]byte(html.UnescapeString(cli.Body())), &summ)
					if err != nil {
						sugar.Debug("%v", html.UnescapeString(cli.Body()))
						sugar.Errorf("%v: %v", key, err)
					}
					outChan <- summ
				} else {
					sugar.Fatal(err)
				}
			}
		}()
	}

	for _, k := range keys {
		paramChan <- k
	}
	close(paramChan)
	wg.Wait()
}
