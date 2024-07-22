package study

import (
	"encoding/json"
)

type Params struct {
	StudyID string `goptions:"-s, --study, description='Study ID to query'"`
	Proxy   string `goptions:"-p, --proxy, description='Proxy'"`
	Output  string `goptions:"-o, --output, description='Output json'"`
	Open    bool   `goptions:"--open, description='whether open the GUI of chrome'"`
	Threads int    `goptions:"-t, --threads, description='How many threads to use'"`
	Exec    string `goptions:"-e, --exec, description='path to chrome executable'"`
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
