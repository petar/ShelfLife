// Copyright 2011 ShelfLife Authors. All rights reserved.
// Use of this source code is governed by a
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"github.com/petar/GoHTTP/server"
	"github.com/petar/GoHTTP/server/subs"

	"github.com/petar/ShelfLife/api"
)

var (
	flagMongoAddr = flag.String("mongo", "", "IP address of MongoDB server")
	flagBind      = flag.String("bind", "0.0.0.0:3300", "Address to bind server to")
	flagParallel  = flag.Int("parallel", 1, "Number of requests served in parallel")
)

func main() {
	fmt.Fprintf(os.Stderr, "ShelfLife API Server — 2011\n")
	flag.Parse()

	srv, err := server.NewServerEasy(*flagBind)
	if err != nil {
		log.Fatalf("Problem binding server: %s\n", err)
	}

	// Attach server subs and extensions HERE
	api := subs.NewAPI()
	if err := api.Register(social.NewAPI()); err != nil {
		log.Fatalf("Problem registering social API: %s\n", err)
	}
	srv.AddSub("/api/", api)

	fmt.Printf("· Serving %d requests in parallel ...\n", *flagParallel)
	srv.Launch()
}
