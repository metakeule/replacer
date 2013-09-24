replacer
========

fast and simple templating for go

[![Build Status](https://secure.travis-ci.org/metakeule/replacer.png)](http://travis-ci.org/metakeule/replacer)

If you need to simply replace placeholders in a template without escaping or logic,
replacer might be for you.

For the typical scenario - your template never changes on runtime -, replacer is faster than using (strings|bytes).Replace() or regexp.ReplaceAllStringFunc() or the text/template package.

Performance
-----------

Runing benchmarks in the benchmark directory, I get the following results:

replacing 2 placeholders that occur 2500x in the template

    BenchmarkNaive      500    6112228 ns/op  3,7x (strings.Replace)
    BenchmarkReg         50   53740939 ns/op 32,8x (regexp.ReplaceAllStringFunc)
    BenchmarkByte       500    4627244 ns/op  2,8x (bytes.Replace)
    BenchmarkTemplate   100   13838150 ns/op  8,4x (template.Execute)
    BenchmarkReplacer  1000    1640622 ns/op  1,0x (replacer.Replace)

replacing 5000 placeholders that occur 1x in the template

    BenchmarkNaiveM        1   8663141464 ns/op 3941,1x (strings.Replace)
    BenchmarkRegM         50     63944139 ns/op   29,1x (regexp.ReplaceAllStringFunc)
    BenchmarkByteM         1   5955402986 ns/op 2709,3x (bytes.Replace)
    BenchmarkTemplateM   100     13903383 ns/op    6,3x (template.Execute)
    BenchmarkReplacerM  1000      2198166 ns/op    1,0x (replacer.Replace)

replacing 2 placeholders that occur 1x in the template, parsing template each time (you should not do this)

    BenchmarkOnceNaive     500 ns/op 6145737      1,2x (strings.Replace)
    BenchmarkOnceReg        50 ns/op 54311308    11,0x (regexp.ReplaceAllStringFunc)
    BenchmarkOnceByte      500 ns/op 4950499      1,0x (bytes.Replace)
    BenchmarkOnceTemplate    1 ns/op 2066008619 417,3x (template.Execute)
    BenchmarkOnceReplacer  200 ns/op 8725035      1,8x (replacer.Replace)


Usage
-----

```go
package main

import (
    "bytes"
    "fmt"
    "github.com/metakeule/replacer"
)

func main() {
    r := replacer.New()
    
    // reuse r to speed up parsing of
    // different templates on the fly
    // (for concurrency you need to protect it with a mutex)
    err := r.Parse([]byte("@@name@@ @@animal@@"))
    if err != nil {
        panic(err.Error())
    }
    
    m := map[string]string{
        "animal": "Duck",
        "name":   "Donald",
    }

    var buffer bytes.Buffer
    
    // reuse r with a parsed template to speed up replacement
    r.Replace(m, &buffer)
    
    // after the replacement you may use the buffer methods Bytes(), String(), Write() or WriteTo()
    // and reuse the same buffer after calling buffer.Reset()
    fmt.Println(buffer.String())
}
```

results in

```
Donald Duck
```

Limitation
----------

Two placeholders immeadiatly following each other are not allowed, e.g.
    
    @@firstname@@@@lastname@@

However, you should be able to combine them in another placeholder and replace the combination.
As long as 1 byte is between them, it is no problem, e.g.

    @@firstname@@ @@lastname@@

GoDoc
-----

see http://godoc.org/github.com/metakeule/replacer
