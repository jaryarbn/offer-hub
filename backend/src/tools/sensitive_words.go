package tools

import "strings"

var sensitiveWords = []string{
	"敏感词",
	"操你妈",
	"傻逼",
	"妈的",
}

// FilterSensitiveWords masks configured sensitive words while preserving the
// surrounding comment text.
func FilterSensitiveWords(content string) string {
	filtered := content
	for _, word := range sensitiveWords {
		filtered = strings.ReplaceAll(filtered, word, strings.Repeat("*", len([]rune(word))))
	}
	return filtered
}
