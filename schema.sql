CREATE TABLE pastes (id text, content text, time number);
CREATE TABLE logins (account text, token text);
CREATE TABLE accounts (id text, email text, password text);
CREATE TABLE bookmarks (account text, paste text, name text);
CREATE UNIQUE INDEX account_email on accounts(email);
