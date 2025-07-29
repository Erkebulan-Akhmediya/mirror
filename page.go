package main

import (
	"bytes"
	"fmt"
	"golang.org/x/net/html"
	"io"
	"net/url"
	"os"
	"strings"
)

type page struct {
	content                 []byte
	urlPath, host, protocol string
}

func (p *page) toTree() (*html.Node, error) {
	r := p.toReader()
	return html.Parse(r)
}

func (p *page) toReader() io.Reader {
	return bytes.NewReader(p.content)
}

func (p *page) save() error {
	p.urlPath = appendHtml(p.urlPath)
	return p.saveTo(p.urlPath)
}

func (p *page) saveTo(path string) error {
	fmt.Println("downloading to: ./" + path)
	lastSlash := strings.LastIndex(path, "/")
	if lastSlash == -1 {
		return os.WriteFile(path, []byte(p.content), os.ModePerm)
	}
	dir := "./" + path[:lastSlash]
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(p.content), os.ModePerm)
}

func (p *page) urlPaths() ([]string, error) {
	tree, err := p.toTree()
	if err != nil {
		return nil, fmt.Errorf("error parsing page: %v", err)
	}
	return p.urlPathsInTree(tree), nil
}

func (p *page) urlPathsInTree(n *html.Node) []string {
	var urlPaths []string
	if urlPath, found := p.urlPathInNode(n); found {
		urlPaths = append(urlPaths, urlPath)
	}
	for child := n.FirstChild; child != nil; child = child.NextSibling {
		urlPaths = append(urlPaths, p.urlPathsInTree(child)...)
	}
	return urlPaths
}

func (p *page) urlPathInNode(n *html.Node) (string, bool) {
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
		if parsedUrl.Host == p.host || parsedUrl.Host == "" {
			return trimSlash(parsedUrl.Path), true
		}
	}
	return "", false
}
