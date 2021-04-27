package main

import (
	"log"
	"os"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"golavue.dmitriko.com/api"
)

func main() {
	staticRoot := os.Getenv("DJAVUE_STATIC")
	if staticRoot == "" {
		log.Fatal("DJAVUE_STATIC is not set")
	}
	mediaRoot := os.Getenv("DJAVUE_MEDIA")
	if mediaRoot == "" {
		log.Fatal("DJAVUE_MEIDA is not set")
	}
	dbPath := os.Getenv("DJAVUE_DB_PATH")
	if dbPath == "" {
		log.Fatal("DJAVUE_DB_PATH is not set")
	}
	dbw, err := api.NewDBWorker(dbPath)
	if err != nil {
		log.Fatal(err)
	}

	r := gin.Default()
	router := api.SetupRouter(r, dbw, mediaRoot)
	router.Use(static.Serve("/", static.LocalFile(staticRoot, false)))
	router.Run(":8080")
}
