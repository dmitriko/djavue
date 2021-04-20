package main

import (
	"strings"

	"github.com/gin-gonic/gin"
)

func setupRouter(dbw *DBWorker) *gin.Engine {
	r := gin.Default()
	r.POST("/api/token/", func(c *gin.Context) {
		username := c.PostForm("username")
		password := c.PostForm("password")
		if username == "" || password == "" {
			c.JSON(400, gin.H{
				"ok":    false,
				"error": "Missing username or password.",
			})
			return
		}
		user, err := dbw.LoadUserByName(username)
		if err != nil {
			var msg string
			if strings.Contains(err.Error(), "no rows") {
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

	})
	return r
}
