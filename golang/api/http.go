package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

type App struct {
	DBW        *DBWorker
	MEDIA_ROOT string
}

type LoginCall struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (app *App) postApiToken(c *gin.Context) {
	var l LoginCall
	c.BindJSON(&l)
	if l.Username == "" || l.Password == "" {
		c.JSON(400, gin.H{
			"ok":    false,
			"error": "Missing username or password.",
		})
		return
	}
	var user User
	err := app.DBW.LoadUserByName(&user, l.Username)
	if err != nil {
		var msg string
		if IsNotFound(err) {
			msg = "Wrong username or password."
		} else {
			msg = fmt.Sprintf("Could not fetch data: %s", err.Error())
		}
		c.JSON(400, gin.H{
			"ok":    false,
			"error": msg,
		})
		return
	}
	if !user.IsPasswdValid(l.Password) {
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

type ImgResp struct {
	PK string `json:"pk"`
}

type JobResp struct {
	OK     bool      `json:"ok"`
	PK     string    `json:"pk"`
	Images []ImgResp `json:"images"`
}

func (app *App) getApiImage(c *gin.Context) {
	user, ok := c.MustGet("user").(*User)
	if !ok {
		respondErr(c, 401, "Not authorized")
		return
	}
	id := c.Param("id")
	var img Image
	if err := app.DBW.LoadImage(&img, id); err != nil {
		if IsNotFound(err) {
			respondErr(c, 404, "Could not find")
			return
		} else {
			respondErr(c, 500, "Could not fetch")
			return
		}
	}
	if user.ID != img.UserID {
		respondErr(c, 403, "Not allowed")
		return
	}
	c.Header("Content-Type", img.MimeType)
	c.File(img.Path)
}

func (app *App) getApiJob(c *gin.Context) {
	user, ok := c.MustGet("user").(*User)
	if !ok {
		respondErr(c, 401, "Not authorized")
		return
	}
	id := c.Param("id")
	var job Job
	if err := app.DBW.LoadJob(&job, id); err != nil {
		if IsNotFound(err) {
			respondErr(c, 404, "Could not find")
			return
		} else {
			respondErr(c, 500, "Could not fetch job")
			return
		}
	}
	if user.ID != job.UserID {
		respondErr(c, 403, "Not allowed")
		return
	}
	resp := JobResp{
		OK: true,
		PK: job.ID,
	}
	var imgs []Image
	err := app.DBW.Select(&imgs, "select * from images where job_id=?", job.ID)
	if err != nil {
		respondErr(c, 500, "Could not fetch images")
		return
	}
	for _, img := range imgs {
		resp.Images = append(resp.Images, ImgResp{PK: img.ID})
	}

	c.JSON(200, &resp)
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
	case JOB_SQUARE_ORIG:
		if err := performJobSquareOrig(app.DBW, job, app.MEDIA_ROOT, file); err != nil {
			respondErr(c, 400, err.Error())
			return
		}
	case JOB_SQUARE_SMALL:
		if err := performJobSquareSmall(app.DBW, job, app.MEDIA_ROOT, file); err != nil {
			respondErr(c, 400, err.Error())
			return
		}
	case JOB_ALL_THREE:
		if err := performJobAllThree(app.DBW, job, app.MEDIA_ROOT, file); err != nil {
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
	r := gin.Default()
	app := &App{dbw, mediaRoot}
	r.POST("/api/token/", app.postApiToken)
	r.POST("/api/job/", AuthMiddleware(dbw), app.postApiJob)
	r.GET("/api/job/:id/", AuthMiddleware(dbw), app.getApiJob)
	r.GET("/api/image/:id/", AuthMiddleware(dbw), app.getApiImage)
	return r
}
