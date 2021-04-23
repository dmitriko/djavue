package main

import (
	"github.com/gin-gonic/gin"
)

type App struct {
	DBW        *DBWorker
	MEDIA_ROOT string
}

func (app *App) postApiToken(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	if username == "" || password == "" {
		c.JSON(400, gin.H{
			"ok":    false,
			"error": "Missing username or password.",
		})
		return
	}
	user, err := app.DBW.LoadUserByName(username)
	if err != nil {
		var msg string
		if _, ok := err.(NotFoundError); ok {
			msg = "Wrong username or password."
		} else {
			msg = "Could not fetch data."
		}
		c.JSON(400, gin.H{
			"ok":    false,
			"error": msg,
		})
		return
	}
	if !user.IsPasswdValid(password) {
		c.JSON(400, gin.H{
			"ok":    false,
			"error": "Wrong username or password.",
		})
		return
	}
	c.JSON(200, gin.H{
		"ok":    true,
		"token": user.Token,
	})

}

func (app *App) postApiJob(c *gin.Context) {
	c.JSON(200, gin.H{
		"ok": true})
}

func setupRouter(dbw *DBWorker, mediaRoot string) *gin.Engine {
	r := gin.Default()
	app := &App{dbw, mediaRoot}
	r.POST("/api/token/", app.postApiToken)
	r.POST("/api/job/", app.postApiJob)
	return r
}
