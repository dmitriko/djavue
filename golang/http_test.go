package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createTestUser(name, password string, dbw *DBWorker) (*User, error) {
	user, _ := NewUser(name, password)
	return user, dbw.SaveNewUser(user)
}

type Resp struct {
	OK    bool   `json:"ok"`
	Error string `json:"error,omitempty"`
	Token string `json:"token,omitempty"`
	JobID string `json:"job_id,omitempty"`
}

func TestApiGetToken(t *testing.T) {
	dbw, err := testDBWorker()
	defer removeWorker(dbw)
	assert.Nil(t, err)
	assert.Nil(t, dbw.createTables())
	router := setupRouter(dbw, "/tmp/foo")
	createTestUser("foo", "bar", dbw)
	w := httptest.NewRecorder()
	// Empty request
	resp := Resp{}
	data := url.Values{}
	req, _ := http.NewRequest("POST", "/api/token/", strings.NewReader(data.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	router.ServeHTTP(w, req)
	assert.Equal(t, 400, w.Code)
	decoder := json.NewDecoder(w.Body)
	err = decoder.Decode(&resp)
	assert.Nil(t, err)
	assert.False(t, resp.OK)
	assert.Equal(t, "Missing username or password.", resp.Error)
	// Only username
	resp = Resp{}
	data.Set("username", "foo")
	req, _ = http.NewRequest("POST", "/api/token/", strings.NewReader(data.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, 400, w.Code)
	decoder = json.NewDecoder(w.Body)
	err = decoder.Decode(&resp)
	assert.Nil(t, err)
	assert.False(t, resp.OK)
	assert.Equal(t, "Missing username or password.", resp.Error)
	// wrong password
	resp = Resp{}
	data = url.Values{}
	data.Set("username", "foo")
	data.Set("password", "dummy")
	req, _ = http.NewRequest("POST", "/api/token/", strings.NewReader(data.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	router.ServeHTTP(w, req)
	assert.Equal(t, 400, w.Code)
	decoder = json.NewDecoder(w.Body)
	err = decoder.Decode(&resp)
	assert.Nil(t, err)
	assert.False(t, resp.OK)
	assert.Equal(t, "Wrong username or password.", resp.Error)
	// wrong username
	resp = Resp{}
	data = url.Values{}
	data.Set("username", "spam")
	data.Set("password", "egg")
	req, _ = http.NewRequest("POST", "/api/token/", strings.NewReader(data.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, 400, w.Code)
	decoder = json.NewDecoder(w.Body)
	err = decoder.Decode(&resp)
	assert.Nil(t, err)
	assert.False(t, resp.OK)
	assert.Equal(t, "Wrong username or password.", resp.Error)
	// correct credentials
	resp = Resp{}
	data = url.Values{}
	data.Set("username", "foo")
	data.Set("password", "bar")
	req, _ = http.NewRequest("POST", "/api/token/", strings.NewReader(data.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	decoder = json.NewDecoder(w.Body)
	err = decoder.Decode(&resp)
	assert.Nil(t, err)
	assert.True(t, resp.OK)
	assert.True(t, len(resp.Token) > 10)
}

func TestApiJobPost(t *testing.T) {
	dbw, err := testDBWorker()
	defer removeWorker(dbw)
	assert.Nil(t, err)
	assert.Nil(t, dbw.createTables())
	user, _ := createTestUser("foo", "bar", dbw)
	router := setupRouter(dbw, "/tmp/foo")

	buf, contentType, err := createJobForm("test_data/img.png", "kind", "original")
	require.Nil(t, err)
	req, _ := http.NewRequest("POST", "/api/job/", buf)
	req.Header.Add("Content-Type", contentType)
	req.Header.Add("Authorization", "Token "+user.Token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	decoder := json.NewDecoder(w.Body)
	var resp Resp
	err = decoder.Decode(&resp)
	require.Nil(t, err)
	require.True(t, resp.OK)
	assert.True(t, resp.JobID != "")
	var job Job
	assert.Nil(t, dbw.LoadJob(&job, resp.JobID))

}

// Returns buffer with form, content type and error
func createJobForm(filePath string, fieldName, fieldValue string) (*bytes.Buffer, string, error) {
	var buf bytes.Buffer
	var err error
	w := multipart.NewWriter(&buf)
	file, err := os.Open(filePath)
	if err != nil {
		return &buf, "", err
	}
	var fw io.Writer
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="file"; filename="%s"`, file.Name()))
	h.Set("Content-Type", "image/png")
	if fw, err = w.CreatePart(h); err != nil {
		return &buf, "", err
	}
	if _, err = io.Copy(fw, file); err != nil {
		return &buf, "", err
	}
	if fw, err = w.CreateFormField(fieldName); err != nil {
		return &buf, "", err
	}
	if _, err = fw.Write([]byte(fieldValue)); err != nil {
		return &buf, "", err
	}
	w.Close()

	return &buf, w.FormDataContentType(), err
}
