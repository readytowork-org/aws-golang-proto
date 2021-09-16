package model

type IVSChannel struct {
	Name string;
	ChannelType string;
	EnableAuthorization bool;
}

type IVSChannelInput struct {
  StreamName  string `json:"streamName" binding:"required"`
}
