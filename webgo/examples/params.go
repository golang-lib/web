// Copyright 2014 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

package main

import (
	"fmt"

	"github.com/gopkg/web/webgo"
)

func index() string {
	return page
}

func process(ctx *webgo.Context) string {
	return fmt.Sprintf("%v\n", ctx.Params)
}

func main() {
	webgo.Get("/", index)
	webgo.Post("/process", process)
	webgo.Run("0.0.0.0:9999")
}

const page = `
<html>
	<head>
		<title>Multipart Test</title>
	</head>
	<body>
		<form action="/process" method="POST">
			<label for="a"> Please write some text </label>
			<input id="a" type="text" name="a"/>
			<br>

			<label for="b"> Please write some more text </label>
			<input id="b" type="text" name="b"/>
			<br>

			<label for="c"> Please write a number </label>
			<input id="c" type="text" name="c"/>
			<br>

			<label for="d"> Please write another number </label>
			<input id="d" type="text" name="d"/>
			<br>

			<input type="submit" name="Submit" value="Submit"/>
		</form>
	</body>
</html>
`
