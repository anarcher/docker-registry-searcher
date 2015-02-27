package main

import (
	"flag"
	"github.com/drone/config"
	"github.com/gin-gonic/gin"
	"github.com/mitchellh/goamz/aws"
	"github.com/mitchellh/goamz/s3"
	"log"
	"os"
)

var (
	configFile = flag.String("c", "/etc/docker-registry-searcher.toml", "docker-registry-search.toml")
)

var (
	debug_mode            = config.Bool("debug", false)
	gin_mode              = config.String("gin-mode", gin.DebugMode)
	ip                    = config.String("ip", "")
	port                  = config.String("port", "8080")
	aws_access_key_id     = config.String("aws-access-key", "")
	aws_secret_access_key = config.String("aws-secret-key", "")
	aws_region            = config.String("aws-region", "")
	s3_bucket             = config.String("s3-bucket", "")
	search_result_max     = config.Int("search-result-max", 1000)
	search_result_limit   = config.Int("search-result-limit", 2000)
)

func main() {
	flag.Parse()
	config.SetPrefix("DS_")
	if _, err := os.Stat(*configFile); err != nil {
		log.Println(err)
		config.Parse("")
	} else {
		config.Parse(*configFile)
	}

	printConfigs()

	awsAuth, err := aws.GetAuth(*aws_access_key_id, *aws_secret_access_key)
	if err != nil {
		log.Fatal(err)
		return
	}

	region, ok := aws.Regions[*aws_region]
	if !ok {
		log.Fatalf("aws_region is wrong:%v", *aws_region)
		return
	}
	s3Client := s3.New(awsAuth, region)
	bucket := s3Client.Bucket(*s3_bucket)

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(gin.ErrorLogger())

	r.GET("/v1/search", func(c *gin.Context) {
		log.Printf("DEBUG_MODE: %v", *debug_mode)
		q := c.Request.URL.Query().Get("q")
		//n := c.Request.URL.Query().Get("n")
		//page := c.Request.URL.Query().Get("page")

		repos, err := LoadS3Repositories(bucket, *search_result_max, *search_result_limit)
		if err != nil {
			c.Fail(500, err)
			return
		}

		var repoNames []string
		repoNames, err = repos.Search(q)
		if err != nil {
			c.Fail(500, err)
			return
		}

		if *debug_mode == true {
			log.Printf("%v", repoNames)
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

	addr := *ip + ":" + *port
	r.Run(addr)
}

func printConfigs() {
	if *debug_mode == false {
		return
	}
	log.Printf("DEBUG: %v\n", *debug_mode)
	log.Printf("GIN_MODE: %v\n", *gin_mode)
	log.Printf("ADDR: %v:%v\n", *ip, *port)
	log.Printf("AWS_ACCESSKEY: %v\n", *aws_access_key_id)
	log.Printf("AWS_SECRETKEY: %v\n", *aws_secret_access_key)
	log.Printf("AWS_REGION: %v\n", *aws_region)
	log.Printf("S3_BUCKEY: %v\n", *s3_bucket)
	log.Printf("SEARCH_RESULT_MAX: %v\n", *search_result_max)
	log.Printf("SEARCH_RESULT_LIMIT: %v\n", *search_result_limit)

}
