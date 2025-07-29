package main

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
