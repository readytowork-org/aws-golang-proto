package services

import (
	// "aws-golang-proto/model"

	// "aws-golang-proto/model"
	// "log"

	"fmt"

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


func (s3Client *s3Services) GetObjList(bucketName string) (*s3.ListObjectsV2Output,error) {
	// accountId is the id of the root user (this is same for all other IAM users as well)
	// accountId can be extracted from aws-cli with following command
	// `aws sts get-caller-identity --query Account --output text`
	// should be stored in .env file
	accountId := "876923632685"
	return s3Client.S3.ListObjectsV2(ctx,&s3.ListObjectsV2Input{
		Bucket: aws.String(bucketName),
		Prefix: aws.String(fmt.Sprintf("ivs/v1/%v",accountId)), // ivs/v1/<aws user account id>
	})
}
