// Copyright 2013 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package webgo is a lightweight webgo framework for Go.

It's ideal for writing simple, performant backend webgo services.

Example

	package main

	import (
		"dev.visiontek.wh/webgo"
	)

	func hello(val string) string { return "hello " + val }

	func main() {
		webgo.Get("/(.*)", hello)
		webgo.Run("0.0.0.0:9999")
	}

Getting parameters

	package main

	import (
		"dev.visiontek.wh/webgo"
	)

	func hello(ctx *webgo.Context, val string) {
		for k,v := range ctx.Params {
			println(k, v)
		}
	}

	func main() {
		webgo.Get("/(.*)", hello)
		webgo.Run("0.0.0.0:9999")
	}

In this example, if you visit `http://localhost:9999/?a=1&b=2`,
you'll see the following printed out in the terminal:

	a 1
	b 2
*/
package webgo
