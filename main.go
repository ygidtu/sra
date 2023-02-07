package main

import (
	"encoding/json"
	"github.com/voxelbrain/goptions"
	"github.com/ygidtu/sra/details"
	"github.com/ygidtu/sra/ebi"
	"github.com/ygidtu/sra/ena"
	"github.com/ygidtu/sra/search"
	"github.com/ygidtu/sra/study"
	"go.uber.org/zap"
)

var (
	sugar      *zap.SugaredLogger
	buildStamp = "dev"
	gitHash    = "dev"
	goVersion  = "dev"
	version    = "dev"
)

func init() {
	SetLogger(false)
}

type Params struct {
	Version bool          `goptions:"-v, --version, description='Show version'"`
	Debug   bool          `goptions:"--debug, description='Show debug info'"`
	Help    goptions.Help `goptions:"-h, --help, description='Show this help'"`

	goptions.Verbs
	Details details.Params `goptions:"detail"`
	Ebi     ebi.Params     `goptions:"ebi"`
	Search  search.Params  `goptions:"search"`
	Study   study.Params   `goptions:"study"`
	Ena     ena.Params     `goptions:"ena"`
}

func DefaultParams() *Params {
	return &Params{
		Details: details.DefaultParam(),
		Ebi:     ebi.DefaultParam(),
		Search:  search.DefaultParam(),
		Study:   study.DefaultParam(),
		Ena:     ena.DefaultParam(),
	}
}

func main() {
	options := DefaultParams()

	goptions.ParseAndFail(options)
	SetLogger(options.Debug)

	data, _ := json.MarshalIndent(options, "", "    ")
	sugar.Debug(string(data))

	if options.Version {
		sugar.Infof("Current version: %s", version)
		sugar.Infof("Git Commit Hash: %s", gitHash)
		sugar.Infof("UTC Build Time : %s", buildStamp)
		sugar.Infof("Golang Version : %s", goVersion)
	} else if options.Verbs == "search" {
		search.Search(&options.Search, sugar)
	} else if options.Verbs == "detail" {
		details.Details(&options.Details, sugar)
	} else if options.Verbs == "ebi" {
		ebi.Ebi(&options.Ebi, sugar)
	} else if options.Verbs == "study" {
		study.Study(&options.Study, sugar)
	} else if options.Verbs == "ena" {
		ena.Ena(&options.Ena, sugar)
	} else {
		goptions.PrintHelp()
	}
}
