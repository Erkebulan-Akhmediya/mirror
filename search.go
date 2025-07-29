package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type search struct {
	seen     map[string]bool
	host     string
	protocol string
}

func newSearch(parsedUrl *url.URL) *search {
	return &search{
		seen:     make(map[string]bool),
		host:     parsedUrl.Host,
		protocol: parsedUrl.Scheme,
	}
}

func (s *search) dfs(urlPath string) {
	if s.seen[urlPath] {
		return
	}
	s.seen[urlPath] = true
	p, err := s.fetchPage(urlPath)
	if err != nil {
		log.Println("Error fetching page:", err)
		return
	}
	urlPaths, err := p.urlPaths()
	if err != nil {
		log.Println("Error finding url paths:", err)
		return
	}
	for _, link := range urlPaths {
		s.dfs(link)
	}
	if err = p.save(); err != nil {
		log.Println("Error saving page:", err)
	}
}

func (s *search) fetchPage(urlPath string) (*page, error) {
	fullUrl := (&url.URL{Path: urlPath, Scheme: s.protocol, Host: s.host}).String()
	res, err := http.Get(fullUrl)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = res.Body.Close()
	}()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed to %q with status %d: %s", fullUrl, res.StatusCode, res.Status)
	}

	cth := res.Header.Get("Content-Type")
	ct := strings.Split(cth, ";")[0]
	if ct != "text/html" {
		return nil, fmt.Errorf("unsupported content type: %s", ct)
	}

	content, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return &page{content: content, urlPath: urlPath, host: s.host, protocol: s.protocol}, nil
}
