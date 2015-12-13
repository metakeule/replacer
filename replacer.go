package replacer

import (
	"bytes"
	"errors"
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
	original []byte
	places   places
	// parseBuffer bytes.Buffer
	delimiter []byte
	lenDel    int
}

func (r *Replacer) SetDelimiter(del delimiter) {
	r.delimiter = delimiterBytes[del]
	r.lenDel = len(r.delimiter)
}

func (r *Replacer) Delimiter() []byte { return r.delimiter }

// returns a new replacer
func New() Replacer {
	r := Replacer{}
	r.SetDelimiter(DefaultDelimiter)
	return r
}

func (r *Replacer) Replace(buffer *bytes.Buffer, m map[string]string) {
	var (
		last int
		repl string
		ok   bool
	)
	for _, place := range r.places {
		buffer.Write(r.original[last:place.pos])
		repl, ok = m[place.placeholder]
		if ok {
			buffer.WriteString(repl)
		}
		last = place.pos
	}
	buffer.Write(r.original[last:len(r.original)])
}

func (r *Replacer) Set(buffer *bytes.Buffer, m map[string]io.WriterTo) (errors map[string]error) {
	var (
		last int
		repl io.WriterTo
		ok   bool
		err  error
	)
	errors = map[string]error{}
	for _, place := range r.places {
		buffer.Write(r.original[last:place.pos])
		repl, ok = m[place.placeholder]
		if ok {
			_, err = repl.WriteTo(buffer)
			if err != nil {
				// fmt.Printf("error: %s", err.Error())
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
		panic("parse error: " + err.Error())
	}
	return r
}

func (r *Replacer) Parse(in []byte) error {
	lenIn := len(in)
	lenDel := r.lenDel
	r.places = make([]place, 0, 22)
	r.original = make([]byte, 0, lenIn)
	var (
		found       = -1
		start       int
		startPlaceH int
		end         int
		pos         int
	)
	for i := 0; i < lenIn; i++ {
		found = bytes.Index(in[i:], r.delimiter)
		if -1 < found {
			if i != 0 && found == 0 {
				return errors.New("Syntax error: can't have 2 or more placeholders side by side: " + string(in[:i+lenDel]))
			}
			start = found + i
			r.original = append(r.original, in[i:start]...)
			startPlaceH = start + lenDel
			found = bytes.Index(in[startPlaceH:], r.delimiter)
			if -1 == found {
				r.original = append(r.original, in[startPlaceH:]...)
				break
			} else {
				end = found + start + lenDel
				pos = len(r.original)
				r.places = append(r.places, place{pos, string(in[startPlaceH:end])})
				i = end + 1
			}
		} else {
			r.original = append(r.original, in[i:]...)
			break
		}
	}
	sort.Sort(r.places)
	return nil
}
