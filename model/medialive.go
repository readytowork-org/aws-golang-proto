package model

import "github.com/aws/aws-sdk-go-v2/service/medialive/types"

type Channel struct {
	Name string;
	ChannelClass string;
	InputAttachments []types.InputAttachment
	Destinations []types.OutputDestination
	EncoderSettings *types.EncoderSettings
	InputSpecification *types.InputSpecification
}


type Input struct {
	Id string;
	Name string;
	Type string;
	DestinationUrl []string;
	InputSecurityGroupsId []*string
}