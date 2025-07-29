package main

import (
	"bytes"
	"fmt"
	"golang.org/x/net/html"
	"io"
	"os"
	"strings"
)

type page struct {
	content                 []byte
	urlPath, host, protocol string
}

func (p *page) toTree() (*node, error) {
	r := p.toReader()
	root, err := html.Parse(r)
	if err != nil {
		return nil, err
	}
	domRoot := node(*root)
	return &domRoot, nil
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
	return p.saveTo(p.urlPath)
}

func (p *page) saveTo(path string) error {
	fmt.Println("downloading to: ./" + path)
	lastSlash := strings.LastIndex(path, "/")
	if lastSlash == -1 {
		return os.WriteFile(path, p.content, os.ModePerm)
	}
	dir := "./" + path[:lastSlash]
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}
	return os.WriteFile(path, p.content, os.ModePerm)
}

func (p *page) urlPaths() ([]string, error) {
	tree, err := p.toTree()
	if err != nil {
		return nil, fmt.Errorf("error parsing page: %v", err)
	}
	return p.urlPathsInTree(tree), nil
}

func (p *page) urlPathsInTree(n *node) []string {
	var urlPaths []string
	if urlPath, found := n.urlPath(p.host); found {
		urlPaths = append(urlPaths, urlPath)
	}
	for child := n.FirstChild; child != nil; child = child.NextSibling {
		domChild := node(*child)
		urlPaths = append(urlPaths, p.urlPathsInTree(&domChild)...)
	}
	return urlPaths
}
