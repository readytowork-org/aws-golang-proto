package services

import (
	// "aws-golang-proto/model"

	// "aws-golang-proto/model"
	// "log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	// "github.com/aws/aws-sdk-go-v2/service/s3/types"

	"github.com/aws/aws-sdk-go/aws/endpoints"
)

type s3Services struct {
	S3 *s3.Client
}


func NewS3Service(config aws.Config) *s3Services{
	return &s3Services{
		S3:s3.New(s3.Options{Credentials: config.Credentials,Region: *aws.String(endpoints.UsWest2RegionID)}),
	}
}


func (s3Client *s3Services) GetObjList() (*s3.ListObjectsV2Output,error) {
	return s3Client.S3.ListObjectsV2(ctx,&s3.ListObjectsV2Input{
		Bucket: aws.String("ivs-console-stream-archive"),
		Prefix: aws.String("ivs/v1/876923632685"), // ivs/v1/<aws user account id>
	})
}
