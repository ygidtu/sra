package details

import "encoding/json"

type File struct {
	ID       string
	URL      string
	Location string
	Size     string
	Type     string
}

func (file *File) Json() string {
	b, _ := json.MarshalIndent(file, "", "\t")
	return string(b)
}
