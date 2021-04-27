package main

import (
	"log"
	"os"

	"golavue.dmitriko.com/api"
)

func main() {
	dbPath := os.Getenv("DJAVUE_DB_PATH")
	if dbPath == "" {
		log.Fatal("DJAVUE_DB_PATH is not set")
	}
	dbw, err := api.NewDBWorker(dbPath)
	if err != nil {
		log.Fatal(err)
	}
	if err := dbw.CreateTables(); err != nil {
		//		log.Fatal(err)
	}
	user, _ := api.NewUser("foo", "f00baRRR")
	if err := dbw.SaveNewUser(user); err != nil {
		log.Fatal(err)
	}
}
