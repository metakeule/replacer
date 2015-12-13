package benchmark

import (
	"bytes"
	"fmt"
	"github.com/metakeule/replacer"
	"github.com/metakeule/replacer/places"
	"io"
	"regexp"
	"strings"
	"testing"
)

var (
	StringT   = "a string with @@replacement1@@ and @@replacement2@@ that c@ntinues"
	StringTx  = "a string with <@replacement1@> and <@replacement2@> that c@ntinues"
	TemplateT = "a string with {{.replacement1}} and {{.replacement2}} that c@ntinues"
	ByteT     = []byte(StringT)
	ByteTx    = []byte(StringTx)
	Expected  = "a string with repl1 and repl2 that c@ntinues"

	StringN   = ""
	StringNx  = ""
	TemplateN = ""
	ByteN     = []byte{}
	ByteNx    = []byte{}
	ExpectedN = ""

	StringM   = ""
	StringMx  = ""
	TemplateM = ""
	ByteM     = []byte{}
	ByteMx    = []byte{}
	ExpectedM = ""
)

var (
	Map = map[string]string{
		"@@replacement1@@": "repl1",
		"@@replacement2@@": "repl2",
	}

	Strings = []string{"@@replacement1@@", "repl1", "@@replacement2@@", "repl2"}

	StringMap = map[string]string{
		"replacement1": "repl1",
		"replacement2": "repl2",
	}

	ByteMap = map[string][]byte{
		"@@replacement1@@": []byte("repl1"),
		"@@replacement2@@": []byte("repl2"),
	}

	seekerMap = map[string]io.ReadSeeker{
		"replacement1": strings.NewReader("repl1"),
		"replacement2": strings.NewReader("repl2"),
	}

	MapM       = map[string]string{}
	StringMapM = map[string]string{}
	ByteMapM   = map[string][]byte{}
	StringsM   = []string{}
)

var (
	mapperNaive = &Naive{}
	naive2      = &Naive2{}
	mapperReg   = &Regexp{Regexp: regexp.MustCompile("(@@[^@]+@@)")}
	byts        = &Bytes{}
	repl        = replacer.New()
	templ       = NewTemplate()
)

func PrepareM() {
	MapM = map[string]string{}
	ByteMapM = map[string][]byte{}
	StringMapM = map[string]string{}
	StringsM = []string{}
	s := []string{}
	sx := []string{}
	r := []string{}
	t := []string{}
	for i := 0; i < 5000; i++ {
		s = append(s, fmt.Sprintf(`a string with @@replacement%v@@`, i))
		sx = append(sx, fmt.Sprintf(`a string with <@replacement%v@>`, i))
		t = append(t, fmt.Sprintf(`a string with {{.replacement%v}}`, i))
		r = append(r, fmt.Sprintf("a string with repl%v", i))
		key := fmt.Sprintf("replacement%v", i)
		val := fmt.Sprintf("repl%v", i)
		MapM["@@"+key+"@@"] = val
		ByteMapM["@@"+key+"@@"] = []byte(val)
		StringMapM[key] = val
		StringsM = append(StringsM, "@@"+key+"@@", val)
	}
	StringM = strings.Join(s, "")
	StringMx = strings.Join(sx, "")
	TemplateM = strings.Join(t, "")
	ExpectedM = strings.Join(r, "")
	ByteM = []byte(StringM)
	ByteMx = []byte(StringMx)
}

func PrepareN() {
	s := []string{}
	sx := []string{}
	r := []string{}
	t := []string{}
	for i := 0; i < 2500; i++ {
		s = append(s, StringT)
		sx = append(sx, StringTx)
		r = append(r, Expected)
		t = append(t, TemplateT)
	}
	TemplateN = strings.Join(t, "")
	StringN = strings.Join(s, "")
	StringNx = strings.Join(sx, "")
	ExpectedN = strings.Join(r, "")
	ByteN = []byte(StringN)
	ByteNx = []byte(StringNx)

}

