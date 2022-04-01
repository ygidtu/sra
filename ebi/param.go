package ebi

import "encoding/json"

type Params struct {
	StudyID string `goptions:"-s, --study, description='Study ID to query'"`
	Proxy   string `goptions:"-p, --proxy, description='Proxy'"`
	Output  string `goptions:"-o, --output, description='Output json'"`
}

func (param *Params) String() string {
	str, _ := json.Marshal(param)
	return string(str)
}

func DefaultParam() Params {
	return Params{
		Output: "./ebi_output.csv",
	}
}
