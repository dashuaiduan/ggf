package main

import (
	"fmt"
	"gotest/gee"
	"log"
	"net/http"
	_ "net/http"
	"time"
)

func onlyForV2() gee.HandlerFunc {
	return func(c *gee.Context) {
		// Start timer
		t := time.Now()
		// if a server error occurred
		//c.String(500, "Internal Server Error")
		// Calculate resolution time
		log.Printf("[%d] %s in %v for group v2", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}

func main() {
	engine := gee.New()
	engine.Use(gee.Logger())

	engine.GET("/", indexHandler)
	engine.GET("/hello", helloHandler)
	r := engine.Group("/v1")
	r1 := r.Group("/v11")
	fmt.Println(r1)
	r.Use(onlyForV2())
	{
		r.GET("/index", v1_index)
		r.GET("/hello", v1_hello)
	}
	r.GET("/panic", func(c *gee.Context) {
		names := []string{"geektutu"}
		c.String(http.StatusOK, names[100])
	})

	engine.Static("/assets", "D:/gotest/assets")
	engine.Run(":8080")
}

func indexHandler(c *gee.Context) {
	//fmt.Fprintf(c.Writer, "666.Path = %q\n", c.Req.URL.Path)
	//c.String(200,"sdfdg")
	//a := make(map[string]string)
	//a["aa"] = "aaaaaaaaaaaaa"
	//a["bb"] = "bbbbbbbbbbb"
	//c.JSON(200,gee.H{"sdf":"sdf","fff":"fff"})
	//w.Header().Set()
	//w.Write([]byte("dfsfg"))
	//w.WriteHeader()
	//c.HTML(400,"<div>666</div><div>777</div>")
	//println(c.Query("user"))
	c.String(200, "index")
}

func helloHandler(c *gee.Context) {
	//fmt.Fprint(c.Writer, c.Req.Header)
	c.String(200, c.PostForm("name"))
	//for k, v := range request.Header {
	//	fmt.Fprintf(w, "Header[%s] = %s\n", k, v)
	//}
}
func v1_index(c *gee.Context) {
	arr := []int{1, 2, 3}
	println(arr[8])
	println(111)
	c.String(200, "v1_index")
}
func v1_hello(c *gee.Context) {
	c.String(200, "v1_hello")
}
