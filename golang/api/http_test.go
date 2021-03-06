package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"os"
	"os/exec"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createTestUser(name, password string, dbw *DBWorker) (*User, error) {
	user, _ := NewUser(name, password)
	return user, dbw.SaveNewUser(user)
}

func setupTestRouter(dbw *DBWorker, mediaRoot string) *gin.Engine {
	r := gin.New()
	return SetupRouter(r, dbw, mediaRoot)
}

type Resp struct {
	OK    bool   `json:"ok"`
	Error string `json:"error,omitempty"`
	Token string `json:"token,omitempty"`
	JobID string `json:"job_id,omitempty"`
}

func NewJsonRequest(path string, data map[string]string) (*http.Request, error) {
	body, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", path, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	return req, nil
}

func TestApiGetToken(t *testing.T) {
	dbw, err := testDBWorker()
	defer removeWorker(dbw)
	assert.Nil(t, err)
	assert.Nil(t, dbw.createTables())
	router := setupTestRouter(dbw, "/tmp/foo")
	createTestUser("foo", "bar", dbw)
	w := httptest.NewRecorder()
	// Empty request
	var resp Resp
	req, _ := NewJsonRequest("/api/token/", map[string]string{})
	router.ServeHTTP(w, req)
	assert.Equal(t, 400, w.Code)
	decoder := json.NewDecoder(w.Body)
	err = decoder.Decode(&resp)
	assert.Nil(t, err)
	assert.False(t, resp.OK)
	assert.Equal(t, "Missing username or password.", resp.Error)
	// Only username
	resp = Resp{}
	req, _ = NewJsonRequest("/api/token/", map[string]string{"username": "foo"})
	req.Header.Add("Content-Type", "application/json")
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
	req, _ = NewJsonRequest("/api/token/", map[string]string{"username": "foo", "password": "dummy"})
	router.ServeHTTP(w, req)
	assert.Equal(t, 400, w.Code)
	decoder = json.NewDecoder(w.Body)
	err = decoder.Decode(&resp)
	assert.Nil(t, err)
	assert.False(t, resp.OK)
	assert.Equal(t, "Wrong username or password.", resp.Error)
	// wrong username
	resp = Resp{}
	req, _ = NewJsonRequest("/api/token/", map[string]string{"username": "spam", "password": "egg"})
	req.Header.Add("Content-Type", "application/json")
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
	req, _ = NewJsonRequest("/api/token/", map[string]string{"username": "foo", "password": "bar"})
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	decoder = json.NewDecoder(w.Body)
	err = decoder.Decode(&resp)
	assert.Nil(t, err)
	assert.True(t, resp.OK)
	assert.True(t, len(resp.Token) > 10)
}

func TestApiJobPostOrig(t *testing.T) {
	dbw, err := testDBWorker()
	defer removeWorker(dbw)
	assert.Nil(t, err)
	assert.Nil(t, dbw.createTables())
	user, _ := createTestUser("foo", "bar", dbw)
	os.MkdirAll("/tmp/foo", os.ModePerm)
	defer os.Remove("/tmp/foo")
	router := setupTestRouter(dbw, "/tmp/foo")

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
	var imgs []Image
	assert.Nil(t, dbw.Select(&imgs, "select * from images where job_id=?", job.ID))
	assert.Equal(t, 1, len(imgs))
	img := imgs[0]
	assert.Equal(t, int64(167314), img.Size)
	assert.Equal(t, 1236, img.Width)
	assert.Equal(t, 624, img.Height)
	stat, err := os.Stat(img.Path)
	assert.Nil(t, err)
	assert.Equal(t, img.Size, stat.Size())
}

func TestApiJobPostSquareOrig(t *testing.T) {
	dbw, err := testDBWorker()
	defer removeWorker(dbw)
	assert.Nil(t, err)
	assert.Nil(t, dbw.createTables())
	user, _ := createTestUser("foo", "bar", dbw)
	os.MkdirAll("/tmp/foo", os.ModePerm)
	defer os.Remove("/tmp/foo")
	router := setupTestRouter(dbw, "/tmp/foo")

	buf, contentType, err := createJobForm("test_data/img.png", "kind", "square_original")
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
	var imgs []Image
	assert.Nil(t, dbw.Select(&imgs, "select * from images where job_id=?", job.ID))
	require.Equal(t, 1, len(imgs))
	img := imgs[0]
	assert.NotEqual(t, int64(167314), img.Size)
	assert.Equal(t, 1236, img.Width)
	assert.Equal(t, 1236, img.Height)
	stat, err := os.Stat(img.Path)
	assert.Nil(t, err)
	assert.Equal(t, img.Size, stat.Size())
}

