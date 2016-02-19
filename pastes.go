package main

import (
	"database/sql"
	"time"
)

type Paste struct {
	Id      string
	Content []byte
	Time    int64
}

func NewPaste() *Paste {
	return &Paste{
		Id:   randomHash(5),
		Time: epoch_ms(),
	}
}

func epoch_ms() int64 {
	now := time.Now()
	return now.UnixNano() / 1000000
}

type pasteStore struct {
	*sql.DB
}

func (s *pasteStore) Insert(paste *Paste) error {
	_, err := s.DB.Exec("insert into pastes values (?,?,?)", paste.Id, paste.Content, paste.Time)
	return err
}

func (s *pasteStore) Get(id string) (*Paste, error) {
	rows, err := s.DB.Query("select id, content, time from pastes WHERE id = ?", id)

	if err != nil {
		return nil, err
	}

	paste := &Paste{}

	for rows.Next() {
		rows.Scan(&paste.Id, &paste.Content, &paste.Time)
	}

	return paste, nil
}
