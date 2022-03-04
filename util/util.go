package util

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
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

func RandomString(length int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, length)
	rand.Read(b)
	return fmt.Sprintf("%x", b)[:length]
}
