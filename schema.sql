CREATE TABLE pastes (id text, content text);
CREATE TABLE logins (id text, account text);
CREATE TABLE accounts (id text, email text, password text, salt text);
CREATE TABLE owned_pastes (paste_id text, account_id text, title text);

