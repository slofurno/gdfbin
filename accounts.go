package main

import (
	"database/sql"

	"golang.org/x/crypto/bcrypt"
)

type Account struct {
	Id       string
	Email    string
	Password []byte
}

type accountStore struct {
	*sql.DB
}

func NewAccount(email string, password string) (*Account, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return nil, err
	}

	account := &Account{
		Id:       randomHash(10),
		Email:    email,
		Password: hashed,
	}

	return account, nil
}

func (s *accountStore) Insert(account *Account) error {
	command := "insert into accounts values (?,?,?)"
	_, err := s.DB.Exec(command, account.Id, account.Email, account.Password)
	return err
}
