package storage

import (
	"encoding/json"
	"io/ioutil"

	"github.com/orcaman/concurrent-map/v2"
)

const (
	dbFile = "db"
)

var (
	pairs = cmap.New[string]()
)

func Init() {
  pairsStr, err := ioutil.ReadFile(dbFile)
  if err != nil {
    panic(err)
  }

  m := make(map[string]string)
  json.Unmarshal(pairsStr, &m)
  pairs.MSet(m)
}

func GetURL(id string) (string, bool) {
  if len(id) <= 0 {
    return "", false
  }

  url, ok := pairs.Get(id)
  if !ok {
    return "", false
  }

  return url, true
}

func SetURL(id, link string) error {
  pairs.Set(id, link)

  jsonBytes, err := pairs.MarshalJSON()
  if err != nil {
    return err
  }

  err = ioutil.WriteFile(dbFile, jsonBytes, 0666)
  if err != nil {
    return err
  }

  return nil
}
