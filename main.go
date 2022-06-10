package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/thinkerou/favicon"

	"blog/database_handler"

	"github.com/gin-gonic/gin"
)

var bio []string
var password string

var posts []database_handler.Post

func main() {
	// Uncomment for release
	//gin.SetMode(gin.ReleaseMode)

	loadAbout() // load about page content from file

	// need to load past 5 blog posts here TODO
	top_five := database_handler.RetrieveTopFive()
	posts = append(posts, top_five...)

	password = os.Getenv("BLOG_PASS")
	if password == "" {
		panic("Couldn't load the password from environment. Set the blog post password")
	}

	// router setup
	router := gin.Default()
	//router.SetTrustedProxies(nil)
	router.LoadHTMLGlob("./templates/*.html")

	// load static folders
	router.Static("/css", "./templates/css")
	router.Static("./content", "./content")

	router.Use(favicon.New("favicon.ico")) // set favicon middleware

	// routing
	// GET
	router.GET("/about", about)
	router.GET("/", index)

	// POST
	router.POST("/api/post", submit)
	// end routing

	router.Run(":8080")
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
		IMAGE string `json:"imageLink"`
	}

	var content post_structure
	if err := c.BindJSON(&content); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "Could not bind data"})
		return
	}

	if content.PASS == password {
		if PostBlog(content.TITLE, content.BODY, content.IMAGE) {
			c.JSON(http.StatusOK, gin.H{"status": "post successfully uploaded"})
		} else {
			c.JSON(http.StatusRequestTimeout, "Can't connect to database")
		}
	} else { // error binding the query to type
		c.JSON(http.StatusForbidden, "Wrong password")
		return
	}
}

func PostBlog(title string, body string, pic string) bool {
	// TODO post to db or whatever storage choice
	// can prob just store as json locally
	new_post := database_handler.Post{Title: title, Timestamp: time.Now().Format("Monday Jan 2 2006 3:04PM"),
		Body: body, Image: pic}

	err := database_handler.Add(title, new_post.Timestamp, base64.RawURLEncoding.EncodeToString([]byte(body)), pic)
	if err != nil {
		fmt.Println(err)
		return false
	}

	// new post goes at head to display in order of most recent
	posts = append([]database_handler.Post{new_post}, posts...)

	fmt.Printf("Blog post received. Now showing %d posts\n", len(posts))

	return true
}

func loadAbout() {
	bio_bytes, err := ioutil.ReadFile("./content/bio.txt")
	bio = strings.Split(string(bio_bytes), "\\n")

	if err != nil {
		log.Output(1, err.Error())
	}
}
