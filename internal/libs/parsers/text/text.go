package text

import (
	"fmt"
	"strings"
)

type Text struct{}

func New() *Text {
	return &Text{}
}

func (t *Text) Parse(attrs map[string]interface{}) ([]byte, error) {
	all := make([]string, 0, len(attrs))
	for k, v := range attrs {
		all = append(all, fmt.Sprintf("%s: %v", k, v))
	}
	line := strings.Join(all, ", ")
	return []byte(line), nil
}
