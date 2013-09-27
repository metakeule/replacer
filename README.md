replacer
========

fast and simple templating for go

[![Build Status](https://secure.travis-ci.org/metakeule/replacer.png)](http://travis-ci.org/metakeule/replacer)

If you need to simply replace placeholders in a template without escaping or logic,
replacer might be for you.

For the typical scenario - your template never changes on runtime -, replacer is faster than using (strings|bytes).Replace(r)() or regexp.ReplaceAllStringFunc() or the text/template package.

Performance
-----------

Runing benchmarks in the benchmark directory, I get the following results:

replacing 2 placeholders that occur 2500x in the template

    BenchmarkNaive      500    3035929 ns/op  3,8x (strings.Replace)
    BenchmarkNaive2    1000    2595076 ns/op  3,2x (strings.Replacer)
    BenchmarkReg        100   25882258 ns/op 32,3x (regexp.ReplaceAllStringFunc)
    BenchmarkByte      1000    2210725 ns/op  2,8x (bytes.Replace)
    BenchmarkTemplate   500    6373070 ns/op  7,9x (template.Execute)
    BenchmarkReplacer  2000     802490 ns/op  1,0x (replacer.Replace)

replacing 5000 placeholders that occur 1x in the template

    BenchmarkNaiveM        1   4317513185 ns/op 3929,6x (strings.Replace)
    BenchmarkNaive2M     500      6329720 ns/op    5,8x (strings.Replacer)
    BenchmarkRegM         50     31198202 ns/op   28,4x (regexp.ReplaceAllStringFunc)
    BenchmarkByteM      1000      1475455 ns/op    1,3x (bytes.Replace)
    BenchmarkTemplateM   500      6435381 ns/op    5,9x (template.Execute)
    BenchmarkReplacerM  2000      1098709 ns/op    1,0x (replacer.Replace)

replacing 2 placeholders that occur 1x in the template, parsing template each time (you should not do this)

    BenchmarkOnceNaive    1000     3037135 ns/op   1,2x (strings.Replace)
    BenchmarkOnceNaive2   1000     2600541 ns/op   1,1x (strings.Replacer)
    BenchmarkOnceReg        50    26540129 ns/op  10,7x (regexp.ReplaceAllStringFunc)
    BenchmarkOnceByte     1000     2471198 ns/op   1,0x (bytes.Replace)
    BenchmarkOnceTemplate    5   978977017 ns/op 396,2x (template.Execute)
    BenchmarkOnceReplacer  500     4572535 ns/op   1,9x (replacer.Replace)

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
    r.Replace(&buffer, m)
    
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

Status
------

The package is stable and ready for consumption.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.

see LICENSE file.

GoDoc
-----

see http://godoc.org/github.com/metakeule/replacer
