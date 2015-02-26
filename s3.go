package main

import (
	"github.com/mitchellh/goamz/s3"
	"strings"
)

const (
	PATH_REGISTRY_REPOSITORIES_LIBRARY = "/registry/repositories/library"
	DELIMITER                          = "/"
)

type S3Repositories struct {
	List   []string
	bucket *s3.Bucket
}

func LoadS3Repositories(bucket *s3.Bucket, max int, limit int) (*S3Repositories, error) {
	repo := &S3Repositories{
		bucket: bucket,
	}
	err := repo.Read(max, limit)

	if err != nil {
		return nil, err
	}

	return repo, nil
}

func (r *S3Repositories) Read(max, limit int) error {
	var marker string

	for {
		list, err := r.bucket.List(PATH_REGISTRY_REPOSITORIES_LIBRARY, DELIMITER, marker, max)
		if err != nil {
			return err
		}

		for _, dir := range list.CommonPrefixes {
			r.List = append(r.List, dir)
		}

		if list.IsTruncated == true {
			marker = list.NextMarker
			if limit > 0 && len(r.List) >= limit {
				break
			}
		} else {
			break
		}

	}

	return nil
}

func (r *S3Repositories) Search(q string) (result []string, err error) {
	for _, i := range r.List {
		if strings.Contains(i, q) {
			result = append(result, i)
		}
	}

	return
}

func (r S3Repositories) InfosByNames(names []string) (infos []map[string]string) {
	for _, name := range names {
		infos = append(infos, map[string]string{"name": name, "description": ""})
	}

	return
}
