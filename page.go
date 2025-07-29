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

// the page struct represents an HTML page.
// it stores content of the page as well as its meta-data.
type page struct {
	content                 []byte
	urlPath, host, protocol string
}

// converts the page to its tree representation.
func (p *page) toTree() (*html.Node, error) {
	r := p.toReader()
	return html.Parse(r)
}

// creates a reader from the page.
func (p *page) toReader() io.Reader {
	return bytes.NewReader(p.content)
}

// the method saves the page using its url path as a file path.
func (p *page) save() error {
	p.urlPath = appendHtml(p.urlPath)
	return p.saveTo(p.urlPath)
}

// the method saves a file under the path creating directories if necessary.
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

// the method extracts url paths from the page.
func (p *page) urlPaths() ([]string, error) {
	tree, err := p.toTree()
	if err != nil {
		return nil, fmt.Errorf("error parsing page: %v", err)
	}
	return p.urlPathsInTree(tree), nil
}

// the method recursively visits each node to extract all url paths from the tree.
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

// the method extracts file path from the node.
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
