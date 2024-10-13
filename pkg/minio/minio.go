package minio

import (
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var MinioClient *minio.Client

type PutObjectOptions = minio.PutObjectOptions
type MakeBucketOptions = minio.MakeBucketOptions
type RemoveObjectOptions = minio.RemoveObjectOptions

func NewMinioClient(endpoint, accessKeyID, secretAccessKey string, useSSL bool) error {
	var err error
	MinioClient, err = minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return err
	}
	return nil
}
