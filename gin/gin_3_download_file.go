package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func FileDownload(c *gin.Context) {
	fmt.Println(">>>>>>>>>>>>>>>>>>>>>>")
	filename := c.PostForm("filename")
	fmt.Println(filename)
	c.FileAttachment("./README.md", "README.md")
	c.Status(http.StatusOK)
}

func main() {
	router := gin.Default()
	gin.DisableConsoleColor()
	router.Use(CORSMiddleware())
	router.POST("/download", FileDownload)
	router.Run(":8080")
}
