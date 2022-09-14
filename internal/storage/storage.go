package storage

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
)

type Storage struct {
	DBFile string
	Pairs  map[string]map[string]string
	Mutex  *sync.RWMutex
}

func New(dbFile string) (*Storage, error) {
	s := &Storage{
		DBFile: dbFile,
		Pairs:  make(map[string]map[string]string),
		Mutex:  &sync.RWMutex{},
	}

	if dbFile == "" {
		return s, nil
	}

	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	pairsStr, err := os.ReadFile(dbFile)
	if errors.Is(err, os.ErrNotExist) {
		return s, nil
	}
	if err != nil {
		return s, err
	}

	err = json.Unmarshal(pairsStr, &s.Pairs)
	if err != nil {
		return s, err
	}

	return s, nil
}

func (s *Storage) GetURL(id string) (string, bool) {
	if len(id) <= 0 {
		return "", false
	}

	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	for user := range s.Pairs {
		url, ok := s.Pairs[user][id]
		if ok {
			return url, true
		}
	}

	return "", false
}

func (s *Storage) SetURL(user, id, link string) error {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	if len(s.Pairs[user]) == 0 {
		s.Pairs[user] = make(map[string]string)
	}
	s.Pairs[user][id] = link

	if s.DBFile == "" {
		return nil
	}

	jsonStr, err := json.Marshal(s.Pairs)
	if err != nil {
		return err
	}

	err = os.WriteFile(s.DBFile, []byte(jsonStr), 0777)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) ListURL(user string) (map[string]string, bool) {
	s.Mutex.Lock()
	list, ok := s.Pairs[user]
	s.Mutex.Unlock()
	if !ok {
		return list, false
	}

	return list, true
}
