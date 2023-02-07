package ebi

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/headzoo/surf/browser"
	"github.com/ygidtu/sra/client"
	"go.uber.org/zap"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

const (
	URL      = "https://www.ebi.ac.uk/arrayexpress/rss/v2/experiments/"
	FILE_URL = "https://www.ebi.ac.uk/arrayexpress/files"
)

var (
	surf *browser.Browser
)

func queryURL(url_ string, queries map[string]string) (string, error) {
	// Use url.Parse() to parse a string into a *url.URL type. If your URL is
	// already an url.URL type you can skip this step.
	urlA, err := url.Parse(url_)
	if err != nil {
		return url_, err
	}

	// Use the Query() method to get the query string params as a url.Values map.
	values := urlA.Query()

	for key, value := range queries {
		values.Add(key, value)
	}

	urlA.RawQuery = values.Encode()
	return urlA.String(), nil
}

func pathURL(url_ string, paths ...string) (string, error) {
	// Use url.Parse() to parse a string into a *url.URL type. If your URL is
	// already a url.URL type you can skip this step.
	urlA, err := url.Parse(url_)
	if err != nil {
		return url_, err
	}

	for _, path := range paths {
		urlA.Path = filepath.Join(urlA.Path, path)
	}

	return urlA.String(), nil
}

func Ebi(options *Params, sugar *zap.SugaredLogger) {
	sugar.Info(options.String())
	if cli, err := client.SetSurfClient(options.Proxy); err == nil {
		surf = cli
	} else {
		sugar.Fatal(err)
	}

	u, err := queryURL(URL, map[string]string{"query": options.StudyID})
	if err != nil {
		sugar.Fatal(err)
	}

	sugar.Debug(u)
	if err := surf.Open(u); err != nil {
		sugar.Fatal(err)
	}

	surf.Dom().Find("guid").Each(func(i int, selection *goquery.Selection) {
		sampleId := strings.Trim(filepath.Base(selection.Text()), "/")

		sugar.Debugf("SampleID is %s", sampleId)

		if sampleId != "" {
			u, err := pathURL(FILE_URL, sampleId, fmt.Sprintf("%s.sdrf.txt", sampleId))
			if err != nil {
				sugar.Fatal(err)
			}
			if err := surf.Open(u); err != nil {
				sugar.Fatal(err)
			}

			f, err := os.OpenFile(options.Output, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.ModePerm)
			defer f.Close()

			if err != nil {
				sugar.Fatal(err)
			}
			_, err = surf.Download(f)
			if err != nil {
				sugar.Fatal(err)
			}
		}
	})

}
