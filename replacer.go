package replacer

import (
	"bytes"
	"fmt"
	"sort"
)

type delimiter int

const (
	DefaultDelimiter delimiter = iota
	HashDelimiter
	DollarDelimiter
)

var delimiterBytes = map[delimiter][]byte{
	DefaultDelimiter: []byte("@@"),
	HashDelimiter:    []byte("##"),
	DollarDelimiter:  []byte("$$"),
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

// instead of the real struct we export an interface
type Replacer interface {
	// Parses the input for placeholders and caches the result
	// should be called once per template
	// must be protected if used concurrently on the same replacer
	// returns an error if 2 placeholders are directly following
	// each other without a byte between them
	Parse([]byte) error

	// Replaces the placeholders that are keys in the given map
	// and writes the resulting bytes to the given buffer.
	// Be aware that the placeholders must not include the delimiter.
	// Bring in your own buffer, allows you to reused it
	Replace(map[string]string, *bytes.Buffer)

	// set the delimiter which surrounds the placeholders
	// valid delimiters are: (delimiter => example)
	//     DefaultDelimiter =>  "@@example@@"
	//     HashDelimiter    =>  "##example##"
	//     DollarDelimiter  =>  "$$example$$"
	SetDelimiter(delimiter)
}

type replace struct {
	original    []byte
	places      places
	parseBuffer *bytes.Buffer
	delimiter   []byte
	lenDel      int
}

func (r *replace) SetDelimiter(del delimiter) {
	r.delimiter = delimiterBytes[del]
	r.lenDel = len(r.delimiter)
}

// returns a new replacer
func New() Replacer {
	r := &replace{}
	r.parseBuffer = &bytes.Buffer{}
	r.SetDelimiter(DefaultDelimiter)
	return r
}

func (r *replace) Replace(m map[string]string, buffer *bytes.Buffer) {
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

func (r *replace) Parse(in []byte) error {
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
