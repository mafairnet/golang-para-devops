package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

func createEC2(ctx context.Context, region string) (string, error) {

	//Cargamos la configuracion de los archovos configuration y credentials despues de haber isntalado la AWS CLI
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		return "", fmt.Errorf("LoadDefaultConfig error: %s", err)
	}

	//Creamos un objeto para poder acceder a los metodos del SDK para EC2
	ec2Client := ec2.NewFromConfig(cfg)

	//Para poder conectarse al EC2 se necseitan archivos de llave, si ya tenemos uno creado lo obtenemos
	existingKeyPairs, err := ec2Client.DescribeKeyPairs(ctx, &ec2.DescribeKeyPairsInput{
		KeyNames: []string{configuration.KeyFileName},
		// or:
		//Filters: []types.Filter{
		//	{
		//		Name:   aws.String("key-name"),
		//		Values: []string{"go-aws-ec2"},
		//	},
		//},
	})
	//Si hay algun error lo mostramos
	if err != nil && !strings.Contains(err.Error(), "InvalidKeyPair.NotFound") {
		return "", fmt.Errorf("DescribeKeyPairs error: %s", err)
	}

	//En caso contrario, creamos el archivo llave si no existe
	if existingKeyPairs == nil || len(existingKeyPairs.KeyPairs) == 0 {
		keyPair, err := ec2Client.CreateKeyPair(ctx, &ec2.CreateKeyPairInput{
			KeyName: aws.String(configuration.KeyFileName),
		})
		if err != nil {
			return "", fmt.Errorf("CreateKeyPair error: %s", err)
		}

		err = os.WriteFile(configuration.DefaulRegion, []byte(*keyPair.KeyMaterial), 0600)
		if err != nil {
			return "", fmt.Errorf("WriteFile (keypair) error: %s", err)
		}
	}

	//Obtenemos la configuracion de la imagen a partir de la cual queremos crear la EC2
	describeImages, err := ec2Client.DescribeImages(ctx, &ec2.DescribeImagesInput{
		Filters: []types.Filter{
			{
				Name:   aws.String("name"),
				Values: []string{configuration.ImageName},
			},
			{
				Name:   aws.String("virtualization-type"),
				Values: []string{configuration.ImageVirtualizationType},
			},
		},
		Owners: []string{configuration.ImageOwner}, // see https://ubuntu.com/server/docs/cloud-images/amazon-ec2
	})

	//Si no existe marcamos un error
	if err != nil {
		return "", fmt.Errorf("DescribeImages error: %s", err)
	}
	if len(describeImages.Images) == 0 {
		return "", fmt.Errorf("describeImages has empty length (%d)", len(describeImages.Images))
	}

	//Creamos la instancia con la imagen a utilizar y el tipo de instancia
	runInstance, err := ec2Client.RunInstances(ctx, &ec2.RunInstancesInput{
		ImageId:      describeImages.Images[0].ImageId,
		InstanceType: types.InstanceTypeT3Micro,
		KeyName:      aws.String(configuration.KeyFileName),
		MinCount:     aws.Int32(1),
		MaxCount:     aws.Int32(1),
	})

	//Si hay algun error lo mostramos
	if err != nil {
		return "", fmt.Errorf("RunInstance error: %s", err)
	}
	if len(runInstance.Instances) == 0 {
		return "", fmt.Errorf("RunInstance has empty length (%d)", len(runInstance.Instances))
	}

	//Retornamos el ID de la instancia para referencia futura
	return *runInstance.Instances[0].InstanceId, nil
}
