package services

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/medialive"
	"github.com/aws/aws-sdk-go-v2/service/medialive/types"
	"github.com/aws/aws-sdk-go/aws/endpoints"
)

type mediaLiveServices struct {
	ML *medialive.Client
}

func NewMediaLiveService(config aws.Config) *mediaLiveServices{
	return &mediaLiveServices{
		ML:medialive.New(medialive.Options{Credentials: config.Credentials,Region: *aws.String(endpoints.ApNortheast1RegionID)}),
	}
}


type IMediaLiveInput struct {
	Id string;
	Name string;
	Type string;
	DestinationUrl []string;
	InputSecurityGroupsId []*string
}

type MediaLiveServices interface{
	CreateInput(params IMediaLiveInput)
	DeleteInput(inputId string) (*medialive.DeleteInputOutput,error)
	ListInputSecurityGroups() (*medialive.ListInputSecurityGroupsOutput,error)
}

var ctx = context.Background()

func (mlClient *mediaLiveServices) CreateInput(params IMediaLiveInput) (*medialive.CreateInputOutput,error) {
	inputDestination := make([]types.InputDestinationRequest, 2);

	for i, dest := range(params.DestinationUrl) {
		inputDestination[i].StreamName = &dest
	}

	input, err := mlClient.ML.CreateInput(ctx,&medialive.CreateInputInput{
		Name: aws.String(params.Name),
		Type: types.InputType(params.Type),
		InputSecurityGroups: aws.ToStringSlice(params.InputSecurityGroupsId),
		Destinations: inputDestination,
	});

	return input, err
}

func (mlClient *mediaLiveServices) DeleteInput(inputId string) (*medialive.DeleteInputOutput,error){
	return mlClient.ML.DeleteInput(ctx,&medialive.DeleteInputInput{InputId: aws.String(inputId)})
}

func (mlClient *mediaLiveServices) ListInputSecurityGroups() (*medialive.ListInputSecurityGroupsOutput,error){
	return mlClient.ML.ListInputSecurityGroups(ctx, &medialive.ListInputSecurityGroupsInput{})
}