func TestApiJobPostSquareSmall(t *testing.T) {
	dbw, err := testDBWorker()
	defer removeWorker(dbw)
	assert.Nil(t, err)
	assert.Nil(t, dbw.createTables())
	user, _ := createTestUser("foo", "bar", dbw)
	os.MkdirAll("/tmp/foo", os.ModePerm)
	defer os.Remove("/tmp/foo")
	router := setupTestRouter(dbw, "/tmp/foo")

	buf, contentType, err := createJobForm("test_data/img.png", "kind", "square_small")
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
	var imgs []Image
	assert.Nil(t, dbw.Select(&imgs, "select * from images where job_id=?", job.ID))
	require.Equal(t, 1, len(imgs))
	img := imgs[0]
	assert.NotEqual(t, int64(167314), img.Size)
	assert.Equal(t, 256, img.Width)
	assert.Equal(t, 256, img.Height)
	stat, err := os.Stat(img.Path)
	assert.Nil(t, err)
	assert.Equal(t, img.Size, stat.Size())
}

func TestApiJobPostAllThree(t *testing.T) {
	dbw, err := testDBWorker()
	defer removeWorker(dbw)
	assert.Nil(t, err)
	assert.Nil(t, dbw.createTables())
	user, _ := createTestUser("foo", "bar", dbw)
	os.MkdirAll("/tmp/foo", os.ModePerm)
	defer os.Remove("/tmp/foo")
	router := setupTestRouter(dbw, "/tmp/foo")

	buf, contentType, err := createJobForm("test_data/img.png", "kind", "all_three")
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
	var imgs []Image
	assert.Nil(t, dbw.Select(&imgs, "select * from images where job_id=?", job.ID))
	require.Equal(t, 3, len(imgs))
}

func TestApiJobGet(t *testing.T) {
	dbw, err := testDBWorker()
	defer removeWorker(dbw)
	require.Nil(t, err)
	require.Nil(t, dbw.createTables())
	user, _ := createTestUser("foo", "bar", dbw)
	job, _ := NewJob(user.ID, JOB_ORIG)
	require.Nil(t, dbw.SaveNewJob(job))
	img, _ := NewImage(job, "/tmp", "foo.png", "image/png")
	require.Nil(t, dbw.SaveNewImage(img))
	req, _ := http.NewRequest("GET", fmt.Sprintf("/api/job/%s/", job.ID), nil)
	req.Header.Add("Authorization", "Token "+user.Token)

	w := httptest.NewRecorder()
	router := setupTestRouter(dbw, "/tmp/foo")
	router.ServeHTTP(w, req)
	require.Equal(t, 200, w.Code)
	decoder := json.NewDecoder(w.Body)
	var resp JobResp
	err = decoder.Decode(&resp)
	require.Nil(t, err)
	require.True(t, resp.OK)
	assert.Equal(t, job.ID, resp.PK)
	assert.Equal(t, 1, len(resp.Images))
}

func TestApiImageGet(t *testing.T) {
	dbw, err := testDBWorker()
	defer removeWorker(dbw)
	assert.Nil(t, err)
	assert.Nil(t, dbw.createTables())
	user, _ := createTestUser("foo", "bar", dbw)
	os.MkdirAll("/tmp/foo", os.ModePerm)
	defer os.Remove("/tmp/foo")
	job, _ := NewJob(user.ID, JOB_ORIG)
	img, _ := NewImage(job, "/tmp/foo", "foo.png", "image/png")
	require.Nil(t, dbw.SaveNewJob(job))
	require.Nil(t, dbw.SaveNewImage(img))
	cpCmd := exec.Command("cp", "test_data/img.png", img.Path)
	require.Nil(t, cpCmd.Run())
	req, _ := http.NewRequest("GET", fmt.Sprintf("/api/image/%s/", img.ID), nil)
	req.Header.Add("Authorization", "Token "+user.Token)

	w := httptest.NewRecorder()
	router := setupTestRouter(dbw, "/tmp/foo")
	router.ServeHTTP(w, req)
	require.Equal(t, 200, w.Code)
	contentType, ok := w.HeaderMap["Content-Type"]
	require.True(t, ok)
	assert.Equal(t, "image/png", contentType[0])

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
