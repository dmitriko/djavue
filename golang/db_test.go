package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Creates DBWorker with pointed to temp file
func testDBWorker() (*DBWorker, error) {
	tmpFile, _ := ioutil.TempFile("", "db")
	return NewDBWorker(tmpFile.Name())
}

// Closes DBWorker and unlinks its file
func removeWorker(dbw *DBWorker) {
	dbw.Close()
	os.Remove(dbw.Path)
}

func TestNewDBWorker(t *testing.T) {
	dbw, err := testDBWorker()
	if assert.Nil(t, err) {
		defer removeWorker(dbw)
	}
}

func TestDBWriteRead(t *testing.T) {
	dbw, err := testDBWorker()
	if assert.Nil(t, err) {
		defer removeWorker(dbw)
	}
	err = dbw.Write(sqlStmt("create table foo(id int, name text)"))
	assert.Nil(t, err)
	assert.Nil(t, dbw.Write(sqlStmt("insert into foo values (?, ?)", 1, "foo"),
		sqlStmt("insert into foo values (?, ?)", 2, "bar")))
	var name string
	err = dbw.QueryRow("select name from foo where id = ?", 2).Scan(&name)
	assert.Nil(t, err)
	assert.Equal(t, name, "bar")
}

func TestMutithreadWrite(t *testing.T) {
	var wg sync.WaitGroup
	dbw, err := testDBWorker()
	assert.Nil(t, err)
	assert.Nil(t, dbw.Write(sqlStmt("create table foo(id int, name text)")))
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func(i int) {
			err = dbw.Write(sqlStmt("insert into foo values (?, ?)", i, fmt.Sprintf("name%d", i)))
			assert.Nil(t, err)
			wg.Done()
		}(i)
	}
	wg.Wait()
	var name string
	err = dbw.QueryRow("select name from foo where id = ?", 1).Scan(&name)
	assert.Nil(t, err)
	if assert.Equal(t, "name1", name) {
		removeWorker(dbw)
	} else {
		println(dbw.Path)
	}
}
