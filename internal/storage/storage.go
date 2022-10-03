package storage

import (
	"database/sql"
	"encoding/json"
	"errors"
	"os"
	"sync"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type Storage struct {
	DBFile        string
	DBCredentials string
	Pairs         map[string]map[string]string
	Mutex         *sync.RWMutex
	ErrDupKey     error
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
		ErrDupKey:     fmt.Error("duplicate key"),
	}

	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	if dbCredentials != "" {
		db, err := sqlx.Connect("postgres", s.DBCredentials)
		if err != nil {
			return s, fmt.Errorf("db connect: %w", err)
		}
		defer db.Close()

		db.MustExec(`
			CREATE TABLE IF NOT EXISTS shortener (
				usr text,
				short text unique,
				long text
			);
		`)

		link := Link{}
		rows, err := db.Queryx("SELECT * FROM shortener")
		if err != nil {
			return s, fmt.Errorf("read rows: %w", err)
		}
		defer rows.Close()
		for rows.Next() {
			err := rows.StructScan(&link)
			if err != nil {
				return s, fmt.Errorf("rows struct scan: %w", err)
			}
			if len(s.Pairs[link.User]) == 0 {
				s.Pairs[link.User] = make(map[string]string)
			}
			s.Pairs[link.User][link.ID] = link.Link
		}
		err = rows.Err()
		if err != nil {
			return s, fmt.Errorf("rows error: %w", err)
		}

		return s, nil
	}

	if dbFile != "" {
		pairsStr, err := os.ReadFile(dbFile)
		if errors.Is(err, os.ErrNotExist) {
			return s, nil
		}
		if err != nil {
			return s, fmt.Errorf("read file: %w", err)
		}

		err = json.Unmarshal(pairsStr, &s.Pairs)
		if err != nil {
			return s, fmt.Errorf("json unmarshal: %w", err)
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
			return fmt.Errorf("sql connect: %w", err)
		}
		defer db.Close()

		_, err = db.Exec("INSERT INTO shortener VALUES ($1, $2, $3)", user, id, link)
		if err != nil {
			if err, ok := err.(*pq.Error); ok {
			    if err.Code == "23505" {
			    	return s.ErrDupKey
			    }
			}
		}

		if err != nil {
			return fmt.Errorf("db error: %w", err)
		}
		return nil
	}

	if s.DBFile != "" {
		jsonStr, err := json.Marshal(s.Pairs)
		if err != nil {
			return fmt.Errorf("json marshal: %w", err)
		}

		err = os.WriteFile(s.DBFile, []byte(jsonStr), 0777)
		if err != nil {
			return fmt.Errorf("write file: %w", err)
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
	defer db.Close()

    if err = db.Ping(); err != nil {
        return false
    }

	return true
}
