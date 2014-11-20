// Copyright 2014 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gopkg/web/webgo"
)

func hello(ctx *webgo.Context, num string) {
	flusher, _ := ctx.ResponseWriter.(http.Flusher)
	flusher.Flush()
	n, _ := strconv.ParseInt(num, 10, 64)
	for i := int64(0); i < n; i++ {
		ctx.WriteString("<br>hello world</br>")
		flusher.Flush()
		time.Sleep(1e9)
	}
}

func main() {
	webgo.Get("/([0-9]+)", hello)
	webgo.Run("0.0.0.0:9999")
}
