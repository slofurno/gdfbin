package main

import (
	"database/sql"
	"fmt"
)

type Bookmark struct {
	Account string
	Paste   string
	Name    string
	Time    int64
}

type bookmarkStore struct {
	*sql.DB
}

func (s *bookmarkStore) Insert(bookmark *Bookmark) error {
	_, err := s.DB.Exec("insert into bookmarks values (?,?,?)", bookmark.Account, bookmark.Paste, bookmark.Name)
	return err
}

func (s *bookmarkStore) Remove(bookmark *Bookmark) error {
	query := `
	DELETE from bookmarks
	WHERE bookmarks.account = ?
	AND bookmarks.name = ?`

	_, err := s.DB.Exec(query, bookmark.Account, bookmark.Name)

	return err
}

func (s *bookmarkStore) GetHistory(bookmark *Bookmark) []*Paste {
	query := `
	SELECT pastes.id, pastes.time
	FROM bookmarks
	INNER JOIN pastes
	ON bookmarks.paste = pastes.id
	WHERE bookmarks.account = ?
	AND bookmarks.name = ?
	ORDER BY time DESC`

	rows, err := s.DB.Query(query, bookmark.Account, bookmark.Name)
	pastes := []*Paste{}

	if err != nil {
		return pastes
	}

	for rows.Next() {
		paste := &Paste{}
		err := rows.Scan(&paste.Id, &paste.Time)
		if err == nil {
			pastes = append(pastes, paste)
		}
	}
	return pastes
}

func (s *bookmarkStore) GetPaste(bookmark *Bookmark) *Paste {
	query := `
	SELECT pastes.content FROM bookmarks
	INNER JOIN pastes
	ON bookmarks.paste = pastes.id
	WHERE bookmarks.account = ?
	AND bookmarks.name = ?
	ORDER BY time DESC
	LIMIT 1`

	row := s.DB.QueryRow(query, bookmark.Account, bookmark.Name)
	paste := &Paste{}

	err := row.Scan(&paste.Content)

	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	return paste
}

func (s *bookmarkStore) Get(account *Account) ([]*Bookmark, error) {
	/*
		query := `
		SELECT * FROM
		(SELECT bookmarks.paste, bookmarks.name, pastes.time
		FROM pastes
		INNER JOIN bookmarks
		ON pastes.id = bookmarks.paste
		WHERE bookmarks.account = ?
		ORDER BY pastes.time ASC) as x
		GROUP BY x.name`
	*/

	query := `
	SELECT bookmarks.paste, bookmarks.name, pastes.time
	FROM pastes
	INNER JOIN bookmarks
	ON pastes.id = bookmarks.paste
	WHERE bookmarks.account = ?
	ORDER BY pastes.time DESC`

	rows, err := s.DB.Query(query, account.Id)

	if err != nil {
		return nil, err
	}

	var bookmarks []*Bookmark
	seen := make(map[string]bool)

	for rows.Next() {
		bookmark := &Bookmark{}
		rows.Scan(&bookmark.Paste, &bookmark.Name, &bookmark.Time)

		if seen[bookmark.Name] == false {
			bookmarks = append(bookmarks, bookmark)
			seen[bookmark.Name] = true
		}
	}

	return bookmarks, nil
}
