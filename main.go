package main

import (
	"log"
	"net/url"
	"os"
)

func main() {
	parsedUrl, err := url.Parse(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	parentDir := "./" + parsedUrl.Host
	if err = os.Mkdir(parentDir, os.ModePerm); err != nil {
		log.Fatal(err)
	}
	if err = os.Chdir(parentDir); err != nil {
		log.Fatal(err)
	}
	urlPath := trimSlash(parsedUrl.Path)
	newSearch(parsedUrl).dfs(urlPath)
	if err = os.Chdir("../"); err != nil {
		log.Fatal(err)
	}
}
