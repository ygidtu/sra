package details

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type Accession struct {
	ID                string
	Tag               string
	ExternalID        string
	SubmittedBy       string
	Study             string
	Abstract          string
	BioProject        string
	BioSample         string
	Sample            string
	SampleDescription string
	Organism          string
	Library           map[string]string
	Runs              []*Run
}

func (accn *Accession) Get() error {
	sugar.Infof("Get runs from %s", accn.ID)

	sra, _ := url.Parse(SRA)
	sra.Path = filepath.Join(sra.Path, fmt.Sprintf("%s[accn]", accn.ID))

	// get content
	err := surf.Open(sra.String())
	if err != nil {
		return err
	}

	accn.Library = make(map[string]string)
	accn.Runs = make([]*Run, 0, 0)

	surf.Dom().Find("#ResultView").Each(func(i int, selection *goquery.Selection) {
		selection.Find("div.sra-full-data").Each(func(i int, sel *goquery.Selection) {
			vals := strings.Split(sel.Text(), ":")

			if len(vals) > 1 {
				switch strings.TrimSpace(vals[0]) {
				case "External Id":
					accn.ExternalID = strings.TrimSpace(vals[1])
				case "Submitted by":
					accn.SubmittedBy = strings.TrimSpace(vals[1])
				case "Study":
					accn.Abstract = strings.TrimSpace(sel.Find("div.expand > div").Text())
					accn.BioProject = strings.TrimSpace(sel.Find("div[class='expand-body'] > a:nth-child(1)").Text())
					accn.Study = strings.TrimSpace(sel.Find("span").Text())
					accn.Study = strings.Split(accn.Study, accn.BioProject)[0]
				case "Sample":
					accn.SampleDescription = strings.ReplaceAll(strings.TrimSpace(sel.Find("span").Text()), " ", " ")

					texts := strings.Split(sel.Find("span > div.expand-body").Text(), "•")
					if len(texts) > 0 {
						accn.BioSample = strings.TrimSpace(texts[0])
						accn.SampleDescription = strings.Split(accn.SampleDescription, accn.BioSample)[0]
					}
					if len(texts) > 1 {
						accn.Sample = strings.TrimSpace(texts[1])
					}

					accn.Organism = strings.TrimSpace(sel.Find("div.expand-body > span").Text())
				case "Library":
					sel.Find("div.expand-body > div").Each(func(i int, s *goquery.Selection) {
						vals = strings.Split(s.Text(), ":")
						if len(vals) > 1 {
							accn.Library[strings.TrimSpace(vals[0])] = strings.TrimSpace(vals[1])
						}
					})
				default:

				}
			}
		})

		selection.Find("table > tbody > tr").Each(func(i int, selection *goquery.Selection) {
			accn.Runs = append(accn.Runs, &Run{ID: strings.TrimSpace(selection.Find("td").First().Text())})
		})
	})
	return nil
}

func (accn *Accession) GetRun(threads int) {
	input := make(chan Ncbi)
	var wg sync.WaitGroup
	for i := 0; i < threads; i++ {
		wg.Add(1)
		go get(input, &wg)
	}

	for _, run := range accn.Runs {
		input <- run
	}

	close(input)
	wg.Wait()
}

func (accn *Accession) Json() string {
	b, _ := json.MarshalIndent(accn, "", "\t")
	return string(b)
}

func (accn *Accession) Save(path string) error {
	return ioutil.WriteFile(path, []byte(accn.Json()), os.ModePerm)
}
