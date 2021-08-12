package services

import (
	"aws-golang-proto/utils"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/mediastore"
	awsconfig "github.com/tkuchiki/aws-sdk-go-config"
)

type mediaStoreServices struct {
	MS *mediastore.MediaStore
}

func NewMediaStoreService() *mediaStoreServices{
	var cfg aws.Config

	accessKey := utils.GetEnvWithKey("ACCESS_KEY")
	secretKey := utils.GetEnvWithKey("SECRET_KEY")
	region := utils.GetEnvWithKey("REGION")

	cfg.Credentials = awsconfig.NewCredentials(awsconfig.Option{
		AccessKey: accessKey,
		SecretKey: secretKey,
	})

	mySession := session.Must(session.NewSession(&cfg))
	return &mediaStoreServices{
		MS:mediastore.New(mySession,aws.NewConfig().WithCredentials(cfg.Credentials).WithRegion(region)),
	}
}

type MediaStoreServices interface{
	DescribeContainer(containerName string) (*mediastore.DescribeContainerOutput,error)
}

func (ms *mediaStoreServices) DescribeContainer(containerName string) (*mediastore.DescribeContainerOutput,error) {
	return ms.MS.DescribeContainer(&mediastore.DescribeContainerInput{ContainerName: aws.String(containerName)})
}