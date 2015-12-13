package places

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"testing"
)

var _template2 = []byte{}
var expected = ""
var _map = map[string]string{}
var _mapReader = map[string]io.ReadSeeker{}

func Prepare() {
	_map = map[string]string{}
	_mapReader = map[string]io.ReadSeeker{}
	orig2 := []string{}
	exp := []string{}
	for i := 0; i < 5; i++ {
		orig2 = append(orig2, fmt.Sprintf(`a string with <@replacement%v@>`, i))
		exp = append(exp, fmt.Sprintf("a string with repl%v", i))
		_map[fmt.Sprintf("replacement%v", i)] = fmt.Sprintf("repl%v", i)
		_mapReader[fmt.Sprintf("replacement%v", i)] = strings.NewReader(fmt.Sprintf("repl%v", i))
	}
	expected = strings.Join(exp, "")
	_template2 = []byte(strings.Join(orig2, ""))
}

func TestFindAndReplace(t *testing.T) {
	Prepare()
	var buffer bytes.Buffer
	if FindAndReplace(_template2, &buffer, _mapReader); buffer.String() != expected {
		t.Errorf("unexpected result: %#v, expected: %#v", buffer.String(), expected)
	}
}

func TestAjacentPlaceholders(t *testing.T) {
	var buffer bytes.Buffer
	tt := []byte("a string with <@replacement0@><@replacement1@> after")
	ph := Find(tt)
	exp := "a string with repl0repl1 after"

	if Replace(tt, &buffer, ph, _mapReader); buffer.String() != exp {
		t.Errorf("unexpected result: %#v, expected: %#v", buffer.String(), exp)
	}
}

func TestReplaceString(t *testing.T) {
	Prepare()
	var buffer bytes.Buffer
	if ReplaceString(_template2, &buffer, Find(_template2), _map); buffer.String() != expected {
		t.Errorf("unexpected result: %#v, expected: %#v", buffer.String(), expected)
	}
}
