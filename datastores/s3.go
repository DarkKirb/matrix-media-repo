package datastores

import (
	"strconv"
	"sync"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/turt2live/matrix-media-repo/common/config"
)

var s3clients = &sync.Map{}

type s3 struct {
	client       *minio.Client
	storageClass string
	bucket       string
	putWithMd5   bool
}

func getS3(ds config.DatastoreConfig) (*s3, error) {
	if val, ok := s3clients.Load(ds.Id); ok {
		return val.(*s3), nil
	}

	endpoint := ds.Options["endpoint"]
	bucket := ds.Options["bucketName"]
	accessKeyId := ds.Options["accessKeyId"]
	accessSecret := ds.Options["accessSecret"]
	region := ds.Options["region"]
	storageClass, hasStorageClass := ds.Options["storageClass"]
	useSslStr, hasSsl := ds.Options["ssl"]
	useMd5Str, hasMd5 := ds.Options["useMD5"]

	if !hasStorageClass {
		storageClass = "STANDARD"
	}

	useSsl := true
	if hasSsl && useSslStr != "" {
		useSsl, _ = strconv.ParseBool(useSslStr)
	}

	useMd5 := false
	if hasMd5 && useMd5Str != "" {
		useMd5, _ = strconv.ParseBool(useMd5Str)
	}

	var err error
	var client *minio.Client
	client, err = minio.New(endpoint, &minio.Options{
		Region: region,
		Secure: useSsl,
		Creds:  credentials.NewStaticV4(accessKeyId, accessSecret, ""),
	})
	if err != nil {
		return nil, err
	}

	s3c := &s3{
		client:       client,
		storageClass: storageClass,
		bucket:       bucket,
		putWithMd5:   useMd5,
	}
	s3clients.Store(ds.Id, s3c)
	return s3c, nil
}
