package benchmark

import (
	"bytes"
	"fmt"
	"text/template"
	"time"
)

type t struct {
	*template.Template
}

func NewTemplate() *t {
	return &t{}
}

func (ø *t) Parse(s string) (err error) {
	ø.Template, err = template.New(fmt.Sprintf("t-%v", time.Now().UnixNano())).Parse(s)
	if err != nil {
		panic(err.Error())
	}
	return
}

func (ø *t) Replace(data map[string]string, buf *bytes.Buffer) error {
	return ø.Template.Execute(buf, data)
}
