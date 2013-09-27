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
