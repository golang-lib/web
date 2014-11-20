// Copyright 2014 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

package main

import (
	"fmt"
	"html"

	"github.com/gopkg/web/webgo"
)

var cookieName = "cookie"

var notice = `
<div>%v</div>
`
var form = `
<form method="POST" action="update">
	<div class="field">
		<label for="cookie"> Set a cookie: </label>
		<input id="cookie" name="cookie"> </input>
	</div>

	<input type="submit" value="Submit"></input>
	<input type="submit" name="submit" value="Delete"></input>
</form>
`

func index(ctx *webgo.Context) string {
	cookie, _ := ctx.Request.Cookie(cookieName)
	var top string
	if cookie == nil {
		top = fmt.Sprintf(notice, "The cookie has not been set")
	} else {
		var val = html.EscapeString(cookie.Value)
		top = fmt.Sprintf(notice, "The value of the cookie is '"+val+"'.")
	}
	return top + form
}

func update(ctx *webgo.Context) {
	if ctx.Params["submit"] == "Delete" {
		ctx.SetCookie(webgo.NewCookie(cookieName, "", -1))
	} else {
		ctx.SetCookie(webgo.NewCookie(cookieName, ctx.Params["cookie"], 0))
	}
	ctx.Redirect(301, "/")
}

func main() {
	webgo.Get("/", index)
	webgo.Post("/update", update)
	webgo.Run("0.0.0.0:9999")
}
