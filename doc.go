// Copyright 2013 Marc René Arns. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package replacer performs fast replacements of placeholders in a []byte.
It has no logic and no escaping.

For the typical scenario - your template never changes on runtime -,
replacer is faster than using (strings|bytes).Replace() or
regexp.ReplaceAllStringFunc() or the text/template package.

You might run the benchmarks in the benchmark directory and
have a look at the example directory.
*/
package replacer
