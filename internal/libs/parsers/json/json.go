package json

import "encoding/json"

type JSON struct{}

func New() *JSON {
	return &JSON{}
}

func (j *JSON) Parse(attrs map[string]interface{}) ([]byte, error) {
	return json.Marshal(attrs)
}
