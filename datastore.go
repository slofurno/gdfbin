package main

import "database/sql"

type DataStore struct {
	Pastes    *pasteStore
	Accounts  *accountStore
	Bookmarks *bookmarkStore
	Logins    *loginStore
	Pastes2   []PasteStore
}

func NewDataStore(db *sql.DB) *DataStore {
	store := &DataStore{}
	pastestore := &pasteStore{DB: db}
	pastecache := NewPasteCache()

	store.Pastes = pastestore
	store.Pastes2 = []PasteStore{pastecache, pastestore}

	store.Accounts = &accountStore{DB: db}
	store.Logins = &loginStore{DB: db}
	store.Bookmarks = &bookmarkStore{DB: db}
	return store
}
