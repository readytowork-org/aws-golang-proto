package model

type IVSChannel struct {
	Name string;
	ChannelType string;
	EnableAuthorization bool;
	RecordingConfigurationArn string;
}

type IVSChannelInput struct {
  StreamName  string `json:"streamName" binding:"required"`
}
