package enaP

import (
	"encoding/json"
)

type Params struct {
	Key     string `goptions:"-i, --input, description='KeyID or list of KeyID'"` // obligatory,
	ENA     string `goptions:"-u, --url, description='ENA的官方API链接'"`
	Proxy   string `goptions:"-x, --proxy, description='代理链接地址，比如：http://127.0.0.1:7890'"`
	Output  string `goptions:"-o, --output, description='输出文件夹'"`
	Threads int    `goptions:"-t, --threads, description='所使用的线程数'"`
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
