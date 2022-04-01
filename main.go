package main

import (
	"encoding/json"
	"github.com/voxelbrain/goptions"
	"github.com/ygidtu/sra/details"
	"github.com/ygidtu/sra/ebi"
	"github.com/ygidtu/sra/search"
	"github.com/ygidtu/sra/study"
	"go.uber.org/zap"
)

var (
	sugar *zap.SugaredLogger
)

func init() {
	SetLogger(false)
}

type Params struct {
	Version bool          `goptions:"-v, --version, description='Show version'"`
	Debug   bool          `goptions:"--debug, description='Show debug info'"`
	Help    goptions.Help `goptions:"--help, description='Show this help'"`

	goptions.Verbs
	Details details.Params `goptions:"detail"`
	Ebi     ebi.Params     `goptions:"ebi"`
	Search  search.Params  `goptions:"search"`
	Study   study.Params   `goptions:"study"`
}

func DefaultParams() *Params {
	return &Params{
		Details: details.DefaultParam(),
		Ebi:     ebi.DefaultParam(),
		Search:  search.DefaultParam(),
		Study:   study.DefaultParam(),
	}
}

func main() {
	options := DefaultParams()

	goptions.ParseAndFail(options)
	SetLogger(options.Debug)

	data, _ := json.MarshalIndent(options, "", "    ")
	sugar.Debug(string(data))

	if options.Verbs == "search" {
		search.Search(&options.Search, sugar)
	} else if options.Verbs == "dtails" {
		details.Details(&options.Details, sugar)
	} else if options.Verbs == "ebi" {
		ebi.Ebi(&options.Ebi, sugar)
	} else if options.Verbs == "study" {
		study.Study(&options.Study, sugar)
	} else {
		goptions.PrintHelp()
	}
}
