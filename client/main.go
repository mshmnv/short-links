package main

import (
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"log"
	"net/http"
	ps "proto"
)

func main() {
	conn, err := grpc.Dial("server_container:4040", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Unable to connect: %s", err)
	}
	client := ps.NewShortClient(conn)

	router := gin.Default()

	router.Static("/css", "html/css")
	router.LoadHTMLGlob("html/*.html")

	// main page
	router.GET("/", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "index.html", gin.H{
			"title": "short.",
		})
	})

	// get link and give url
	router.GET("/get", func(ctx *gin.Context) {
		req := &ps.LinkRequest{Link: ctx.Query("link")}
		if response, err := client.Get(ctx, req); err == nil {
			ctx.HTML(http.StatusOK, "get.html", gin.H{
				"title": "short.",
				"url":   response.Url,
			})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
	})
	// save url and link in bd and give a link to client
	router.POST("/", func(ctx *gin.Context) {
		req := &ps.UrlRequest{Url: ctx.PostForm("url")}
		if response, err := client.Create(ctx, req); err == nil {
			ctx.HTML(http.StatusOK, "create.html", gin.H{
				"title": "short.",
				"url":   req.Url,
				"link":  response.Link,
			})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
	})
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
