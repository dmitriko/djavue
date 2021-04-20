package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func createTestUser(name, password string, dbw *DBWorker) error {
	user, _ := NewUser(name, password)
	return dbw.SaveNewUser(user)
}

type Resp struct {
	OK    bool   `json:"ok"`
	Error string `json:"error,omitempty"`
	Token string `json:"token,omitempty"`
}

func TestApiGetToken(t *testing.T) {
	dbw, err := testDBWorker()
	defer removeWorker(dbw)
	assert.Nil(t, err)
	assert.Nil(t, dbw.createTables())
	router := setupRouter(dbw)
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
