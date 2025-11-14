package main

import (
	"strings"
)

func SlowStringConcat(strs []string) string {
	var result string

	for _, s := range strs {
		result += s
	}

	return result
}

func FastStringConcat(strs []string) string {
	var builder strings.Builder

	for _, s := range strs {
		builder.WriteString(s)
	}

	return builder.String()
}
