package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var configuration = getProgramConfiguration()

func main() {
	var (
		s3Client *s3.Client
		err      error
		out      []byte
	)

	fileNamePtr := flag.String("filename", "", "Nombre del archivo.")
	actionPtr := flag.String("action", "", "Accion a realizar.")

	flag.Parse()

	fileName := *fileNamePtr
	action := *actionPtr

	if fileName == "" || action == "" {
		log.Fatal("No option selected")
		os.Exit(1)
	}

	fmt.Printf("Region:%v\n", configuration.DefaulRegion)
	fmt.Printf("Bucket:%v\n", configuration.BucketName)
	fmt.Printf("File:%v\n", fileName)
	fmt.Printf("Action:%v\n", action)

	//Inicializamos un objeto para acceder a los metodos para S3 del SDK de AWS
	if s3Client, err = initS3Client(context.Background(), configuration.DefaulRegion); err != nil {
		fmt.Printf("initConfig error: %s", err)
		os.Exit(1)
	}

	//Creamos un S3 si no existe
	if err = createS3Bucket(context.Background(), s3Client); err != nil {
		fmt.Printf("createS3Bucket error: %s", err)
		os.Exit(1)
	}

	switch action {
	case "upload":
		//Subimos una archivo
		if err = uploadFileToS3(fileName, context.Background(), s3Client); err != nil {
			fmt.Printf("uploadFileToS3 error: %s", err)
			os.Exit(1)
		}
		fmt.Printf("Uploaded file.\n")

	case "download":
		//Bajamos un archivo
		if out, err = downloadFileFromS3(fileName, context.Background(), s3Client); err != nil {
			fmt.Printf("uploadFileToS3 error: %s", err)
			os.Exit(1)
		}
		fmt.Printf("Downloaded file with contents: %s", out)
	default:
		fmt.Printf("No action selected")
	}
}
