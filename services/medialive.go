package services

import (
	"context"
	"fmt"
	"strings"

	"aws-golang-proto/model"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/medialive"
	"github.com/aws/aws-sdk-go-v2/service/medialive/types"
	msTypes "github.com/aws/aws-sdk-go-v2/service/mediastore/types"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/google/uuid"
)

type mediaLiveServices struct {
	ML *medialive.Client
}

func NewMediaLiveService(config aws.Config) *mediaLiveServices{
	return &mediaLiveServices{
		ML:medialive.New(medialive.Options{Credentials: config.Credentials,Region: *aws.String(endpoints.ApNortheast1RegionID)}),
	}
}

type MediaLiveServices interface{
	// inputs
	CreateInput(params model.Input)
	DeleteInput(inputId string) (*medialive.DeleteInputOutput,error)
	DescribeInput(inputId string) (*medialive.DescribeInputOutput,error)
	// input security groups
	ListInputSecurityGroups() (*medialive.ListInputSecurityGroupsOutput,error)
	// channels
	CreateChannel(
			channelName string, 
			channelClass string, 
			input types.Input, 
			container msTypes.Container,
		)(*medialive.CreateChannelOutput,error)
	StartChannel(channelId string) (*medialive.StartChannelOutput,error)
	StopChannel(channelId string) (*medialive.StopChannelOutput,error)
	DeleteChannel(channelId string) (*medialive.DeleteChannelOutput,error)
	DescribeChannel(channelId string) (*medialive.DescribeChannelOutput,error)
}

var ctx = context.Background()

