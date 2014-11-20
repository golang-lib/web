// Copyright 2014 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

package main

import (
	"log"
	"os"

	"github.com/gopkg/web/webgo"
)

func hello(val string) string { return "hello " + val }

func main() {
	f, err := os.Create("server.log")
	if err != nil {
		println(err.Error())
		return
	}
	logger := log.New(f, "", log.Ldate|log.Ltime)
	webgo.Get("/(.*)", hello)
	webgo.SetLogger(logger)
	webgo.Run("0.0.0.0:9999")
}
