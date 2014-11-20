// Copyright 2014 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

package main

import (
	"log"
	"os"
	"strings"
	"text/template"

	"github.com/gopkg/web/webgo"
)

const page = `
<html>
<meta charset="utf-8"/>
<body>
{{if .Value}}.
Hi {{.Value}}.
<form method="post" action="/logout">
<input type="submit" name="method" value="logout" />
</form>
You will logout after 10 seconds. Then try to reload.
{{else}}
<form method="post" action="/login">
<label for="name">Name:</label>
<input type="text" id="name" name="name" value="" />
<input type="submit" name="method" value="login" />
</form>
{{end}}
</body>
</html>
`

var tmpl = template.Must(template.New("x").Parse(page))

func getSession(ctx *webgo.Context, manager *webgo.SessionManager) *webgo.Session {
	id, _ := ctx.GetSecureCookie("SessionId")
	session := manager.GetSessionById(id)
	ctx.SetSecureCookie("SessionId", session.Id, int64(manager.GetTimeout()))
	ctx.SetHeader("Pragma", "no-cache", true)
	return session
}

func main() {
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	manager := webgo.NewSessionManager(logger)
	manager.OnStart(func(session *webgo.Session) {
		println("started new session")
	})
	manager.OnEnd(func(session *webgo.Session) {
		println("abandon")
	})
	manager.SetTimeout(10)

	webgo.Config.CookieSecret = "7C19QRmwf3mHZ9CPAaPQ0hsWeufKd"
	webgo.Get("/", func(ctx *webgo.Context) {
		session := getSession(ctx, manager)
		tmpl.Execute(ctx, session)
	})
	webgo.Post("/login", func(ctx *webgo.Context) {
		name := strings.Trim(ctx.Params["name"], " ")
		if name != "" {
			logger.Printf("User \"%s\" login", name)

			// XXX: set user own object.
			getSession(ctx, manager).Value = name
		}
		ctx.Redirect(302, "/")
	})
	webgo.Post("/logout", func(ctx *webgo.Context) {
		session := getSession(ctx, manager)
		if session.Value != nil {
			// XXX: get user own object.
			logger.Printf("User \"%s\" logout", session.Value.(string))
			session.Abandon()
		}
		ctx.Redirect(302, "/")
	})
	webgo.Run(":6061")
}