// inputs
func (mlClient *mediaLiveServices) CreateInput(params model.Input) (*medialive.CreateInputOutput,error) {
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

func (mlClient *mediaLiveServices) DescribeInput(inputId string) (*medialive.DescribeInputOutput,error){
	return mlClient.ML.DescribeInput(ctx,&medialive.DescribeInputInput{
		InputId: aws.String(inputId),
	})
}


// input security groups
func (mlClient *mediaLiveServices) ListInputSecurityGroups() (*medialive.ListInputSecurityGroupsOutput,error){
	return mlClient.ML.ListInputSecurityGroups(ctx, &medialive.ListInputSecurityGroupsInput{})
}

// channels
func (mlClient *mediaLiveServices) CreateChannel(
		channelName string, 
		channelClass string, 
		input types.Input, 
		container msTypes.Container,
	) (*medialive.CreateChannelOutput,error){
		var channel model.Channel;

		channel.Name = channelName;
		channel.ChannelClass = channelClass;

		channel.InputAttachments = []types.InputAttachment{
			{
				InputId: input.Id,
				InputAttachmentName: input.Name,
				InputSettings: &types.InputSettings{
					AudioSelectors: []types.AudioSelector{
						{
							Name: aws.String("appears_in_audio_description_audio_selector_name"),
						},
					},
					InputFilter: types.InputFilterAuto,
					DeblockFilter: types.InputDeblockFilterDisabled,
					DenoiseFilter: types.InputDenoiseFilterDisabled,
					FilterStrength: 1,
					SourceEndBehavior: types.InputSourceEndBehaviorContinue,
					Smpte2038DataPreference: types.Smpte2038DataPreferenceIgnore,
				},
			},
		}

		// output destination to mediastore & (later) S3
		// TODO : outputDestination of length 2 containing mediastore and S3 endpoints
		id := uuid.New()
		// destination url for mediastore
		protocol := "mediastoressl"
		host := strings.Split(*container.Endpoint,"://")[1];
		folderName := "GolangFolder";
		fileName := uuid.New().String()
		firstDestinationUrl := fmt.Sprintf("%v://%v/%v/%v",protocol,host,folderName,fileName)
		secondDestinationUrl := fmt.Sprintf("%v://%v/%v/%v_redundant",protocol,host,folderName,fileName)
		channel.Destinations = []types.OutputDestination{
			{
				Id: aws.String(strings.Split(id.String(),"-")[0]),
				Settings: []types.OutputDestinationSettings{
					{
						Url: aws.String(firstDestinationUrl),
					},
					{
						Url: aws.String(secondDestinationUrl),
					},
				},
			},
		}

		// Encoder Settings
		channel.EncoderSettings = &types.EncoderSettings{
			OutputGroups: []types.OutputGroup{{
				OutputGroupSettings: &types.OutputGroupSettings{
					HlsGroupSettings: &types.HlsGroupSettings{
						Destination: &types.OutputLocationRef{
							DestinationRefId: channel.Destinations[0].Id,
						},
					},
				},
				Outputs: []types.Output{
					{
						OutputSettings: &types.OutputSettings{
							HlsOutputSettings: &types.HlsOutputSettings{
								NameModifier: aws.String("_360"),
								HlsSettings: &types.HlsSettings{
									StandardHlsSettings: &types.StandardHlsSettings{
										M3u8Settings: &types.M3u8Settings{
											PmtPid: aws.String("480"),
											VideoPid: aws.String("481"),
											PcrControl: types.M3u8PcrControlPcrEveryPesPacket,
											AudioPids: aws.String("492-498"),
											TimedMetadataBehavior: types.M3u8TimedMetadataBehaviorNoPassthrough,
											TimedMetadataPid: aws.String("502"),
											NielsenId3Behavior: types.M3u8NielsenId3BehaviorNoPassthrough,
										},
									},
								},
							},
						},
						AudioDescriptionNames: aws.ToStringSlice([]*string{
							aws.String("audio_rlaktp"),
						}),
						VideoDescriptionName: aws.String("video_rfkdb9"),
					},
				},
			}},
			AudioDescriptions: []types.AudioDescription{
				{
					Name: aws.String("audio_rlaktp"),
					AudioSelectorName: aws.String("appears_in_audio_description_audio_selector_name"),
					CodecSettings: &types.AudioCodecSettings{
						AacSettings: &types.AacSettings{
							Bitrate: 128000,
							SampleRate: 48000,
							Profile: types.AacProfileLc,
							InputType: types.AacInputTypeNormal,
							RawFormat: types.AacRawFormatNone,
							RateControlMode: types.AacRateControlModeCbr,
							CodingMode: types.AacCodingModeCodingMode10,
						},
					},
				},
			},
			VideoDescriptions: []types.VideoDescription{
				{
					Name: aws.String("video_rfkdb9"),
					Width: 480,
					Height: 360,
					Sharpness: 50,
					ScalingBehavior: types.VideoDescriptionScalingBehaviorDefault,
					RespondToAfd: types.VideoDescriptionRespondToAfdNone,
					CodecSettings: &types.VideoCodecSettings{
						H264Settings: &types.H264Settings{
							Bitrate: 500000,	// 500kbps
							FramerateControl: types.H264FramerateControlInitializeFromSource,
						},
					},
				},
			},
			TimecodeConfig: &types.TimecodeConfig{
				Source: types.TimecodeConfigSourceEmbedded,
			},
		}

		channel.InputSpecification = &types.InputSpecification{
			Codec: types.InputCodecAvc,
			Resolution: types.InputResolutionSd,
			MaximumBitrate: types.InputMaximumBitrateMax20Mbps,
		}

		return mlClient.ML.CreateChannel(ctx,&medialive.CreateChannelInput{
			Name: aws.String(channel.Name),
			InputAttachments: channel.InputAttachments,
			Destinations: channel.Destinations,
			EncoderSettings: channel.EncoderSettings,
			InputSpecification: channel.InputSpecification,
		})
}

func (mlClient *mediaLiveServices) StartChannel(channelId string) (*medialive.StartChannelOutput,error){
	return mlClient.ML.StartChannel(ctx,&medialive.StartChannelInput{
		ChannelId: aws.String(channelId),
	})
}

func (mlClient *mediaLiveServices) StopChannel(channelId string) (*medialive.StopChannelOutput,error){
	return mlClient.ML.StopChannel(ctx,&medialive.StopChannelInput{
		ChannelId: aws.String(channelId),
	})
}

func (mlClient *mediaLiveServices) DeleteChannel(channelId string) (*medialive.DeleteChannelOutput,error){
	return mlClient.ML.DeleteChannel(ctx,&medialive.DeleteChannelInput{
		ChannelId: aws.String(channelId),
	})
}

func (mlClient *mediaLiveServices) DescribeChannel(channelId string) (*medialive.DescribeChannelOutput,error){
	return mlClient.ML.DescribeChannel(ctx,&medialive.DescribeChannelInput{
		ChannelId: aws.String(channelId),
	})
}
