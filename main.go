package main

import (
	"flag"
	"fmt"
	"log"
)

func main() {
	var (
		dbFile  = flag.String("db", "wordle.db", "sqlite db file")
		hostKey = flag.String("key", "key.pem", "key")
		port    = flag.String("port", "22", "port")
	)
	flag.Parse()

	repo, err := newRepo(*dbFile)
	if err != nil {
		log.Fatal(err)
	}

	server, err := newServer(repo, *hostKey, *port)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("listening on :%s\n", *port)
	log.Fatal(server.ListenAndServe())
}
