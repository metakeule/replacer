replacer
========

fast and simple templating for go

[![Build Status](https://secure.travis-ci.org/metakeule/replacer.png)](http://travis-ci.org/metakeule/replacer)

If you need to simply replace placeholders in a template without escaping or logic, replacer might be for you.

For the typical scenario - your template never changes on runtime -, replacer is faster than using (strings|bytes).Replace(r)() or regexp.ReplaceAllStringFunc() or the text/template package.

There is also a new subpackage called places which is similar in performance for normal rendering but faster for parsing. It has a different API though.

Performance
-----------

Runing benchmarks in the benchmark directory, I get the following results (go1.4, linux64):

replacing 2 placeholders that occur 2500x in the template

    BenchmarkNaive     500    2780308 ns/op     4 allocs/op   7,54x (strings.Replace)
    BenchmarkNaive2   2000    1098662 ns/op    13 allocs/op   2,98x (strings.Replacer)
    BenchmarkReg        50   22916309 ns/op  5024 allocs/op  62,18x (regexp.ReplaceAllStringFunc)
    BenchmarkByte     1000    2087658 ns/op     4 allocs/op   5,66x (bytes.Replace)
    BenchmarkTemplate  300    5566514 ns/op 15002 allocs/op  15,10x (template.Execute)
    BenchmarkReplacer 3000     376034 ns/op     0 allocs/op   1,02x (replacer.Replace)
    BenchmarkPlaces   3000     368525 ns/op     0 allocs/op   1,00x (places.ReplaceString)
                                


replacing 5000 placeholders that occur 1x in the template

    BenchmarkNaiveM       1 4286720867 ns/op 10000 allocs/op 6673,23x (strings.Replace)
    BenchmarkNaive2M    500    4019384 ns/op 11007 allocs/op    6,26x (strings.Replacer)
    BenchmarkRegM        50   27298490 ns/op  5025 allocs/op   42,50x (regexp.ReplaceAllStringFunc)
    BenchmarkByteM     1000    1626838 ns/op     4 allocs/op    2,53x (bytes.Replace)
    BenchmarkTemplateM  300    5667141 ns/op 15002 allocs/op    8,82x (template.Execute)
    BenchmarkReplacerM 2000     643043 ns/op     0 allocs/op    1,00x (replacer.Replace)
    BenchmarkPlacesM   2000     642376 ns/op     0 allocs/op    1,00x (places.ReplaceString)
                                

replacing 2 placeholders that occur 1x in the template, parsing template each time (you should not do this until you need it)

    BenchmarkOnceNaive    500   2759530 ns/op     4 allocs/op   2,45x (strings.Replace)
    BenchmarkOnceNaive2  2000   1127832 ns/op    13 allocs/op   1,00x (strings.Replacer)
    BenchmarkOnceReg       50  23076371 ns/op  5024 allocs/op  20,46x (regexp.ReplaceAllStringFunc)
    BenchmarkOnceByte    1000   2336374 ns/op     6 allocs/op   2,07x (bytes.Replace)
    BenchmarkOnceTemplate   2 917598881 ns/op 60058 allocs/op 813,60x (template.Execute)
    BenchmarkOnceReplacer 500   3510982 ns/op  5025 allocs/op   3,11x (replacer.Replace)
    BenchmarkOncePlaces  1000   1808015 ns/op    26 allocs/op   1,60x (places.ReplaceString)


places
========

Usage
-----

```go
package main

import (
    "bytes"
    "fmt"
    "github.com/metakeule/replacer/places"
)

func main() {
    // parse the template once
    template := places.NewTemplate([]byte("<@name@>: <@animal@>"))    
    
    
    // reuse it to speed up replacement
    var buffer bytes.Buffer
    template.ReplaceString(&buffer, map[string]string{"animal": "Duck","name": "Donald"})

    // there are alternative methods for Bytes, io.ReadSeeker etc.
    
    // after the replacement you may use the buffer methods Bytes(), String(), Write() or WriteTo()
    // and reuse the same buffer after calling buffer.Reset()
    fmt.Println(buffer.String())
}
```


results in

```
Donald: Duck
```

Documentation (GoDoc)
---------------------

see https://godoc.org/github.com/metakeule/replacer/places

replacer
========

For compatibility there is still the more limited and slower parsing `replacer` library:

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
