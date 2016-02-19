package main

import "database/sql"

type DataStore struct {
	Pastes *pasteStore
}

func NewDataStore(db *sql.DB) *DataStore {
	store := &DataStore{}
	store.Pastes = &pasteStore{DB: db}
	return store
}
