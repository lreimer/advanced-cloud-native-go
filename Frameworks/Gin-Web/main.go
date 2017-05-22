package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	engine := gin.Default()

	engine.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	// configuration for static files and templates
	engine.LoadHTMLGlob("./templates/*.html")
	engine.StaticFile("/favicon.ico", "./favicon.ico")

	engine.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "Advanced Cloud Native Go with Gin Framework",
		})
	})

	// the hello message endpoint with JSON response from map
	engine.GET("/hello", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Hello Gin Framework."})
	})

	// get all books
	engine.GET("/api/books", func(c *gin.Context) {
		c.JSON(http.StatusOK, AllBooks())
	})

	// create new book
	engine.POST("/api/books", func(c *gin.Context) {
		var book Book
		if c.BindJSON(&book) == nil {
			isbn, created := CreateBook(book)
			if created {
				c.Header("Location", "/api/books/"+isbn)
				c.Status(http.StatusCreated)
			} else {
				c.Status(http.StatusConflict)
			}
		}
	})

	// get book by ISBN
	engine.GET("/api/books/:isbn", func(c *gin.Context) {
		isbn := c.Params.ByName("isbn")
		book, found := GetBook(isbn)
		if found {
			c.JSON(http.StatusOK, book)
		} else {
			c.AbortWithStatus(http.StatusNotFound)
		}
	})

	// update existing book
	engine.PUT("/api/books/:isbn", func(c *gin.Context) {
		isbn := c.Params.ByName("isbn")

		var book Book
		if c.BindJSON(&book) == nil {
			exists := UpdateBook(isbn, book)
			if exists {
				c.Status(http.StatusOK)
			} else {
				c.Status(http.StatusNotFound)
			}
		}
	})

	// delete book
	engine.DELETE("/api/books/:isbn", func(c *gin.Context) {
		isbn := c.Params.ByName("isbn")
		DeleteBook(isbn)
		c.Status(http.StatusOK)
	})

	// run server on PORT
	engine.Run(port())
}

func port() string {
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8080"
	}
	return ":" + port
}
