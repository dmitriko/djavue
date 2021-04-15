package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewUser(t *testing.T) {
	user, err := NewUser("foo", "bar")
	assert.Nil(t, err)
	assert.Equal(t, "foo", user.Username)
	assert.NotEqual(t, "bar", user.Password)
}

func TestUserPassword(t *testing.T) {
	user, err := NewUser("foo", "bar")
	assert.Nil(t, err)
	assert.True(t, user.IsPasswdValid("bar"))
	assert.False(t, user.IsPasswdValid("dummy"))
}

func TestUserNewSave(t *testing.T) {
	dbw, err := testDBWorker()
	if assert.Nil(t, err) {
		defer removeWorker(dbw)
	}
	assert.Nil(t, dbw.CreateUserTable())
	user, err := NewUser("foo", "bar")
	assert.Nil(t, err)
	assert.Nil(t, dbw.SaveNewUser(user))
	var u User
	assert.Nil(t, dbw.QueryRow("select id, username, password, token from users where id=?", user.ID).Scan(
		&u.ID, &u.Username, &u.Password, &u.Token))
	assert.True(t, u.IsPasswdValid("bar"))
	assert.Equal(t, "foo", u.Username)
}

func TestUserLoad(t *testing.T) {
	dbw, err := testDBWorker()
	if assert.Nil(t, err) {
		defer removeWorker(dbw)
	}
	assert.Nil(t, dbw.CreateUserTable())
	user1, err := NewUser("foo", "bar")
	assert.Nil(t, err)
	assert.Nil(t, dbw.SaveNewUser(user1))
	user2, err := NewUser("spam", "egg")
	assert.Nil(t, err)
	assert.Nil(t, dbw.SaveNewUser(user2))
	u, err := dbw.LoadUser(user2.ID)
	assert.Nil(t, err)
	assert.Equal(t, "spam", u.Username)
}

func TestUserUpdate(t *testing.T) {
	dbw, err := testDBWorker()
	if assert.Nil(t, err) {
		defer removeWorker(dbw)
	}
	assert.Nil(t, dbw.CreateUserTable())
	user, err := NewUser("foo", "bar")
	assert.Nil(t, err)
	assert.Nil(t, dbw.SaveNewUser(user))
	token, _ := NewToken()
	user.Token = token
	assert.Nil(t, dbw.SaveUser(user))
	u, err := dbw.LoadUser(user.ID)
	assert.Nil(t, err)
	assert.Equal(t, token, u.Token)

}
