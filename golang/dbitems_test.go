package main

import (
	"fmt"
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
	var u2 User
	err = dbw.LoadUser(&u2, user2.ID)
	assert.Nil(t, err)
	assert.Equal(t, "spam", u2.Username)
	var u1 User
	err = dbw.LoadUserByName(&u1, user1.Username)
	assert.Nil(t, err)
	assert.Equal(t, u1.ID, user1.ID)
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
	var u User
	err = dbw.LoadUser(&u, user.ID)
	assert.Nil(t, err)
	assert.Equal(t, token, u.Token)

}

func TestJobNew(t *testing.T) {
	user, _ := NewUser("foo", "bar")
	job, err := NewJob(user.ID, JOB_ORIG)
	assert.Nil(t, err)
	assert.Equal(t, job.UserID, user.ID)
}

func TestJobNewSave(t *testing.T) {
	dbw, err := testDBWorker()
	if assert.Nil(t, err) {
		defer removeWorker(dbw)
	}
	_ = dbw.CreateUserTable()
	_ = dbw.CreateJobTable()
	user, _ := NewUser("foo", "bar")
	job, _ := NewJob(user.ID, JOB_ALL_THREE)
	dbw.SaveNewUser(user)
	dbw.SaveNewJob(job)
	var job_id, user_id, state string
	assert.Nil(t, dbw.QueryRow("select id, user_id, state from jobs where id = ?", job.ID).Scan(&job_id, &user_id, &state))
}

func TestJobSave(t *testing.T) {
	dbw, err := testDBWorker()
	if assert.Nil(t, err) {
		defer removeWorker(dbw)
	}
	_ = dbw.CreateUserTable()
	_ = dbw.CreateJobTable()
	user, _ := NewUser("foo", "bar")
	job, _ := NewJob(user.ID, "original")
	dbw.SaveNewUser(user)
	dbw.SaveNewJob(job)
	assert.Equal(t, int64(0), job.State)
	job.State = 1
	assert.Nil(t, dbw.SaveJob(job))
	var job_state int64
	assert.Nil(t, dbw.QueryRow("select state from jobs where id = ?", job.ID).Scan(&job_state))
	assert.Equal(t, int64(1), job_state)
}

func TestJobLoad(t *testing.T) {
	dbw, err := testDBWorker()
	if assert.Nil(t, err) {
		defer removeWorker(dbw)
	}
	_ = dbw.CreateUserTable()
	_ = dbw.CreateJobTable()
	user, _ := NewUser("foo", "bar")
	job, _ := NewJob(user.ID, JOB_SQUARE_ORIG)
	dbw.SaveNewUser(user)
	dbw.SaveNewJob(job)
	var j Job
	err = dbw.LoadJob(&j, job.ID)
	assert.Nil(t, err)
	assert.Equal(t, user.ID, j.UserID)
}

func TestImageNew(t *testing.T) {
	user, _ := NewUser("foo", "bar")
	job, _ := NewJob(user.ID, JOB_ORIG)
	img, _ := NewImage(job, "/tmp", "foo.jpg", "image/jpeg")
	assert.Equal(t, fmt.Sprintf("/tmp/%s_foo.jpg", img.ID), img.Path)
}

func TestImageSaveNew(t *testing.T) {
	dbw, err := testDBWorker()
	if assert.Nil(t, err) {
		defer removeWorker(dbw)
	}
	_ = dbw.CreateUserTable()
	_ = dbw.CreateJobTable()
	assert.Nil(t, dbw.CreateImageTable())
	user, _ := NewUser("foo", "bar")
	job, _ := NewJob(user.ID, JOB_ORIG)
	img, _ := NewImage(job, "/tmp/", "foo.jpg", "image/jpeg")
	dbw.SaveNewUser(user)
	dbw.SaveNewJob(job)
	assert.Nil(t, dbw.SaveNewImage(img))
	var path string
	dbw.QueryRow("select path from images where id=?", img.ID).Scan(&path)
	assert.Equal(t, img.Path, path)
}

func TestImageLoad(t *testing.T) {
	dbw, err := testDBWorker()
	if assert.Nil(t, err) {
		defer removeWorker(dbw)
	}
	_ = dbw.CreateUserTable()
	_ = dbw.CreateJobTable()
	assert.Nil(t, dbw.CreateImageTable())
	user, _ := NewUser("foo", "bar")
	job, _ := NewJob(user.ID, JOB_ORIG)
	img, _ := NewImage(job, "/tmp/", "foo.jpg", "image/jpeg")
	dbw.SaveNewUser(user)
	dbw.SaveNewJob(job)
	assert.Nil(t, dbw.SaveNewImage(img))
	var imgLoaded Image
	err = dbw.LoadImage(&imgLoaded, img.ID)
	assert.Nil(t, err)
	assert.Equal(t, img.ID, imgLoaded.ID)
}
