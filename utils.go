package main

import "strings"

func cleanInput(text string) []string {
	text_lower := strings.ToLower(text)
	var fields []string
	fields = strings.Fields(text_lower)
	return fields
}
