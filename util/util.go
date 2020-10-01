package util

import (
	"encoding/json"
	"log"
	"os"
	"strings"
)

var markdownReplacer = strings.NewReplacer(
	"\\", "\\\\",
	"`", "\\`",
	"*", "\\*",
	"[", "\\[",
	"_", "\\_",
)

func DumpJson(filePath string, data interface{}) {
	log.Println("Saving json.")
	file, err := os.Create(filePath)
	defer file.Close()
	if err != nil {
		log.Println(err)
	}
	if err := json.NewEncoder(file).Encode(&data); err != nil {
		log.Println(err)
	}
}

func EscapeMarkdown(s string) string {
	return markdownReplacer.Replace(s)
}
