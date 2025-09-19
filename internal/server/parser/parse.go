package parser

import "strings"

func ExtractIDFromURI(uri string) string {
	parts := strings.Split(uri, "://")

	if len(parts) == 2 {
		return parts[1]
	}

	return ""
}

func CheckForProtectedTag(data map[string]any) bool {
	if tagNames, ok := data["tag_names"].([]any); ok {
		for _, tag := range tagNames {
			if tagStr, ok := tag.(string); ok && tagStr == "protected" {
				return true
			}
		}
	}

	return false
}
