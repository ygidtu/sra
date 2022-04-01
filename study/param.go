package study

import (
	"encoding/json"
)

type Params struct {
	StudyID string `goptions:"-s, --study, description='Study ID to query'"`
	Proxy   string `goptions:"-p, --proxy, description='Proxy'"`
	Output  string `goptions:"-o, --output, description='Output json'"`
	Open    bool   `goptions:"--open, description='是否打开chrome的图形化界面'"`
	Threads int    `goptions:"-t, --threads, description='How many threads to use'"`
}

func (param *Params) String() string {
	str, _ := json.Marshal(param)
	return string(str)
}

func DefaultParam() Params {
	return Params{
		Output:  "./output.csv",
		Threads: 1,
	}
}