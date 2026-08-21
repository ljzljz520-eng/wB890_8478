package main

import (
	"flag"
	"fmt"
	"log"
	"memorialstation/api"
	"memorialstation/storage"
	"net/http"
	"os"
)

func main() {
	path := flag.String("db", "memorial.db", "database path")
	addr := flag.String("addr", ":8080", "listen address")
	flag.Parse()
	store, err := storage.Open(*path)
	if err != nil {
		log.Fatal(err)
	}
	defer store.Close()
	server := api.New(store)
	if os.Getenv("MEMORIAL_ONCE") == "1" {
		fmt.Println("memorial station ready")
		return
	}
	log.Printf("memorial station listening on %s", *addr)
	if err := http.ListenAndServe(*addr, server.Handler()); err != nil {
		log.Fatal(err)
	}
}
