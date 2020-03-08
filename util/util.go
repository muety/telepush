package util

import (
	"encoding/json"
	"log"
	"os"
)

func DumpJson(filePath string, data interface{}) {
	log.Println("Saving json.")
	file, err := os.Create(filePath)
	defer file.Close()
	if err != nil {
		log.Println(err)
	}

	encoder := json.NewEncoder(file)
	err = encoder.Encode(&data)
	if err != nil {
		log.Println(err)
	}
}
