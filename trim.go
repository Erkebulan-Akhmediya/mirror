package main

import (
	"net/url"
	"strings"
)

func trimToPath(urlStr string) (string, error) {
	urlStr = strings.TrimSpace(urlStr)
	parsedUrl, err := url.Parse(urlStr)
	if err != nil {
		return "", err
	}
	return trimSlash(parsedUrl.Path), nil
}

func trimSlash(s string) string {
	if s == "/" || s == "" {
		return ""
	}
	if s[0] == '/' {
		s = s[1:]
	}
	if s[len(s)-1] == '/' {
		s = s[:len(s)-1]
	}
	return s
}
