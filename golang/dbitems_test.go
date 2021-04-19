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

func TestJobNew(t *testing.T) {
	user, _ := NewUser("foo", "bar")
	job, err := NewJob(user.ID)
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
	job, _ := NewJob(user.ID)
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
	job, _ := NewJob(user.ID)
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
	job, _ := NewJob(user.ID)
	dbw.SaveNewUser(user)
	dbw.SaveNewJob(job)
	j, err := dbw.LoadJob(job.ID)
	assert.Nil(t, err)
	assert.Equal(t, user.ID, j.UserID)
}

func TestImageNew(t *testing.T) {
	user, _ := NewUser("foo", "bar")
	job, _ := NewJob(user.ID)
	img, _ := NewImage(job, "foo.jpg", "image/jpeg")
	assert.Equal(t, "foo.jpg", img.Path)
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
	job, _ := NewJob(user.ID)
	img, _ := NewImage(job, "/tmp/foo.jpg", "image/jpeg")
	dbw.SaveNewUser(user)
	dbw.SaveNewJob(job)
	assert.Nil(t, dbw.SaveNewImage(img))
	var path string
	dbw.QueryRow("select path from images where id=?", img.ID).Scan(&path)
	assert.Equal(t, "/tmp/foo.jpg", path)
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
	job, _ := NewJob(user.ID)
	img, _ := NewImage(job, "/tmp/foo.jpg", "image/jpeg")
	dbw.SaveNewUser(user)
	dbw.SaveNewJob(job)
	assert.Nil(t, dbw.SaveNewImage(img))
	imgLoaded, err := dbw.LoadImage(img.ID)
	assert.Nil(t, err)
	assert.Equal(t, img.ID, imgLoaded.ID)
}
