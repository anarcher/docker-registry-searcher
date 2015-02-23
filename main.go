package main

import (
	"github.com/drone/config"
	"github.com/gin-gonic/gin"
	"gopkg.in/amz.v2/aws"
	"gopkg.in/amz.v2/s3"
)

var (
	gin_mode              = config.String("gin_mode", gin.DebugMode)
	addr                  = config.String("addr", ":8080")
	aws_access_key_id     = config.String("aws_access_key_id", "")
	aws_secret_access_key = config.String("aws_secret_access_key", "")
)

func main() {

	awsAuth := aws.Auth{*aws_access_key_id, *aws_secret_access_key}

	s3Client := s3.New(awsAuth, aws.USEast)

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.GET("/v1/search", func(c *gin.Context) {
		q := c.Params.ByName("q")
		c.String(200, "Hello world: "+q)
	})

	r.Run(*addr)
}
