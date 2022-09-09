package storage

import (
	"encoding/json"
	"os"
	"sync"
)

type Storage struct {
	DBFile string
	Pairs map[string]string
	Mutex *sync.RWMutex
}

func New(dbFile string) *Storage {
	s := &Storage{
		DBFile: dbFile,
		Pairs: make(map[string]string),
		Mutex: &sync.RWMutex{},
	}

	if dbFile == "" {
		return s
	}

	s.Mutex.RLock()
	pairsStr, err := os.ReadFile(dbFile)
	if err != nil {
		s.Mutex.RUnlock()
		panic(err)
	}

	err = json.Unmarshal(pairsStr, &s.Pairs)
	s.Mutex.RUnlock()
	if err != nil {
		panic(err)
	}

	return s
}

func (s *Storage) GetURL(id string) (string, bool) {
	if len(id) <= 0 {
		return "", false
	}

	s.Mutex.RLock()
	url, ok := s.Pairs[id]
	s.Mutex.RUnlock()
	if !ok {
		return "", false
	}

	return url, true
}

func (s *Storage) SetURL(id, link string) error {
	s.Mutex.Lock()
	s.Pairs[id] = link

	if s.DBFile == "" {
		s.Mutex.Unlock()
		return nil
	}

	jsonStr, err := json.Marshal(s.Pairs)
	if err != nil {
		s.Mutex.Unlock()
		return err
	}

	err = os.WriteFile(s.DBFile, []byte(jsonStr), 0777)
	s.Mutex.Unlock()
	if err != nil {
		return err
	}

	return nil
}
