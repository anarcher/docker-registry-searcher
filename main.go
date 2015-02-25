package main

import (
	"flag"
	"fmt"
	"github.com/drone/config"
	"github.com/gin-gonic/gin"
	"github.com/mitchellh/goamz/aws"
	"github.com/mitchellh/goamz/s3"
	"log"
)

var (
	configFile = flag.String("c", "/etc/docker-registry-searcher.toml", "docker-registry-search.toml")
)

var (
	gin_mode              = config.String("gin_mode", gin.DebugMode)
	addr                  = config.String("addr", ":8080")
	aws_access_key_id     = config.String("aws_access_key_id", "")
	aws_secret_access_key = config.String("aws_secret_access_key", "")
	aws_bucket            = config.String("aws_bucket", "")
	search_result_max     = config.Int("search_result_max", 100)
)

const (
	PATH_REGISTRY_REPOSITORIES_LIBRARY = "/registry/repositories/library"
	DELIMITER                          = "/"
)

func main() {

	config.Parse(*configFile)

	awsAuth, err := aws.GetAuth(*aws_access_key_id, *aws_secret_access_key)
	if err != nil {
		log.Fatal(err)
		return
	}
	s3Client := s3.New(awsAuth, aws.USEast)
	bucket := s3Client.Bucket(*aws_bucket)

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(gin.ErrorLogger())

	r.GET("/v1/search", func(c *gin.Context) {
		q := c.Params.ByName("q")
		list, err := bucket.List(PATH_REGISTRY_REPOSITORIES_LIBRARY, DELIMITER, q, *search_result_max)
		if err != nil {
			c.Fail(500, err)
			return
		}
		fmt.Println(list)
		c.String(200, "Hello world: "+q)
		return
	})

	r.Run(*addr)
}
