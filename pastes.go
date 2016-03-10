package main

import (
	"database/sql"
	"github.com/bradfitz/gomemcache/memcache"
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

type PasteStore interface {
	Insert(*Paste) error
	Get(string) (*Paste, error)
}

type pasteCache struct {
	mc *memcache.Client
}

func NewPasteCache() *pasteCache {
	mc := memcache.New("localhost:11211")
	return &pasteCache{mc: mc}
}

func (s *pasteCache) Insert(paste *Paste) error {
	item := &memcache.Item{
		Key:   paste.Id,
		Value: paste.Content,
	}

	return s.mc.Add(item)
}

func (s *pasteCache) Get(id string) (*Paste, error) {
	item, err := s.mc.Get(id)

	if err != nil {
		return nil, err
	}

	paste := &Paste{
		Id:      item.Key,
		Content: item.Value,
	}

	return paste, nil
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
