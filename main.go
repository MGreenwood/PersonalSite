package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/thinkerou/favicon"

	"github.com/gin-gonic/gin"
)

var bio []string
var password string

var posts []blog_post

type blog_post struct {
	Title     string
	Body      string
	Timestamp string
}

func main() {
	// Uncomment for release
	//gin.SetMode(gin.ReleaseMode)

	loadAbout() // load about page content from file

	// need to load past 5 blog posts here TODO

	password = os.Getenv("BLOG_PASS")
	if password == "" {
		panic("Couldn't load the password from environment. Set the blog post password")
	}

	// router setup
	router := gin.Default()
	router.LoadHTMLGlob("./templates/*.html")
	router.Static("/css", "./templates/css")
	router.Static("/content", "./content/*.pdf")
	router.Use(favicon.New("favicon.ico")) // set favicon middleware

	// routing
	// GET
	router.GET("/about", about)
	router.GET("/", index)

	// POST
	router.POST("/api/post", submit)
	// end routing

	router.Run("localhost:8080")
}

func about(c *gin.Context) {
	c.HTML(http.StatusOK, "about.html", gin.H{
		"title": "About",
		"bio":   ([]string)(bio),
	})
}

func index(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title":      "Home",
		"blog_posts": posts,
	})
}

func submit(c *gin.Context) {
	type post_structure struct {
		TITLE string `json:"title" binding:"required"`
		BODY  string `json:"body" binding:"required"`
		PASS  string `json:"password" binding:"required"`
	}

	var content post_structure
	if err := c.BindJSON(&content); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "Could not bind data"})
		return
	}
	if content.PASS == password {
		PostBlog(content.TITLE, content.BODY)
		c.JSON(http.StatusOK, gin.H{"status": "post successfully uploaded"})
	} else { // error binding the query to type
		c.JSON(http.StatusForbidden, "Wrong password")
		return
	}

}

func PostBlog(title string, body string) {
	// TODO post to db or whatever storage choice
	// can prob just store as json locally

	posts = append(posts, blog_post{title, body, time.Now().Format("Monday Jan 2 2006")})

	fmt.Println("Blog post received")
	fmt.Printf("Now showing %d posts", len(posts))
}

func loadAbout() {
	bio_bytes, err := ioutil.ReadFile("./content/bio.txt")
	bio = strings.Split(string(bio_bytes), "\\n")

	if err != nil {
		log.Output(1, err.Error())
	}
}
