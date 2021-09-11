package main

import (
	"fmt"
	. "gee"
	"log"
	"net/http"
	"time"
)

func Logger() HandlerFunc {
	return func(c *Context) {
		// Start timer
		t := time.Now()
		// Process request
		c.Next()
		// Calculate resolution time
		log.Printf("[%d] %s in %v", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}

func onlyForV2() HandlerFunc {
	return func(c *Context) {
		// Start timer
		t := time.Now()
		// if a server error occurred
		//c.Fail(500, "Internal Server Error")
		// Calculate resolution time
		log.Printf("[%d] %s in %v for group v2", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}

func main() {
	// http.HandleFunc("/", indexHandler)
	// http.HandleFunc("/hello", helloHandler)
	r := New()
	r.GET("/index", func(c *Context) {
		c.HTML(http.StatusOK, "<h1>Index Page</h1>")
	})
	r.Use(Logger())
	v1 := r.Group("/v1")
	{
		v1.GET("/", indexHandler)
		v1.GET("/feng", func(c *Context) {
			// expect /hello?name=geektutu
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
		})
	}

	v2 := r.Group("/v2")
	v2.Use(onlyForV2())
	{
		v2.GET("/hello/:name", func(c *Context) {
			// expect /hello/geektutu
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
		})

		v2.GET("/assets/*filepath", func(c *Context) {
			c.JSON(http.StatusOK, H{"filepath": c.Param("filepath")})
		})
		v2.POST("/login", func(c *Context) {
			c.JSON(http.StatusOK, H{
				"username": c.PostFrom("username"),
				"password": c.PostFrom("password"),
			})
		})
	}

	log.Fatal(r.Run(":8888"))
}

func indexHandler(c *Context) {
	fmt.Println(c.Writer, "URL.Path = %q\n", c.Req.URL.Path)
}

func helloHandler(c *Context) {
	for k, v := range c.Req.Header {
		fmt.Fprintf(c.Writer, "Header[%q] = %q\n", k, v)
	}
}
