package main

import (
	"fmt"
	"golang.org/x/net/html"
	"log"
	"net/url"
)

type search struct {
	seen     map[string]bool
	host     string
	protocol string
}

func (s *search) dfs(fullUrl string) {
	urlPath, err := trimToPath(fullUrl)
	if err != nil {
		return
	}
	if s.seen[urlPath] {
		return
	}
	s.seen[urlPath] = true
	p, err := s.fetchPage(urlPath)
	if err != nil {
		log.Println("Error fetching page:", err)
		return
	}
	if err = p.save(); err != nil {
		log.Println("Error saving page:", err)
	}
	tree, err := p.toTree()
	if err != nil {
		log.Println("Error parsing page:", err)
		return
	}
	links := s.findLinks(tree)
	for _, link := range links {
		s.dfs(link)
	}
}

func (s *search) findLinks(node *html.Node) []string {
	res := s.findLink(node)
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		res = append(res, s.findLinks(child)...)
	}
	return res
}

func (s *search) findLink(node *html.Node) []string {
	if node.Type != html.ElementNode || node.Data != "a" {
		return nil
	}
	var res []string
	for _, a := range node.Attr {
		if a.Key != "href" {
			continue
		}
		parsedUrl, err := url.Parse(a.Val)
		if err != nil {
			fmt.Println("Error parsing url:", err)
			break
		}
		if parsedUrl.Host == s.host {
			res = append(res, parsedUrl.String())
			break
		}
		if parsedUrl.Host == "" {
			parsedUrl.Host = s.host
			parsedUrl.Scheme = s.protocol
			res = append(res, parsedUrl.String())
			break
		}
	}
	return res
}
