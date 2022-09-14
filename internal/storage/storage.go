package storage

import (
	"database/sql"
	"encoding/json"
	"errors"
	"os"
	"sync"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Storage struct {
	DBFile        string
	DBCredentials string
	Pairs         map[string]map[string]string
	Mutex         *sync.RWMutex
}

type Link struct {
	User string `db:"usr"`
	ID   string `db:"short"`
	Link string `db:"long"`
}

func New(dbFile, dbCredentials string) (*Storage, error) {
	s := &Storage{
		DBFile:        dbFile,
		DBCredentials: dbCredentials,
		Pairs:         make(map[string]map[string]string),
		Mutex:         &sync.RWMutex{},
	}

	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	if dbCredentials != "" {
		db, err := sqlx.Connect("postgres", s.DBCredentials)
		if err != nil {
			return s, err
		}
		defer db.Close()

		db.MustExec(`
			CREATE TABLE IF NOT EXISTS shortener (
				usr text,
				short text,
				long text
			);
		`)

		link := Link{}
		rows, err := db.Queryx("SELECT * FROM shortener")
		if err != nil {
			return s, err
		}
		for rows.Next() {
			err := rows.StructScan(&link)
			if err != nil {
				return s, err
			}
			if len(s.Pairs[link.User]) == 0 {
				s.Pairs[link.User] = make(map[string]string)
			}
			s.Pairs[link.User][link.ID] = link.Link
		}
		err = rows.Err()
		if err != nil {
			return s, nil
		}

		return s, nil
	}

	if dbFile != "" {
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

	if s.DBCredentials != "" {
		db, err := sqlx.Connect("postgres", s.DBCredentials)
		if err != nil {
			return err
		}
		defer db.Close()

		db.MustExec("INSERT INTO shortener VALUES ($1, $2, $3)", user, id, link)

		return nil
	}

	if s.DBFile != "" {
		jsonStr, err := json.Marshal(s.Pairs)
		if err != nil {
			return err
		}

		err = os.WriteFile(s.DBFile, []byte(jsonStr), 0777)
		if err != nil {
			return err
		}
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

func (s *Storage) Ping() bool {
	db, err := sql.Open("postgres", s.DBCredentials)
	if err != nil {
		return false
	}
	db.Close()
	return true
}
