package storage

import (
	"encoding/json"
	"io/ioutil"
	"sync"
)

var (
	DbFile = ""
	pairs  = make(map[string]string)
	mutex  = &sync.RWMutex{}
)

func Init() {
	if DbFile == "" {
		return
	}

	pairsStr, err := ioutil.ReadFile(DbFile)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(pairsStr, &pairs)
	if err != nil {
		panic(err)
	}
}

func GetURL(id string) (string, bool) {
	if len(id) <= 0 {
		return "", false
	}

	mutex.RLock()
	url, ok := pairs[id]
	mutex.RUnlock()
	if !ok {
		return "", false
	}

	return url, true
}

func SetURL(id, link string) error {
	mutex.Lock()
	pairs[id] = link
	mutex.Unlock()

	if DbFile == "" {
		return nil
	}

	jsonStr, err := json.Marshal(pairs)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(DbFile, []byte(jsonStr), 0777)
	if err != nil {
		return err
	}

	return nil
}
