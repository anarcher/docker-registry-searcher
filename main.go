package main

import (
	"flag"
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
	search_result_max     = config.Int("search_result_max", 1000)
	search_result_limit   = config.Int("search_result_limit", 2000)
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
		//n := int(c.Params.ByName("n"))
		//page := int(c.Params.ByName("page"))

		repos, err := LoadS3Repositories(bucket, *search_result_max, *search_result_limit)
		if err != nil {
			c.Fail(500, err)
		}

		var repoNames []string
		repoNames, err = repos.Search(q)
		if err != nil {
			c.Fail(500, err)
		}

		repoInfos := repos.InfosByNames(repoNames)
		total := len(repoInfos)

		//I didn't use paging.
		c.JSON(200, gin.H{
			"num_pages":   1,
			"num_results": total,
			"results":     repoInfos,
			"page_size":   total,
			"query":       q,
			"page":        1,
		})
		return
	})

	r.Run(*addr)
}
