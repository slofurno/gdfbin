package main

import "database/sql"

type Bookmark struct {
	Account string
	Paste   string
	Name    string
}

type bookmarkStore struct {
	*sql.DB
}

func (s *bookmarkStore) Insert(bookmark *Bookmark) error {
	_, err := s.DB.Exec("insert into bookmarks values (?,?,?)", bookmark.Account, bookmark.Paste, bookmark.Name)
	return err
}

func (s *bookmarkStore) Get(account *Account) ([]*Bookmark, error) {
	query := "SELECT paste, name FROM bookmarks WHERE bookmarks.account = ?"
	rows, err := s.DB.Query(query, account.Id)

	if err != nil {
		return nil, err
	}

	var bookmarks []*Bookmark

	for rows.Next() {
		bookmark := &Bookmark{}
		rows.Scan(&bookmark.Paste, &bookmark.Name)
		bookmarks = append(bookmarks, bookmark)
	}

	return bookmarks, nil
}
