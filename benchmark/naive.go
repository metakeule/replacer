package benchmark

import (
	"strings"
)

type Naive struct {
	Map      map[string]string
	Template string
}

func (ø *Naive) Replace() (s string) {
	s = ø.Template
	for k, v := range ø.Map {
		s = strings.Replace(s, k, v, -1)
	}
	return
}
