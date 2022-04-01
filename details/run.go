package details

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"net/url"
	"os"
	"strings"
)

type Run struct {
	ID    string
	Files []*File
}

func (run *Run) Get() error {
	sugar.Infof("Get files of %s from run", run.ID)

	sra, _ := url.Parse(RUN)
	values := sra.Query()
	values.Add("run", run.ID)
	sra.RawQuery = values.Encode()

	// get content
	err := surf.Open(sra.String())
	if err != nil {
		return err
	}

	run.Files = make([]*File, 0, 0)
	surf.Dom().Find("div[class='section menu-content menu1 download readable'] > div").Each(func(i int, selection *goquery.Selection) {
		type_, size := "", ""

		selection.Find("table > tbody > tr").Each(func(i int, selection *goquery.Selection) {
			locationIdx, urlIdx := 1, 2
			if type_ == "" {
				type_ = strings.TrimSpace(selection.Find("td:nth-child(1)").Text())
				locationIdx, urlIdx = 3, 4
			}
			if size == "" {
				size = strings.TrimSpace(selection.Find("td:nth-child(2)").Text())
			}

			file := &File{
				ID:       run.ID,
				Type:     type_,
				Size:     size,
				Location: strings.TrimSpace(selection.Find(fmt.Sprintf("td:nth-child(%d)", locationIdx)).Text()),
				URL:      strings.TrimSpace(selection.Find(fmt.Sprintf("td:nth-child(%d)", urlIdx)).Text()),
			}

			run.Files = append(run.Files, file)
		})
	})

	return nil
}

func (run *Run) Json() string {
	b, _ := json.MarshalIndent(run, "", "\t")
	return string(b)
}

func (run *Run) Save(path string) error {
	return ioutil.WriteFile(path, []byte(run.Json()), os.ModePerm)
}
