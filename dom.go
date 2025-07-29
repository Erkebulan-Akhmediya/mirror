package main

import (
	"fmt"
	"golang.org/x/net/html"
	"net/url"
)

type node html.Node

func (n *node) urlPath(host string) (string, bool) {
	if n.Type != html.ElementNode || n.Data != "a" {
		return "", false
	}
	for _, a := range n.Attr {
		if a.Key != "href" {
			continue
		}
		parsedUrl, err := url.Parse(a.Val)
		if err != nil {
			fmt.Println("Error parsing url:", err)
			return "", false
		}
		if parsedUrl.Host == host || parsedUrl.Host == "" {
			return trimSlash(parsedUrl.Path), true
		}
	}
	return "", false
}
