package benchmark

import (
	"regexp"
)

type Regexp struct {
	Regexp   *regexp.Regexp
	Map      map[string]string
	Replacer func(string) string
	Template string
}

func (ø *Regexp) Setup() {
	ø.Replacer = func(found string) string { return ø.Map[found] }
}

func (ø *Regexp) Replace() string {
	return ø.Regexp.ReplaceAllStringFunc(ø.Template, ø.Replacer)
}
