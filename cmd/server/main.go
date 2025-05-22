package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/valyala/fasthttp"
	"github.com/yqnk/urlshortener/internal/server"
)

func main() {
	service, err := server.NewService("./urlshortener.db")
	if err != nil {
		log.Fatalf("failed to initialize service: %v", err)
	}

	handler := server.NewHandler(service)

	host := flag.String("host", "localhost", "a string")
	port := flag.String("port", "3333", "a string")

	flag.Parse()

	addr := *host + ":" + *port
	fmt.Printf("Running at %s...\n", addr)
	if err := fasthttp.ListenAndServe(addr, handler); err != nil {
		log.Fatalf("failted to start server: %v", err)
	}
}
