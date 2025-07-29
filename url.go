package main

// the function removes trailing and leading slashes.
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

// the function appends .html file extension if necessary.
func appendHtml(urlPath string) string {
	if urlPath == "" {
		return "index.html"
	}
	if len(urlPath) < 5 || urlPath[len(urlPath)-5:] != ".html" {
		return urlPath + ".html"
	}
	return urlPath
}
