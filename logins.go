package main

import "database/sql"

type Login struct {
	Account string
	Token   string
}

func NewLogin(account *Account) *Login {
	return &Login{
		Account: account.Id,
		Token:   randomHash(10),
	}
}

type loginStore struct {
	*sql.DB
}

func (s *loginStore) Insert(login *Login) error {
	_, err := s.DB.Exec("insert into logins values (?,?)", login.Account, login.Token)
	return err
}

func (s *loginStore) GetAccount(token string) (*Account, error) {
	command := "select accounts.id, accounts.email from accounts INNER JOIN logins ON accounts.id = logins.account WHERE logins.token = ?"
	rows, err := s.DB.Query(command, token)

	if err != nil {
		return nil, err
	}

	var account *Account

	for rows.Next() {
		account = &Account{}
		rows.Scan(&account.Id, &account.Email)
	}

	return account, nil
}
