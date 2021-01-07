package store

import (
	"encoding/gob"
	"github.com/muety/webhook2telegram/model"
	"github.com/orcaman/concurrent-map"
	"log"
	"os"
)

type GobStore struct {
	data     cmap.ConcurrentMap
	filePath string
}

func NewGobStore(filePath string) *GobStore {
	//gob.Register(model.StoreObject{})
	//gob.Register(model.StoreMessageObject{})

	// Backwards compatibility
	gob.RegisterName("main.StoreObject", model.StoreObject{})
	gob.RegisterName("main.StoreMessageObject", model.StoreMessageObject{})

	store := &GobStore{
		data:     cmap.New(),
		filePath: filePath,
	}

	if err := store.load(); err == nil {
		log.Println("read existing gob store from file")
	}

	return store
}

func (s *GobStore) load() error {
	file, err := os.Open(s.filePath)
	defer file.Close()
	if err != nil {
		log.Printf("error: failed to read store from %s\n", s.filePath)
		return nil
	}

	var rawData map[string]interface{}
	if err := gob.NewDecoder(file).Decode(&rawData); err != nil {
		log.Printf("error: failed to decode store data from %s (%v)\n", s.filePath, err)
		return nil
	}

	s.data = cmap.New()
	for k, v := range rawData {
		s.data.Set(k, v)
	}

	return nil
}

func (s *GobStore) dump() error {
	file, err := os.Create(s.filePath)
	defer file.Close()
	if err != nil {
		log.Printf("error: failed to dump store to %s (%v)", s.filePath, err)
		return err
	}

	return gob.NewEncoder(file).Encode(s.data.Items())
}

func (s *GobStore) Get(key string) interface{} {
	if v, ok := s.data.Get(key); ok {
		return v
	}
	return nil
}

func (s *GobStore) Put(key string, value interface{}) {
	s.data.Set(key, value)
	go s.dump()
}

func (s *GobStore) Delete(key string) {
	s.data.Remove(key)
	go s.dump()
}

func (s *GobStore) GetItems() map[string]interface{} {
	return s.data.Items()
}

func (s *GobStore) Flush() error {
	return s.dump()
}
