package main

import (
	"database/sql"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

type DBWorker struct {
	Path string
	DB   *sql.DB
	mu   sync.Mutex
}

type SQL struct {
	Q    string
	Args []interface{}
}

func sqlStmt(query string, args ...interface{}) *SQL {
	return &SQL{Q: query, Args: args}
}

func NewDBWorker(path string) (*DBWorker, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}
	return &DBWorker{Path: path, DB: db}, nil
}

func (dbw *DBWorker) QueryRow(query string, args ...interface{}) *sql.Row {
	dbw.mu.Lock()
	defer dbw.mu.Unlock()
	return dbw.DB.QueryRow(query, args...)
}

func (dbw *DBWorker) WriteOne(query string, args ...interface{}) error {
	dbw.mu.Lock()
	defer dbw.mu.Unlock()
	stmt, err := dbw.DB.Prepare(query)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(args...)
	return err

}

func (dbw *DBWorker) Write(sqls ...*SQL) error {
	if len(sqls) == 1 {
		return dbw.WriteOne(sqls[0].Q, sqls[0].Args...)
	} else {
		dbw.mu.Lock()
		defer dbw.mu.Unlock()
		tx, err := dbw.DB.Begin()
		if err != nil {
			return err
		}
		defer tx.Rollback()
		for _, mutation := range sqls {
			stmt, err := tx.Prepare(mutation.Q)
			if err != nil {
				return err
			}
			_, err = stmt.Exec(mutation.Args...)
			stmt.Close()
			if err != nil {
				return err
			}
		}
		tx.Commit()
	}
	return nil
}

func (dbw *DBWorker) Close() {
	dbw.DB.Close()
}
