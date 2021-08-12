package services

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/mediastore"
	"github.com/aws/aws-sdk-go/aws/endpoints"
)

type mediaStoreServices struct {
	MS *mediastore.Client
}

func NewMediaStoreService(config aws.Config) *mediaStoreServices{
	return &mediaStoreServices{
		MS:mediastore.New(mediastore.Options{Credentials:config.Credentials,Region: endpoints.ApNortheast1RegionID}),
	}
}

type MediaStoreServices interface{
	DescribeContainer(containerName string) (*mediastore.DescribeContainerOutput,error)
}

func (ms *mediaStoreServices) DescribeContainer(containerName string) (*mediastore.DescribeContainerOutput,error) {
	return ms.MS.DescribeContainer(ctx, &mediastore.DescribeContainerInput{ContainerName: aws.String(containerName)})
}