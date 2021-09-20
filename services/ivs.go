package services

import (
	// "aws-golang-proto/model"

	"aws-golang-proto/model"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ivs"
	"github.com/aws/aws-sdk-go-v2/service/ivs/types"

	"github.com/aws/aws-sdk-go/aws/endpoints"
)

type ivsServices struct {
	IVS *ivs.Client
}


func NewInteractiveVideoService(config aws.Config) *ivsServices{
	return &ivsServices{
		IVS:ivs.New(ivs.Options{Credentials: config.Credentials,Region: *aws.String(endpoints.UsWest2RegionID)}),
	}
}


// channel
// creates a channel and an associated stream key
func (ivsClient *ivsServices) CreateChannel(params model.IVSChannel) (*ivs.CreateChannelOutput, error) {

	recConfig, err := ivsClient.IVS.GetRecordingConfiguration(ctx,&ivs.GetRecordingConfigurationInput{
		Arn: aws.String("arn:aws:ivs:us-west-2:876923632685:recording-configuration/djZinlTn6F38"),
	})

	if err != nil {
		log.Println("Error getting recording configurations", err)
	}

	return ivsClient.IVS.CreateChannel(ctx,&ivs.CreateChannelInput{
		Name: aws.String(params.Name),
		Authorized: params.EnableAuthorization,
		LatencyMode: types.ChannelLatencyModeLowLatency,
		Type: types.ChannelType(params.ChannelType),
		RecordingConfigurationArn: recConfig.RecordingConfiguration.Arn,
	})
}

func (ivsClient *ivsServices) StopStream(channelARN string) (*ivs.StopStreamOutput, error) {
	return ivsClient.IVS.StopStream(ctx,&ivs.StopStreamInput{
		ChannelArn: &channelARN,
	})
}

func (ivsClient *ivsServices) ListLiveChannels(nextToken string) (*ivs.ListStreamsOutput, error) {
	var listStreamInput ivs.ListStreamsInput
	if nextToken != "" {
		listStreamInput.NextToken = &nextToken
	}
	return ivsClient.IVS.ListStreams(ctx,&listStreamInput)
}

func (ivsClient *ivsServices) ListStreamKey(channelARN string) (*ivs.ListStreamKeysOutput, error) {
	return ivsClient.IVS.ListStreamKeys(ctx,&ivs.ListStreamKeysInput{
		ChannelArn: &channelARN,
	})
}

// only one key can be associated with one channel, so before creating a stream key for a channel, delete existing key first
func (ivsClient *ivsServices) CreateStreamKey(channelARN string) (*ivs.CreateStreamKeyOutput, error) {
	log.Println("Channel ARN : ",channelARN)
	return ivsClient.IVS.CreateStreamKey(ctx,&ivs.CreateStreamKeyInput{
		ChannelArn: &channelARN,
	})
}

func (ivsClient *ivsServices) GetStreamKey(streamKeyArn string) (*ivs.GetStreamKeyOutput, error) {
	return ivsClient.IVS.GetStreamKey(ctx,&ivs.GetStreamKeyInput{
		Arn: &streamKeyArn,
	})
}


func (ivsClient *ivsServices) GetPlaybackURL(channelArn string) (string, error) {
	channel, err := ivsClient.IVS.GetChannel(ctx,&ivs.GetChannelInput{
		Arn: aws.String(channelArn),
	})

	if err != nil {
		log.Println("Error getting channel information", err)
	}
	
	return *channel.Channel.PlaybackUrl, err
}

func (ivsClient *ivsServices) DeleteChannel(channelArn string) (*ivs.DeleteChannelOutput, error) {
	return ivsClient.IVS.DeleteChannel(ctx,&ivs.DeleteChannelInput{
		Arn: aws.String(channelArn),
	})
}
