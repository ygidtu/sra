package details

import (
	"encoding/json"
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type Study struct {
	ID         string
	Accessions []*Accession
	Files      []*File
}

func (study *Study) Get() error {
	sugar.Infof("Get accessions from %s", study.ID)

	sra, _ := url.Parse(SRA)

	// Use the Query() method to get the query string params as a url.Values map.
	values := sra.Query()
	values.Add("term", study.ID)
	sra.RawQuery = values.Encode()

	// get content
	err := surf.Open(sra.String())
	if err != nil {
		return err
	}

	surf.Dom().Find("div[class='rslt']").Each(func(i int, selection *goquery.Selection) {
		accn := &Accession{
			ID:  strings.ReplaceAll(filepath.Base(selection.Find("a").AttrOr("href", "")), "[accn]", ""),
			Tag: strings.ReplaceAll(selection.Find("a").Text(), "\n", ""),
		}
		study.Accessions = append(study.Accessions, accn)
	})
	return nil
}

func (study *Study) GetFiles() error {
	sugar.Infof("Get files of %s from run", study.ID)

	sra, _ := url.Parse(RUN)
	values := sra.Query()
	values.Add("study", study.ID)
	sra.RawQuery = values.Encode()

	// get content
	err := surf.Open(sra.String())
	if err != nil {
		return err
	}

	surf.Dom().Find("#id-related-files > tbody > tr").Each(func(i int, selection *goquery.Selection) {
		file := &File{
			ID:   strings.TrimSpace(selection.Find("td:nth-child(2)").Text()),
			URL:  selection.Find("td:nth-child(3) > a").AttrOr("href", ""),
			Size: strings.TrimSpace(selection.Find("td:nth-child(4)").Text()),
		}
		study.Files = append(study.Files, file)
	})

	return nil
}

func (study *Study) GetAccessions(threads int) {
	var wg sync.WaitGroup
	input := make(chan Ncbi)

	for i := 0; i < threads; i++ {
		wg.Add(1)
		go get(input, &wg)
	}

	for _, accn := range study.Accessions {
		input <- accn
	}
	close(input)
	wg.Wait()
}

func (study *Study) Json() string {
	b, _ := json.MarshalIndent(study, "", "\t")
	return string(b)
}

func (study *Study) Save(path string) error {
	return ioutil.WriteFile(path, []byte(study.Json()), os.ModePerm)
}
