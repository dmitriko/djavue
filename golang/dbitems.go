package main

import "golang.org/x/crypto/bcrypt"

type User struct {
	ID       string
	Username string
	Password string
	Token    string
}

func (dbw *DBWorker) CreateUserTable() error {
	return dbw.WriteOne("create table users (id text primary key, username text unique, password text, token text unique)")
}

func (dbw *DBWorker) SaveNewUser(user *User) error {
	return dbw.WriteOne("insert into users (id, username, password, token) values (?,?,?,?)",
		user.ID, user.Username, user.Password, user.Token)
}

func NewUser(username, password string) (*User, error) {
	token, err := NewToken()
	if err != nil {
		return nil, err
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		return nil, err
	}
	id, err := NewULIDNow()
	if err != nil {
		return nil, err
	}
	user := &User{ID: id, Username: username, Password: string(hash), Token: token}
	return user, nil
}

func (u *User) IsPasswdValid(passwd string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(passwd))
	if err != nil {
		return false
	}
	return true
}
