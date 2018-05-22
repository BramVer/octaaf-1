package main

import "strings"

type style rune

const (
	mdbold    style = '*'
	mdcursive style = '_'
	mdquote   style = '`'
	mdescape  bool  = true
)

const tokenList = "\\*_`"

// Markdown ensures safe markdown parsing on unsafe input
func Markdown(input string, style style) string {
	return string(style) + MDEscape(input) + string(style)
}

// MDEscape escapes all telegram markdown characters
func MDEscape(input string) string {
	for _, token := range tokenList {
		input = strings.Replace(input, string(token), "\\"+string(token), -1)
	}
	return input
}
