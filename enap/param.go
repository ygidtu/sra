package enaP

import (
	"encoding/json"
)

type Params struct {
	Key     string `goptions:"-i, --input, description='KeyID or list of KeyID'"` // obligatory,
	ENA     string `goptions:"-u, --url, description='The official api of ENA'"`
	Proxy   string `goptions:"-x, --proxy, description='The proxy url, eg: http://127.0.0.1:7890'"`
	Output  string `goptions:"-o, --output, description='The output directory'"`
	Threads int    `goptions:"-t, --threads, description='The number of threads to use, do not use too much threads'"`
	Resume  bool   `goptions:"-c, --resume, description='Skip finished requests'"`
}

func (param *Params) String() string {
	str, _ := json.MarshalIndent(param, "", "    ")
	return string(str)
}

func DefaultParam() Params {
	return Params{
		ENA:     "https://www.ebi.ac.uk/ena/browser/api/summary/",
		Output:  "./output.csv",
		Threads: 1,
	}
}
