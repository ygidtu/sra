package study

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
)

func queryURL(url_ string, queries map[string]string) (string, error) {
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

func open(url string) (*goquery.Document, error) {
	var html string
	if err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.Sleep(2*time.Second),
		chromedp.WaitReady("#ph-rs-table>table", chromedp.ByID),
		chromedp.OuterHTML(`document.querySelector("body")`, &html, chromedp.ByJSPath),
	); err != nil {
		sugar.Errorf("%s: %v", url, err)
	}

	return goquery.NewDocumentFromReader(strings.NewReader(html))
}

func getPage(params map[string]string) (int, *goquery.Document, error) {
	page := 1

	params["o"] = "experiment_s:a;acc_s:a"
	u, err := queryURL(URL, params)
	if err != nil {
		return page, nil, err
	}

	doc, err := open(u)
	if err != nil {
		return page, doc, err
	}

	if doc.Find("#t-pager-pages").Text() != "" {
		i, err := strconv.Atoi(doc.Find("#t-pager-pages").Text())
		if err != nil {
			return page, doc, err
		}
		page = i
	}
	return page, doc, nil
}

func extractData(doc *goquery.Document, output chan [][]string) {
	var data [][]string
	doc.Find("#ph-rs-table>table").Each(func(i int, selection *goquery.Selection) {
		row := make([]string, 0, 0)
		selection.Find("thead > tr > th[xxclass]").Each(func(i int, sel *goquery.Selection) {
			sel.Find("div").Each(func(i int, s *goquery.Selection) {
				if _, ok := s.Attr("class"); !ok {
					row = append(row, s.Text())
				}
			})
		})
		data = append(data, row)

		selection.Find("tbody > tr").Each(func(i int, sel *goquery.Selection) {
			row = make([]string, 0, 0)

			sel.Find("td[class]").Each(func(i int, s *goquery.Selection) {
				row = append(row, s.Text())
			})
			data = append(data, row)
		})
	})

	output <- data
	sugar.Debugf("return %d data", len(data))
}

func getResults(wg *sync.WaitGroup, params chan map[string]string, output chan [][]string) {
	defer wg.Done()
	for {
		param, ok := <-params

		if !ok {
			break
		}

		sugar.Debugf("%v", param)

		u, err := queryURL(URL, param)
		if err != nil {
			sugar.Fatal(err)
		}

		doc, err := open(u)
		if err != nil {
			sugar.Fatal(err)
		}
		extractData(doc, output)
	}
}
