package search

import (
	"encoding/json"
	"github.com/voxelbrain/goptions"
	"time"
)

type Params struct {
	RBP     string        `goptions:"-i, --input, obligatory, description='RBP的list，csv格式，带列名'"`
	SRA     string        `goptions:"-u, --url, description='SRA的官方链接'"`
	Proxy   string        `goptions:"-x, --proxy, description='代理链接地址，比如：http://127.0.0.1:7890'"`
	Output  string        `goptions:"-o, --output, description='输出文件夹'"`
	Param   string        `goptions:"-p, --param, description='额外的查询参数'"`
	Timeout time.Duration `goptions:"-t, --timeout, description='Connection timeout in seconds'"`
	Open    bool          `goptions:"--open, description='是否打开chrome的图形化界面'"`
	Debug   bool          `goptions:"--debug, description='显示debug log'"`
	Help    goptions.Help `goptions:"-h, --help, description='Show this help'"`
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
