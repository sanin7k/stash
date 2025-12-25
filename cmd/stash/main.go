package main

import (
	"fmt"
	"log"
	"os"

	"stash/internal/client"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("usage: stash <backup|restore|list> <file>")
	}

	cmd := os.Args[1]
	addr := "127.0.0.1:9000"

	switch cmd {
	case "backup":
		if len(os.Args) < 3 {
			log.Fatalf("usage: stash backup <file>")
		}
		client.Backup(addr, os.Args[2])
	case "restore":
		if len(os.Args) < 3 {
			log.Fatalf("usage: stash restore <file>")
		}
		client.Restore(addr, os.Args[2])
	case "list":
		client.List(addr)
	default:
		fmt.Println("unknown command")
	}
}
