package main

import (
	"database/sql"
	"time"
)

type Paste struct {
	id      string
	content []byte
	time    int64
}

func NewPaste() *Paste {
	return &Paste{
		id:   randomHash(),
		time: epoch_ms(),
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
	_, err := s.DB.Exec("insert into pastes values (?,?,?)", paste.id, paste.content, paste.time)
	return err
}

func (s *pasteStore) Get(id string) (*Paste, error) {
	rows, err := s.DB.Query("select id, content, time from pastes WHERE id = ?", id)

	if err != nil {
		return nil, err
	}

	paste := &Paste{}

	for rows.Next() {
		rows.Scan(&paste.id, &paste.content, &paste.time)
	}

	return paste, nil
}
