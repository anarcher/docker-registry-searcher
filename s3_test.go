package main

import (
	"github.com/mitchellh/goamz/aws"
	"github.com/mitchellh/goamz/s3"
	"os"
	"testing"
)

func _TestS3Client(t *testing.T) {
	auth, err := aws.EnvAuth()
	if err != nil {
		t.Error(err)
	}

	region, ok := aws.Regions[os.Getenv("AWS_REGION")]
	if !ok {
		t.Error("AWS_REGION is wrong ", region)
	}

	bucketName := os.Getenv("S3_BUCKET")

	client := s3.New(auth, region)
	if client == nil {
		t.Error("client is wrong")
	}

	bucket := client.Bucket(bucketName)

	path := "registry/repositories/library/"
	var list *s3.ListResp
	list, err = bucket.List(path, "_json", "", 10)
	if err != nil {
		t.Error(err)
	}

	t.Log("Delimiter:", list.Delimiter)
	t.Log("Prefix:", list.Prefix)
	t.Log("Marker:", list.Marker)
	t.Log("IsTruncated:", list.IsTruncated)
	t.Log("MaxKeys:", list.MaxKeys)
	t.Log("NextMarker:", list.NextMarker)
	t.Log("MaxKeys:", list.MaxKeys)

	t.Log("CommonPrefixes:", len(list.CommonPrefixes))
	for _, prefix := range list.CommonPrefixes {
		t.Log(prefix)
	}

	t.Log("Contents:", len(list.Contents))
	for _, objects := range list.Contents {
		t.Log(objects.Key)
	}
}
