package utils

import (
	"strings"
)

// function to encode args for search query
func EncodeArgs(args string) string {
	parts := strings.Split(args, " ")
	return strings.Join(parts, "-")
}