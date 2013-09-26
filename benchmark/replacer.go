package benchmark

import (
	"strings"
)

type Naive2 struct {
	Replacements []string
	Template     string
}

func (ø *Naive2) Replace() (s string) {
	r := strings.NewReplacer(ø.Replacements...)
	return r.Replace(ø.Template)
}
