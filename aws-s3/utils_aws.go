package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

//Funcion que inicializa el objeto que nos permitira conectar a AWS y acceder a los metodos del S3
func initS3Client(ctx context.Context, region string) (*s3.Client, error) {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		return nil, fmt.Errorf("Config error: %s", err)
	}
	return s3.NewFromConfig(cfg), nil
}

//Funcion que crea un Bucket S3 si no existe
func createS3Bucket(ctx context.Context, s3Client *s3.Client) error {

	_, err := s3Client.CreateBucket(ctx, &s3.CreateBucketInput{
		Bucket: aws.String(configuration.BucketName),
	})
	if err != nil {
		return fmt.Errorf("CreateBucket error: %s", err)
	}
	return nil
}

//Funcion que sube un archivo que nosotros indiquemos al S#
func uploadFileToS3(fileName string, ctx context.Context, s3Client *s3.Client) error {
	uploader := manager.NewUploader(s3Client)
	_, err := uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket: aws.String(configuration.BucketName),
		Key:    aws.String(fileName),
		Body:   strings.NewReader("this is a test"),
	})
	if err != nil {
		return fmt.Errorf("Upload error: %s", err)
	}
	return nil
}

//Funcion que baja el archivo que indiquemos del S3
func downloadFileFromS3(fileName string, ctx context.Context, s3Client *s3.Client) ([]byte, error) {
	buffer := manager.NewWriteAtBuffer([]byte{})

	downloader := manager.NewDownloader(s3Client)
	numBytes, err := downloader.Download(ctx, buffer, &s3.GetObjectInput{
		Bucket: aws.String(configuration.BucketName),
		Key:    aws.String(fileName),
	})
	if err != nil {
		return buffer.Bytes(), fmt.Errorf("Download error: %s", err)
	}

	if bytesReceived := int64(len(buffer.Bytes())); numBytes != bytesReceived {
		return buffer.Bytes(), fmt.Errorf("Incorrect number of bytes returned. Got %d, but expected %d", numBytes, bytesReceived)
	}

	saveFile(fileName, string(buffer.Bytes()))
	return buffer.Bytes(), nil
}
