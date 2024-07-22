package search

import (
	"encoding/json"
	"time"
)

type Params struct {
	RBP     string        `goptions:"-i, --input, description='The list of RBP in csv format with header'"` // obligatory,
	SRA     string        `goptions:"-u, --url, description='The official api of SRA'"`
	Proxy   string        `goptions:"-x, --proxy, description='The proxy url, eg: http://127.0.0.1:7890'"`
	Output  string        `goptions:"-o, --output, description='The output directory'"`
	Param   string        `goptions:"-p, --param, description='The extra parameters'"`
	Timeout time.Duration `goptions:"-t, --timeout, description='Connection timeout in seconds'"`
	Exec    string        `goptions:"-e, --exec, description='path to chrome executable'"`
	Open    bool          `goptions:"--open, description='whether open the GUI of chrome'"`
}

func (param *Params) String() string {
	str, _ := json.Marshal(param)
	return string(str)
}

func DefaultParam() Params {
	return Params{
		SRA:     "https://www.ncbi.nlm.nih.gov/sra/",
		Output:  "./output",
		Timeout: 10 * time.Second,
		Param:   `"Homo sapiens"[orgn:__txid9606] AND(rna seq[Strategy])`,
	}
}
