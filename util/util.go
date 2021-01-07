package util

import (
	"strings"
)

var markdownReplacer = strings.NewReplacer(
	"\\", "\\\\",
	"`", "\\`",
	"*", "\\*",
	"[", "\\[",
	"_", "\\_",
)

func EscapeMarkdown(s string) string {
	return markdownReplacer.Replace(s)
}
