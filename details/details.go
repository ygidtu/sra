package details

import (
	"github.com/headzoo/surf/browser"
	"github.com/voxelbrain/goptions"
	"github.com/ygidtu/sra/client"
	"go.uber.org/zap"
)

const (
	SRA = "https://www.ncbi.nlm.nih.gov/sra"
	RUN = "https://trace.ncbi.nlm.nih.gov/Traces/sra"
)

var (
	surf  *browser.Browser
	sugar *zap.SugaredLogger
)

func Details(options *Params, sugar_ *zap.SugaredLogger) {
	sugar = sugar_
	sugar.Info(options.String())

	if cli, err := client.SetSurfClient(options.Proxy); err != nil {
		surf = cli
	} else {
		sugar.Fatal(err)
	}

	var obj Ncbi
	if options.StudyID != "" {
		study := &Study{ID: options.StudyID}

		err := study.Get()
		if err != nil {
			sugar.Fatal(err)
		}

		err = study.GetFiles()
		if err != nil {
			sugar.Fatal(err)
		}

		study.GetAccessions(options.Threads)
		for _, accn := range study.Accessions {
			accn.GetRun(options.Threads)
		}
		obj = study
	} else if options.AccessionID != "" {
		accn := &Accession{ID: options.AccessionID}

		err := accn.Get()
		if err != nil {
			sugar.Fatal(err)
		}

		accn.GetRun(options.Threads)

		obj = accn
	} else if options.RunID != "" {
		run := &Run{ID: options.RunID}

		err := run.Get()
		if err != nil {
			sugar.Fatal(err)
		}

		obj = run
	}

	if obj != nil {
		err := obj.Save(options.Output)
		if err != nil {
			sugar.Fatal(err)
		}
	} else {
		goptions.PrintHelp()
	}
}
