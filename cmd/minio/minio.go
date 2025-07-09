package minio

import (
	"context"
	"io"
	"log"
	"os"
	"time"

	"github.com/cheggaaa/pb"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

const (
	// endpoints  = "localhost:9000"
	endpoints  = "185.81.97.192:9000"
	accessKey  = "minioadmin"
	secretKey  = "minioadmin"
	useSSL     = false
	bucketName = "tutorial"
	location   = "us-east-1"
	debug      = false
)

// func main() {

// 	// createBucket(minioClient, ctx, bucketName, location, false)
// 	// objectName := "OBS-Studio-31.0.1-macOS-Apple.dmg"
// 	// filePath := "/Users/mehranmoradi/Downloads/OBS-Studio-31.0.1-macOS-Apple.dmg"
// 	// FPutObject(minioClient, ctx, bucketName, objectName, filePath, false)
// 	// PutObject(minioClient, ctx, bucketName, filePath, objectName, false)
// }

func makeMinioClient() (minioClient *minio.Client) {
	minioClient, err := minio.New(endpoints, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})

	if err != nil {
		log.Println("Cannot make client with error: ", err)
	}
	return minioClient
}

func createBucket(bucketName string, location string, debug bool) {

	minioClient := makeMinioClient()
	ctx := context.Background()

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

func FPutObject(bucketName string, objectName string, filePath string, debug bool) {

	minioClient := makeMinioClient()
	ctx := context.Background()

	const contentType = "application/octet-stream"

	if debug {
		minioClient.TraceOn(os.Stdout)
	}

	// Check object already exists:
	exist := objectExist(bucketName, objectName, false)
	if exist {
		log.Printf("Already object %s exist.\n", objectName)

	} else {
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
}

func PutObject(bucketName string, filePath string, objectName string, debug bool) {

	minioClient := makeMinioClient()
	ctx := context.Background()

	const contentType = "application/octet-stream"

	if debug {
		minioClient.TraceOn(os.Stdout)
	}

	// Check object already exists:
	exist := objectExist(bucketName, objectName, false)

	if exist {
		log.Printf("Already object %s exist.\n", objectName)

	} else {
		// Make open object:
		object, err := os.Open(filePath)
		if err != nil {
			log.Println("Cannot opef file with error: ", err)
		}
		defer object.Close()

		// Make stats object
		fileInfo, err := os.Stat(filePath)
		if err != nil {
			log.Println("Cannot stats file with: ", err)
		}

		// Make progress bar:
		progress := pb.New64(fileInfo.Size())
		progress.SetRefreshRate(50 * time.Microsecond)
		progress.Start()

		// Make upload a file:
		fileUploadInfo, err := minioClient.PutObject(ctx, bucketName, filePath, object, fileInfo.Size(), minio.PutObjectOptions{
			Progress:    progress,
			ContentType: contentType,
		})
		if err != nil {
			log.Println("Cannot upload a file with error: ", err)
		}

		log.Printf("File %s uploaded succesfully: ", fileUploadInfo.Key)
	}

}

func PutObjectApi(reader io.Reader, objectName string, fileSize int, progress *pb.ProgressBar) (err error) {

	minioClient := makeMinioClient()
	ctx := context.Background()

	const contentType = "application/octet-stream"

	if debug {
		minioClient.TraceOn(os.Stdout)
	}

	// Check object already exists:
	exist := objectExist(bucketName, objectName, false)

	if exist {
		log.Printf("Already object %s exist.\n", objectName)

	} else {

		// Make upload a file:
		fileUploadInfo, err := minioClient.PutObject(ctx, bucketName, objectName, reader, int64(fileSize), minio.PutObjectOptions{
			ContentType: contentType,
			Progress:    progress,
		})

		if err != nil {
			log.Println("Cannot upload a file with error: ", err)
			log.Printf("File %s uploaded succesfully: ", fileUploadInfo.Key)
			return err
		}
	}
	return
}

func objectExist(bucketName string, objectName string, debug bool) (exist bool) {

	minioClient := makeMinioClient()
	ctx := context.Background()
	const contentType = "application/octet-stream"
	var objectFind string

	// Make debug for trace:
	if debug {
		minioClient.TraceOn(os.Stdout)
	}

	// List objects options:
	options := minio.ListObjectsOptions{
		Prefix:    objectName,
		UseV1:     true,
		Recursive: true,
	}

	//List all obejcts:
	for object := range minioClient.ListObjects(ctx, bucketName, options) {
		objectFind = object.Key
	}
	if objectFind == objectName {
		return true
	}

	return
}
