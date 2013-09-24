package benchmark

import (
	"bytes"
)

type Bytes struct {
	Map    map[string][]byte
	Buffer *bytes.Buffer
}

func (ø *Bytes) Parse(s string) {
	ø.Buffer = &bytes.Buffer{}
	ø.Buffer.WriteString(s)
}

func (ø *Bytes) Replace() (r []byte) {
	r = ø.Buffer.Bytes()

	for k, v := range ø.Map {
		r = bytes.Replace(r, []byte(k), v, -1)
	}
	return
}
