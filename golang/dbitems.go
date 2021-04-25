package main

import (
	"errors"
	"fmt"
	"mime/multipart"
	"path/filepath"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       string
	Username string
	Password string
	Token    string
}

const (
	JobStateStarted = iota
	JobStateFailed
	JobStateDone
)

const (
	JOB_ORIG         = "original"
	JOB_SQUARE_ORIG  = "square_original"
	JOB_SQUARE_SMALL = "square_small"
	JOB_ALL_THREE    = "all_three"
)

var JobKind = map[string]bool{
	JOB_ORIG:         true,
	JOB_SQUARE_ORIG:  true,
	JOB_SQUARE_SMALL: true,
	JOB_ALL_THREE:    true,
}

type Job struct {
	ID     string
	UserID string `db:"user_id"`
	State  int64
	Kind   string
}

type Image struct {
	ID       string
	UserID   string `db:"user_id"`
	JobID    string `db:"job_id"`
	MimeType string `db:"mime_type"`
	Path     string
	Size     int64
	Width    int64
	Height   int64
}

func (dbw *DBWorker) CreateUserTable() error {
	schema := `create table users (
			id text primary key, 
			username text unique not null, 
			password text not null, 
			token text unique not null)`
	return dbw.WriteOne(schema)
}

func (dbw *DBWorker) CreateJobTable() error {
	schema := `create table jobs (
			id text primary key,
			user_id text not null,
			state int not null,
			kind text not null,
			foreign key (user_id) 
				references users (id)
				)`
	return dbw.WriteOne(schema)
}

func (dbw *DBWorker) CreateImageTable() error {
	schema := `create table images (
			id text primary key,
			user_id text not null,
			job_id text not null,
			path text unique not null,
			mime_type text not null,
			size int not null,
			width int not null,
			height int not null,
			foreign key (user_id)
				references users (id),
			foreign key (job_id)
				references jobs (id)
			)`
	return dbw.WriteOne(schema)
}

func (dbw *DBWorker) createTables() error {
	if err := dbw.CreateUserTable(); err != nil {
		return err
	}
	if err := dbw.CreateJobTable(); err != nil {
		return err
	}
	if err := dbw.CreateImageTable(); err != nil {
		return err
	}
	return nil
}

func (dbw *DBWorker) SaveNewUser(user *User) error {
	return dbw.WriteOne("insert into users (id, username, password, token) values (?,?,?,?)",
		user.ID, user.Username, user.Password, user.Token)
}

func (dbw *DBWorker) SaveUser(user *User) error {
	_, err := dbw.NamedExec(`update users 
                            set username=:username, password=:password, token=:token
							where id=:id`, user)
	return err
}

func (dbw *DBWorker) LoadUser(id string) (*User, error) {
	user := &User{}
	err := dbw.Get(user, "select * from users where id=?", id)
	return user, err
}

func (dbw *DBWorker) LoadUserByName(username string) (*User, error) {
	user := &User{}
	err := dbw.Get(user, "select * from users where username=?", username)
	return user, err
}

func (dbw *DBWorker) LoadUserByToken(token string) (*User, error) {
	user := &User{}
	err := dbw.Get(user, "select * from users where token=?", token)
	return user, err
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
	if passwd == "" {
		return false // no empty passwords allowed
	}
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(passwd))
	if err != nil {
		return false
	}
	return true
}

func NewJob(userID, kind string) (*Job, error) {
	if !JobKind[kind] {
		return nil, errors.New("Wrong job kind.")
	}
	job := &Job{}
	id, err := NewULIDNow()
	if err != nil {
		return job, err
	}
	job.ID = id
	job.UserID = userID
	job.State = JobStateStarted
	job.Kind = kind
	return job, nil
}

func (dbw *DBWorker) SaveNewJob(job *Job) error {
	return dbw.WriteOne("insert into jobs (id, user_id, state, kind) values(?,?,?,?)", job.ID, job.UserID, job.State, job.Kind)
}

func (dbw *DBWorker) SaveJob(job *Job) error {
	return dbw.WriteOne("update jobs set state = ? where id = ?", job.State, job.ID)
}

func (dbw *DBWorker) LoadJob(jobID string) (*Job, error) {
	var j Job
	err := dbw.Get(&j, "select * from jobs where id=?", jobID)
	return &j, err
}

func NewImageFromFileHeader(job *Job, fileHeader *multipart.FileHeader, mediaRoot string) (*Image, error) {
	img := &Image{}
	return img, nil
}

func NewImage(job *Job, mediaRoot, fileName, mimeType string) (*Image, error) {
	img := &Image{}
	id, err := NewULIDNow()
	if err != nil {
		return img, err
	}
	img.ID = id
	img.JobID = job.ID
	img.UserID = job.UserID
	img.Path = filepath.Join(mediaRoot, fmt.Sprintf("%s_%s", id, fileName))
	return img, nil
}

func (dbw *DBWorker) SaveNewImage(img *Image) error {
	_, err := dbw.NamedExec(`insert into images (
		id, job_id, user_id, path, mime_type, size, width, height) values (
		:id, :job_id, :user_id, :path, :mime_type, :size, :width, :height)`, img)
	return err
}

func (dbw *DBWorker) LoadImage(id string) (*Image, error) {
	var img Image
	err := dbw.Get(&img, "select * from images where id = ?", id)
	return &img, err
}
