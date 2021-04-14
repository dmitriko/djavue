package main

import "golang.org/x/crypto/bcrypt"

type User struct {
	ID       string
	Username string
	Password string
	Token    string
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
