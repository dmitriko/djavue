package main

import (
	"log"
	"os"
	"strings"

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
	var user User
	err := app.DBW.LoadUserByName(&user, username)
	if err != nil {
		var msg string
		if IsNotFound(err) {
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

func respondErr(c *gin.Context, code int, msg string) {
	c.AbortWithStatusJSON(code, gin.H{"ok": false, "error": msg})
}

func AuthMiddleware(dbw *DBWorker) gin.HandlerFunc {
	return func(c *gin.Context) {
		h, ok := c.Request.Header["Authorization"]
		if !ok {
			respondErr(c, 401, "Authoriaztion headers not provided.")
			return
		}
		token := strings.Replace(h[0], "Token ", "", -1)
		if token == "" {
			respondErr(c, 401, "Authorization token no valid.")
			return
		}
		var user User
		err := dbw.LoadUserByToken(&user, token)
		if err != nil {
			respondErr(c, 401, "Authorization token no user.")
			return
		}
		c.Set("user", &user)
		c.Next()
	}
}

func (app *App) postApiJob(c *gin.Context) {
	user, ok := c.MustGet("user").(*User)
	if !ok {
		respondErr(c, 401, "Not authorized")
		return
	}
	file, err := c.FormFile("file")
	if err != nil {
		respondErr(c, 400, "No file is received.")
		return
	}
	//	fmt.Printf("headers are %#v", file.Header)
	kind := c.PostForm("kind")
	job, err := NewJob(user.ID, kind)
	if err != nil {
		respondErr(c, 400, err.Error())
		return
	}
	if err := app.DBW.SaveNewJob(job); err != nil {
		println(err.Error())
		respondErr(c, 500, "Could not save job.")
		return
	}

	switch kind {
	case JOB_ORIG:
		if err := performJobOrig(app.DBW, job, app.MEDIA_ROOT, file); err != nil {
			respondErr(c, 400, err.Error())
			return
		}
	default:
		respondErr(c, 400, "Not valid job kind.")
		return
	}

	c.JSON(200, gin.H{
		"ok":     true,
		"job_id": job.ID,
	})
}

func setupRouter(dbw *DBWorker, mediaRoot string) *gin.Engine {
	if _, err := os.Stat(mediaRoot); os.IsNotExist(err) {
		log.Fatalf("%s is not accessible", mediaRoot)
	}
	r := gin.New()
	app := &App{dbw, mediaRoot}
	r.POST("/api/token/", app.postApiToken)
	r.POST("/api/job/", AuthMiddleware(dbw), app.postApiJob)
	return r
}