func TestReplace(t *testing.T) {
	mapperNaive.Map = Map
	mapperNaive.Template = StringT
	if r := mapperNaive.Replace(); r != Expected {
		t.Errorf("unexpected result for %s: %#v", "mapperNaive", r)
	}

	naive2.Replacements = Strings
	naive2.Template = StringT
	if r := naive2.Replace(); r != Expected {
		t.Errorf("unexpected result for %s: %#v", "naive2", r)
	}

	mapperReg.Map = Map
	mapperReg.Template = StringT
	mapperReg.Setup()
	if r := mapperReg.Replace(); r != Expected {
		t.Errorf("unexpected result for %s: %#v", "mapperReg", r)
	}

	byts.Map = ByteMap
	byts.Parse(StringT)
	if r := byts.Replace(); string(r) != Expected {
		t.Errorf("unexpected result for %s: %#v, expected: %#v", "byts", string(r), Expected)
	}

	templ.Parse(TemplateT)
	var tbf bytes.Buffer
	if templ.Replace(StringMap, &tbf); tbf.String() != Expected {
		t.Errorf("unexpected result for %s: %#v, expected: %#v", "template", tbf.String(), Expected)
	}

	err := repl.Parse(ByteT)
	if err != nil {
		panic(err.Error())
	}

	var bf bytes.Buffer
	if repl.Replace(&bf, StringMap); bf.String() != Expected {
		t.Errorf("unexpected result for %s: %#v", "fastreplace2", bf.String())
	}
}

func TestReplaceN(t *testing.T) {
	PrepareN()
	mapperNaive.Map = Map
	mapperNaive.Template = StringN
	if r := mapperNaive.Replace(); r != ExpectedN {
		t.Errorf("unexpected result for %s: %#v", "mapperNaive", r)
	}

	naive2.Replacements = Strings
	naive2.Template = StringN
	if r := naive2.Replace(); r != ExpectedN {
		t.Errorf("unexpected result for %s: %#v", "naive2", r)
	}

	mapperReg.Map = Map
	mapperReg.Template = StringN
	mapperReg.Setup()
	if r := mapperReg.Replace(); r != ExpectedN {
		t.Errorf("unexpected result for %s: %#v", "mapperReg", r)
	}

	templ.Parse(TemplateN)
	var tbf bytes.Buffer
	if templ.Replace(StringMap, &tbf); tbf.String() != ExpectedN {
		t.Errorf("unexpected result for %s: %#v, expected: %#v", "template", tbf.String(), ExpectedN)
	}

	err := repl.Parse(ByteN)

	if err != nil {
		panic(err.Error())
	}

	var bf bytes.Buffer
	if repl.Replace(&bf, StringMap); bf.String() != ExpectedN {
		t.Errorf("unexpected result for %s: %#v", "fastreplace2", bf.String())
	}
}

func TestReplaceM(t *testing.T) {
	PrepareM()
	mapperNaive.Map = MapM
	mapperNaive.Template = StringM
	if r := mapperNaive.Replace(); r != ExpectedM {
		t.Errorf("unexpected result for %s: %#v", "mapperNaive", r)
	}

	naive2.Replacements = StringsM
	naive2.Template = StringM
	if r := naive2.Replace(); r != ExpectedM {
		t.Errorf("unexpected result for %s: %#v", "naive2", r)
	}

	mapperReg.Map = MapM
	mapperReg.Template = StringM
	mapperReg.Setup()
	if r := mapperReg.Replace(); r != ExpectedM {
		t.Errorf("unexpected result for %s: %#v", "mapperReg", r)
	}

	naive2.Replacements = StringsM
	naive2.Template = StringM
	if r := naive2.Replace(); r != ExpectedM {
		t.Errorf("unexpected result for %s: %#v", "naive2", r)
	}

	templ.Parse(TemplateM)
	var tbf bytes.Buffer
	if templ.Replace(StringMapM, &tbf); tbf.String() != ExpectedM {
		t.Errorf("unexpected result for %s: %#v, expected: %#v", "template", tbf.String(), ExpectedM)
	}

	err := repl.Parse(ByteM)
	if err != nil {
		panic(err.Error())
	}

	var bf bytes.Buffer

	if repl.Replace(&bf, StringMapM); bf.String() != ExpectedM {
		t.Errorf("unexpected result for %s: %#v", "fastreplace2", bf.String())
	}
}

func BenchmarkNaive(b *testing.B) {
	b.StopTimer()
	PrepareN()
	mapperNaive.Map = Map
	mapperNaive.Template = StringN
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		mapperNaive.Replace()
	}
}

func BenchmarkNaive2(b *testing.B) {
	b.StopTimer()
	PrepareN()
	naive2.Replacements = Strings
	naive2.Template = StringN
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		naive2.Replace()
	}
}

