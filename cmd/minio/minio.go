package main

import (
	"context"
	"log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func main() {
	endpoints := "localhost:9000"
	accessKey := "minioadmin"
	secretKey := "minioadmin"
	useSSL := false
	ctx := context.Background()

	minioClient, err := minio.New(endpoints, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})

	if err != nil {
		log.Println("Cannit make client with error: ", err)
	}

	// Define bucket name:
	bucketName := "tutorial"
	location := "us-east-1"

	// Check bucket existing check:
	makeBucketError := minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: location, ForceCreate: true})
	if makeBucketError != nil {
		exists, errorBucketExist := minioClient.BucketExists(ctx, bucketName)
		if errorBucketExist == nil && exists {
			log.Printf("Already have bucket name: ", bucketName)
		} else {
			log.Println("Cannot find bucket name: ", errorBucketExist)
		}
	}
	objectName := "2025-07-02 09-34-47.mov"
	filePath := "/Users/mehranmoradi/Desktop/2025-07-02 09-34-47.mov"
	contentType := "application/octet-stream"

	// Upload file:
	info, uploadError := minioClient.FPutObject(ctx, bucketName, objectName, filePath, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if uploadError != nil {
		log.Println("Cannot upload a file with error: ", uploadError)
	}
	log.Printf("Upload successfuly with %s\n", info.Bucket)
}
