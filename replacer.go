package replacer

import (
	"bytes"
	"fmt"
	"io"
	"sort"
)

type delimiter int

const (
	DefaultDelimiter delimiter = iota
	HashDelimiter
	DollarDelimiter
	PercentDelimiter
)

var delimiterBytes = map[delimiter][]byte{
	DefaultDelimiter: []byte(`@@`),
	HashDelimiter:    []byte(`##`),
	DollarDelimiter:  []byte(`$$`),
	PercentDelimiter: []byte(`%%`),
}

type place struct {
	pos         int
	placeholder string
}

type places []place

// fullfill sort.Interface.
func (p places) Len() int           { return len(p) }
func (p places) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p places) Less(i, j int) bool { return p[i].pos < p[j].pos }

type Replacer struct {
	original    []byte
	places      places
	parseBuffer *bytes.Buffer
	delimiter   []byte
	lenDel      int
}

func (r *Replacer) SetDelimiter(del delimiter) {
	r.delimiter = delimiterBytes[del]
	r.lenDel = len(r.delimiter)
}

func (r *Replacer) Delimiter() []byte { return r.delimiter }

// returns a new replacer
func New() *Replacer {
	r := &Replacer{}
	r.parseBuffer = &bytes.Buffer{}
	r.SetDelimiter(DefaultDelimiter)
	return r
}

func (r *Replacer) Replace(buffer *bytes.Buffer, m map[string]string) {
	last := 0
	for _, place := range r.places {
		buffer.Write(r.original[last:place.pos])
		if repl, ok := m[place.placeholder]; ok {
			buffer.WriteString(repl)
		}
		last = place.pos
	}
	buffer.Write(r.original[last:len(r.original)])
}

func (r *Replacer) Set(buffer *bytes.Buffer, m map[string]io.WriterTo) (errors map[string]error) {
	last := 0
	errors = map[string]error{}
	for _, place := range r.places {
		buffer.Write(r.original[last:place.pos])
		if repl, ok := m[place.placeholder]; ok {
			_, err := repl.WriteTo(buffer)
			if err != nil {
				fmt.Printf("error: %s", err.Error())
				errors[place.placeholder] = err
				return
			}
		}
		last = place.pos
	}
	buffer.Write(r.original[last:len(r.original)])
	return
}

func (r *Replacer) MustParse(in []byte) *Replacer {
	err := r.Parse(in)
	if err != nil {
		panic(fmt.Sprintf("parse error: %s", err.Error()))
	}
	return r
}

func (r *Replacer) Parse(in []byte) error {
	r.parseBuffer.Reset()
	lenIn := len(in)
	r.places = []place{}
	for i := 0; i < lenIn; i++ {
		found := bytes.Index(in[i:], r.delimiter)
		if -1 < found {
			if i != 0 && found == 0 {
				return fmt.Errorf("Syntax error: can't have 2 or more placeholders side by side: %#v\n", string(in[:i+r.lenDel]))
			}
			start := found + i
			r.parseBuffer.Write(in[i:start])
			startPlaceH := start + r.lenDel
			found = bytes.Index(in[startPlaceH:], r.delimiter)
			if -1 == found {
				r.parseBuffer.Write(in[startPlaceH:])
				break
			} else {
				end := found + start + r.lenDel
				pos := r.parseBuffer.Len()
				r.places = append(r.places, place{pos, string(in[startPlaceH:end])})
				i = end + 1
			}
		} else {
			r.parseBuffer.Write(in[i:])
			break
		}
	}
	r.original = r.parseBuffer.Bytes()
	sort.Sort(r.places)
	return nil
}
