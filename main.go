package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/figment-networks/near-indexer/near/client"
	"github.com/figment-networks/near-indexer/server"
)

func main() {
	var endpoint string
	var addr string
	var port int

	flag.StringVar(&endpoint, "endpoint", "", "Near service endpoint")
	flag.StringVar(&addr, "addr", "0.0.0.0", "Server listen address")
	flag.IntVar(&port, "port", 5555, "Server listen port")
	flag.Parse()

	if endpoint == "" {
		log.Fatal("endpoint is required")
	}

	rpc := client.New(endpoint)
	srv := server.New(&rpc)

	listenAddr := fmt.Sprintf("%v:%v", addr, port)

	log.Println("starting server on", listenAddr)
	if err := srv.Run(listenAddr); err != nil {
		log.Fatal(err)
	}
}
