package main

import "database/sql"

type DataStore struct {
	Pastes    *pasteStore
	Accounts  *accountStore
	Bookmarks *bookmarkStore
	Logins    *loginStore
}

func NewDataStore(db *sql.DB) *DataStore {
	store := &DataStore{}
	store.Pastes = &pasteStore{DB: db}
	store.Accounts = &accountStore{DB: db}
	store.Logins = &loginStore{DB: db}
	store.Bookmarks = &bookmarkStore{DB: db}
	return store
}
