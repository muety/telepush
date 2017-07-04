package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

var store map[string]StoreObject

func InitEmpty() {
	store = make(map[string]StoreObject)
}

func ReadStoreFromJSON(filePath string) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Println("Could not read store from file. Initializing empty one.")
		InitEmpty()
		return
	}
	err = json.Unmarshal(data, &store)
	if err != nil {
		fmt.Println("Could not read store from file. Initializing empty one.")
		InitEmpty()
		return
	}
}

func FlushStoreToJSON(filePath string) error {
	data, err := json.Marshal(&store)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filePath, []byte(data), 0644)
}

func StoreGet(key string) StoreObject {
	return store[key]
}

func StorePut(key string, value StoreObject) {
	store[key] = value
}
