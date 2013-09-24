package replacer

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

var Template = []byte{}
var Expected = ""
var Map = map[string]string{}

func Prepare() {
	Map = map[string]string{}
	orig := []string{}
	exp := []string{}
	for i := 0; i < 5; i++ {
		orig = append(orig, fmt.Sprintf(`a string with @@replacement%v@@`, i))
		exp = append(exp, fmt.Sprintf("a string with repl%v", i))
		Map[fmt.Sprintf("replacement%v", i)] = fmt.Sprintf("repl%v", i)
	}
	Expected = strings.Join(exp, "")
	Template = []byte(strings.Join(orig, ""))
}

var repl = New()

func TestReplaceMulti(t *testing.T) {
	Prepare()
	repl.Parse(Template)
	var buffer bytes.Buffer
	if repl.Replace(Map, &buffer); buffer.String() != Expected {
		t.Errorf("unexpected result: %#v, expected: %#v", buffer.String(), Expected)
	}
}

func TestReplaceSyntaxError(t *testing.T) {
	ſ := repl.Parse([]byte("before @@one@@@@two@@ after"))
	if ſ == nil {
		t.Errorf("expected syntax error for 2 placeholders side by side, got none")
	}
}
