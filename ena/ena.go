package ena

import (
	"bufio"
	"compress/gzip"
	"github.com/cheggaaa/pb/v3"
	"github.com/headzoo/surf/browser"
	"github.com/ygidtu/sra/client"
	"go.uber.org/zap"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var (
	surf *browser.Browser
)

func parseQuery(urlA *url.URL, key, fields string) {
	// Use the Query() method to get the query string params as an url.Values map.
	values := urlA.Query()

	// Make the changes that you want using the Add(), Set() and Del() methods. If
	// you want to retrieve or check for a specific parameter you can use the Get()
	// and Has() methods respectively.
	values.Add("fields", fields)
	values.Add("result", "read_run")
	values.Add("accession", key)
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
		r := bufio.NewReader(f)

		for {
			line, err := r.ReadString('\n')
			if err != nil {
				break
			}
			keys = append(keys, strings.TrimSpace(line))
		}
	}

	return keys, nil
}

func Ena(options *Params, sugar *zap.SugaredLogger) {
	sugar.Info(options.String())
	if cli, err := client.SetSurfClient(options.Proxy); err == nil {
		surf = cli
	} else {
		sugar.Fatal(err)
	}

	if _, err := os.Stat(options.Output); os.IsNotExist(err) {
		err := os.MkdirAll(options.Output, os.ModePerm)
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
	bar := pb.StartNew(len(keys))

	paramChan := make(chan string)

	var wg sync.WaitGroup
	for i := 0; i < options.Threads; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				key, ok := <-paramChan
				if !ok {
					break
				}
				bar.Increment()

				if key == "" {
					continue
				}

				oFile := filepath.Join(options.Output, key)
				if options.Resume {
					if _, err := os.Stat(oFile); !os.IsNotExist(err) {
						continue
					}
				}

				urlB := *urlA
				parseQuery(&urlB, strings.TrimSpace(key), options.Fields)
				sugar.Debug(urlB.String())
				if err := surf.Open(urlB.String()); err != nil {
					sugar.Fatal(err)
				}

				f, err := os.OpenFile(oFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
				if err != nil {
					sugar.Fatal(err)
				}

				gw := gzip.NewWriter(f)
				_, _ = surf.Download(gw)
				_ = gw.Close()
				_ = f.Close()
			}
		}()
	}

	for _, k := range keys {
		paramChan <- k
	}
	close(paramChan)
	wg.Wait()
}
