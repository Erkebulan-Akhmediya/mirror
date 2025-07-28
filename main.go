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
	parentDir := "./" + parsedUrl.Host
	if err = os.Mkdir(parentDir, os.ModePerm); err != nil {
		log.Fatal(err)
	}
	if err = os.Chdir(parentDir); err != nil {
		log.Fatal(err)
	}
	s.dfs(start)
	if err = os.Chdir("../"); err != nil {
		log.Fatal(err)
	}
}
