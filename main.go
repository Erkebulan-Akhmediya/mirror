package main

import (
	"log"
	"net/url"
	"os"
)

func main() {
	start := os.Args[1]
	parsedUrl, err := url.Parse(start)
	if err != nil {
		log.Fatal(err)
	}
	s := search{
		seen:     make(map[string]bool),
		host:     parsedUrl.Host,
		protocol: parsedUrl.Scheme,
	}
	s.dfs(start)
}
