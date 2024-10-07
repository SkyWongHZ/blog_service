package main

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/go-programming-tour-book/blog-service/pkg/minio"
)

func TestMinIOConnection(t *testing.T) {
	// 初始化配置
	if err := setupSetting(); err != nil {
		t.Fatalf("Failed to setup setting: %v", err)
	}

	// 初始化 MinIO 客户端
	if err := setupMinio(); err != nil {
		t.Fatalf("Failed to setup MinIO client: %v", err)
	}

	// 测试连接和权限
	err := testMinIOConnection()
	if err != nil {
		t.Fatalf("MinIO connection test failed: %v", err)
	}
}

func testMinIOConnection() error {
	// 列出所有 buckets
	buckets, err := minio.MinioClient.ListBuckets(context.Background())
	if err != nil {
		return fmt.Errorf("failed to list buckets: %v", err)
	}
	fmt.Printf("Found %d buckets\n", len(buckets))

	// 尝试创建一个测试 bucket
	bucketName := "test-bucket"
	err = minio.MinioClient.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{})
	if err != nil {
		// 检查 bucket 是否已存在
		exists, errBucketExists := minio.MinioClient.BucketExists(context.Background(), bucketName)
		if errBucketExists == nil && exists {
			fmt.Printf("Bucket %s already exists\n", bucketName)
		} else {
			return fmt.Errorf("failed to create bucket: %v", err)
		}
	} else {
		fmt.Printf("Successfully created bucket %s\n", bucketName)
	}

	// 尝试上传一个小文件
	objectName := "test-object"
	contentType := "application/octet-stream"
	_, err = minio.MinioClient.PutObject(context.Background(), bucketName, objectName, strings.NewReader("test content"), int64(len("test content")), minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		return fmt.Errorf("failed to upload object: %v", err)
	}
	fmt.Printf("Successfully uploaded test object to %s\n", bucketName)

	return nil
}
