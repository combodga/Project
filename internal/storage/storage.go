package storage

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"sync"
)

var (
	dbFile = ""
	pairs  = make(map[string]string)
	mutex  = &sync.RWMutex{}
)

func Init() {
	dbFile = os.Getenv("FILE_STORAGE_PATH")
	if dbFile == "" {
		return
	}

	pairsStr, err := ioutil.ReadFile(dbFile)
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

	if dbFile == "" {
		return nil
	}

	jsonStr, err := json.Marshal(pairs)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(dbFile, []byte(jsonStr), 0666)
	if err != nil {
		return err
	}

	return nil
}
