package replacer

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

var _template = []byte{}
var expected = ""
var _map = map[string]string{}

func Prepare() {
	_map = map[string]string{}
	orig := []string{}
	orig2 := []string{}
	exp := []string{}
	for i := 0; i < 5; i++ {
		orig = append(orig, fmt.Sprintf(`a string with @@replacement%v@@`, i))
		orig2 = append(orig2, fmt.Sprintf(`a string with <@replacement%v@>`, i))
		exp = append(exp, fmt.Sprintf("a string with repl%v", i))
		_map[fmt.Sprintf("replacement%v", i)] = fmt.Sprintf("repl%v", i)
	}
	expected = strings.Join(exp, "")
	_template = []byte(strings.Join(orig, ""))
}

var repl = New()

func TestReplaceMulti(t *testing.T) {
	Prepare()
	repl.Parse(_template)
	var buffer bytes.Buffer
	if repl.Replace(&buffer, _map); buffer.String() != expected {
		t.Errorf("unexpected result: %#v, expected: %#v", buffer.String(), expected)
	}
}

func TestReplaceSyntaxError(t *testing.T) {
	ſ := repl.Parse([]byte("before @@one@@@@two@@ after"))
	if ſ == nil {
		t.Errorf("expected syntax error for 2 placeholders side by side, got none")
	}
}
