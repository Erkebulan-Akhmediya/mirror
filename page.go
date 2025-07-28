package main

import (
	"bytes"
	"fmt"
	"golang.org/x/net/html"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type page struct {
	content []byte
	urlPath string
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
		return nil, fmt.Errorf("request failed with status %d: %s", res.StatusCode, res.Status)
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
	return &page{content: content, urlPath: urlPath}, nil
}

func (p *page) toTree() (*html.Node, error) {
	r := p.toReader()
	return html.Parse(r)
}

func (p *page) toReader() io.Reader {
	return bytes.NewReader(p.content)
}

func (p *page) save() error {
	if p.urlPath == "" {
		p.urlPath = "index.html"
	} else if p.urlPath[len(p.urlPath)-5:] != ".html" {
		p.urlPath += ".html"
	}
	p.saveTo(p.urlPath)
	return nil
}

func (p *page) saveTo(path string) {
	fmt.Println("Saving page:", path)
}
