// Copyright 2011 ShelfLife Authors. All rights reserved.
// Use of this source code is governed by a
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"log"
	"mime"
	"os"
	"github.com/petar/GoHTTP/server"
	"github.com/petar/GoHTTP/server/rpc"
	"github.com/petar/GoHTTP/server/static"

	"github.com/petar/ShelfLife/db"
	"github.com/petar/ShelfLife/sociability"
)

var (
	flagDbAddr    = flag.String("db", "127.0.0.1:22000", "IP address of DB server")
	flagBind      = flag.String("bind", "0.0.0.0:3300", "Address to bind server to")
	flagStatic    = flag.String("static", "", "If non-empty, static files served from this directory")
	flagParallel  = flag.Int("parallel", 1, "Number of requests served in parallel")
)

func main() {
	fmt.Fprintf(os.Stderr, "ShelfLife Server — 2011\n")
	flag.Parse()

	// Connect to database
	log.Printf("Connecting to DB ...")
	db, err := db.NewDb(*flagDbAddr, "shelflife")
	if err != nil {
		log.Fatalf("Problem connecting to db: %s", err)
	}

	// Start web server
	log.Printf("Starting web server ...")
	srv, err := server.NewServerEasy(*flagBind)
	if err != nil {
		log.Fatalf("Problem binding server: %s\n", err)
	}

	// Attach static file server
	mime.AddExtensionType(".x-jquery-tmpl", "text/x-jquery-tmpl")
	if *flagStatic != "" {
		srv.AddSub("/s/", static.NewStaticSub(*flagStatic))	
	}

	// Attach RPC server module
	rpcsub := rpc.NewRPC()
	if err := rpcsub.RegisterName("ss", sociability.NewAPI(db, []byte{1, 2, 3, 4})); err != nil {
		log.Fatalf("Problem registering social API: %s\n", err)
	}
	srv.AddSub("/api/", rpcsub)

	log.Printf("Serving %d requests in parallel ...\n", *flagParallel)
	srv.Launch(*flagParallel)
	<-make(chan int)
}
