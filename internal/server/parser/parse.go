package parser

import "strings"

func ExtractIDFromURI(uri string) string {
	parts := strings.Split(uri, "://")

	if len(parts) == 2 {
		return parts[1]
	}

	return ""
}