func BenchmarkReg(b *testing.B) {
	b.StopTimer()
	PrepareN()
	mapperReg.Map = Map
	mapperReg.Template = StringN
	mapperReg.Setup()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		mapperReg.Replace()
	}
}

func BenchmarkByte(b *testing.B) {
	b.StopTimer()
	PrepareN()
	byts.Map = ByteMap
	byts.Parse(StringN)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		byts.Replace()
	}
}
func BenchmarkTemplate(b *testing.B) {
	b.StopTimer()
	PrepareN()
	templ.Parse(TemplateN)
	var tbf bytes.Buffer
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		templ.Replace(StringMap, &tbf)
		tbf.Reset()
	}
}

func BenchmarkReplacer(b *testing.B) {
	b.StopTimer()
	PrepareN()
	repl.Parse(ByteN)
	var bf bytes.Buffer
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		repl.Replace(&bf, StringMap)
		bf.Reset()
	}
}

func BenchmarkPlaces(b *testing.B) {
	b.StopTimer()
	PrepareN()
	var pl = places.Find(ByteNx)
	var bf bytes.Buffer
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		places.ReplaceString(ByteNx, &bf, pl, StringMap)
		bf.Reset()
	}
}

func BenchmarkNaiveM(b *testing.B) {
	b.StopTimer()
	PrepareM()
	mapperNaive.Map = MapM
	mapperNaive.Template = StringM
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		mapperNaive.Replace()
	}
}

func BenchmarkNaive2M(b *testing.B) {
	b.StopTimer()
	PrepareM()
	naive2.Replacements = StringsM
	naive2.Template = StringM
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		naive2.Replace()
	}
}

func BenchmarkRegM(b *testing.B) {
	b.StopTimer()
	PrepareM()
	mapperReg.Map = MapM
	mapperReg.Template = StringM
	mapperReg.Setup()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		mapperReg.Replace()
	}
}

func BenchmarkByteM(b *testing.B) {
	b.StopTimer()
	PrepareM()
	byts.Map = ByteMap
	byts.Parse(StringM)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		byts.Replace()
	}
}

func BenchmarkTemplateM(b *testing.B) {
	b.StopTimer()
	PrepareM()
	templ.Parse(TemplateM)
	var tbf bytes.Buffer
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		templ.Replace(StringMapM, &tbf)
		tbf.Reset()
	}
}

func BenchmarkReplacerM(b *testing.B) {
	b.StopTimer()
	PrepareM()
	repl.Parse(ByteM)
	var bf bytes.Buffer
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		repl.Replace(&bf, StringMapM)
		bf.Reset()
	}
}

func BenchmarkPlacesM(b *testing.B) {
	b.StopTimer()
	PrepareM()
	// println(StringMx)
	var pl = places.Find(ByteMx)
	var bf bytes.Buffer
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		places.ReplaceString(ByteMx, &bf, pl, StringMapM)
		bf.Reset()
	}
}

func BenchmarkOnceNaive(b *testing.B) {
	for i := 0; i < b.N; i++ {
		mapperNaive.Map = Map
		mapperNaive.Template = StringN
		mapperNaive.Replace()
	}
}

func BenchmarkOnceNaive2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		naive2.Replacements = Strings
		naive2.Template = StringN
		naive2.Replace()
	}
}

func BenchmarkOnceReg(b *testing.B) {
	mapperReg.Setup()
	for i := 0; i < b.N; i++ {
		mapperReg.Map = Map
		mapperReg.Template = StringN
		mapperReg.Replace()
	}
}

func BenchmarkOnceByte(b *testing.B) {
	for i := 0; i < b.N; i++ {
		byts.Parse(StringN)
		byts.Map = ByteMap
		byts.Replace()
	}
}

func BenchmarkOnceTemplate(b *testing.B) {
	for i := 0; i < b.N; i++ {
		templ.Parse(TemplateN)
		var tbf bytes.Buffer
		templ.Replace(StringMap, &tbf)
	}
}

func BenchmarkOnceReplacer(b *testing.B) {
	for i := 0; i < b.N; i++ {
		repl.Parse(ByteN)
		var bf bytes.Buffer
		repl.Replace(&bf, StringMap)
	}
}

func BenchmarkOncePlaces(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var bf bytes.Buffer
		places.FindAndReplaceString(ByteNx, &bf, StringMap)
	}
}
