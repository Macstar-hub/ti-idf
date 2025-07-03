package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/cheggaaa/pb"
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

	createBucket(minioClient, ctx, bucketName, location, false)
	objectName := "2025-07-02 09-34-47.mov"
	filePath := "/Users/mehranmoradi/Desktop/2025-07-02 09-34-47.mov"
	putObject(minioClient, ctx, bucketName, objectName, filePath, false)

}

func createBucket(minioClient *minio.Client, ctx context.Context, bucketName string, location string, debug bool) {
	if debug {
		minioClient.TraceOn(os.Stdout)
	}
	makeBucketError := minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: location, ForceCreate: true})
	if makeBucketError != nil {
		exists, errorBucketExist := minioClient.BucketExists(ctx, bucketName)
		if errorBucketExist == nil && exists == true {
			log.Printf("Already have bucket name: %s\n", bucketName)
		} else {
			log.Println("Cannot find bucket name: ", errorBucketExist)
		}
	}

	exists, errorBucketExist := minioClient.BucketExists(ctx, bucketName)
	if errorBucketExist == nil && exists == true {
		log.Printf("Already have bucket name: %s\n", bucketName)
	}
}

func putObject(minioClient *minio.Client, ctx context.Context, bucketName string, objectName string, filePath string, debug bool) {
	const contentType = "application/octet-stream"

	if debug {
		minioClient.TraceOn(os.Stdout)
	}

	fileInfo, err := os.Stat(filePath)
	if err != nil {
		log.Println("Cannot stats file with error: ", err)
	}

	// Make progress bar:
	progress := pb.New64(fileInfo.Size())
	progress.SetRefreshRate(50 * time.Microsecond)
	progress.Start()

	// Upload file:
	info, uploadError := minioClient.FPutObject(ctx, bucketName, objectName, filePath, minio.PutObjectOptions{
		ContentType: contentType,
		Progress:    progress,
	})

	if uploadError != nil {
		log.Println("Cannot upload a file with error: ", uploadError)
	}
	log.Printf("Upload successfuly with %s\n", info.Bucket)
}
