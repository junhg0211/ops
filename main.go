package main

import (
	"fmt"
	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
	"net/http"
)

func getRenderer() multitemplate.Renderer {
	r := multitemplate.NewRenderer()
	r.AddFromFiles("index", "templates/base.html", "templates/index.html")
	r.AddFromFiles("profile", "templates/base.html", "templates/profile.html")
	r.AddFromFiles("upload", "templates/base.html", "templates/upload.html")
	return r
}

func index(c *gin.Context) {
	c.HTML(http.StatusOK, "index", nil)
}

func profile(c *gin.Context) {
	email := c.Param("email")
	c.HTML(http.StatusOK, "profile", gin.H{
		"username": "Sch",
		"email":    email,
	})
}

func upload(c *gin.Context) {
	c.HTML(http.StatusOK, "upload", nil)
}

func main() {
	router := gin.Default()
	router.Static("/assets", "assets")
	router.HTMLRender = getRenderer()

	router.GET("/", index)
	router.GET("/index", index)
	router.GET("/profile/:email", profile)
	router.GET("/upload", upload)

	err := router.Run(":80")
	if err != nil {
		fmt.Println("running router:", err)
		return
	}
}
