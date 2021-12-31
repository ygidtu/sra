package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type Progress struct {
	Path string
	Data map[string]int
}

// Load progress from json file
func (p *Progress) Load() {
	if _, err := os.Stat(p.Path); !os.IsNotExist(err) {
		data, err := ioutil.ReadFile(p.Path)
		if err != nil {
			log.Warn("failed to reload progress from: ", p.Path)
		}

		if err := json.Unmarshal(data, &p.Data); err != nil {
			log.Warn("failed to unmarshal progress")
		}

	}
}

// Dump progress to json file
func (p *Progress) Dump() error {
	dataStr, _ := json.Marshal(p.Data)
	return ioutil.WriteFile(p.Path, dataStr, 0755)
}

// Has check whether this rbp is processed
func (p *Progress) Has(rbp string) bool {
	_, ok := p.Data[rbp]
	return ok
}

// Add add new progressed rbp
func (p *Progress) Add(rbp string) {
	p.Data[rbp] = 0
}
