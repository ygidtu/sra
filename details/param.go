package details

import (
	"encoding/json"
)

type Params struct {
	StudyID     string `goptions:"-s, --study, description='Study ID to query'"`
	AccessionID string `goptions:"-a, --accession, description='Accession ID to query'"`
	RunID       string `goptions:"-r, --run, description='Run ID to query'"`
	Proxy       string `goptions:"-p, --proxy, description='Proxy'"`
	Output      string `goptions:"-o, --output, description='Output json'"`
	Threads     int    `goptions:"-t, --threads, description='How many threads to use'"`
}

func (param *Params) String() string {
	str, _ := json.Marshal(param)
	return string(str)
}

func DefaultParam() Params {
	return Params{
		Output:  "./output.json",
		Threads: 1,
	}
}
