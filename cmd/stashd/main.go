package main

import (
	"log"
	"net"

	"stash/internal/server"
)

func main() {
	ln, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()

	sem := make(chan struct{}, 16)

	log.Println("stashd listening on :9000")

	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}

		sem <- struct{}{}
		go func() {
			defer func() { <-sem }()
			server.HandleConn(conn)
		}()
	}
}
