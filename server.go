package main

import (
	"aws-golang-proto/services"
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
)

func main() {
	cfg,err := config.LoadDefaultConfig(context.TODO(),config.WithSharedConfigProfile("default"))	// loads from ~/.aws/credentials
	
	if err != nil {
		log.Fatal("Failed to load aws configuration")
	}

	mlService := services.NewMediaLiveService(cfg)
	// msService := services.NewMediaStoreService(cfg)

	// container, err := msService.DescribeContainer("ProgrammaticContainer")

	// if err != nil{
	// 	log.Println("Error fetching container info : ",err)
	// } else {
	// 	log.Println(*(container.Container.Endpoint))
	// }

	newChannel,err := mlService.CreateChannel()
	
	if err != nil {
		log.Fatal("Failed to create channel, error : ",err)
	} else {
		log.Println("Channel created : ",newChannel.Channel.Name)
	}


	channelDescription,err := mlService.DescribeChannel(*(newChannel.Channel.Id))
	
	if err != nil {
		log.Fatal("Failed to fetch the detail of the channel, error : ",err)
	} else {
		log.Printf("Channel Info :\n\t\t\tName : %v\n\t\t\tState : %v",*(channelDescription.Name), channelDescription.State)
		if channelDescription.State == "IDLE"{
			startingChannel, err := mlService.StartChannel(*(channelDescription.Id))
			if err != nil {
				log.Fatal("Failed to start the channel, error : ",err)
			} else {
				log.Println("Channel State : ",startingChannel.State)
			}
		}
	}

	// Stop Channel
	// stoppingChannel, err := mlService.StopChannel(*channelDescription.Id)
	// if err != nil {
	// 	log.Fatal("Failed to stop the channel, error : ",err)
	// } else {
	// 	log.Println("Channel State : ",stoppingChannel.State)
	// }


	// ch,err := ML.ListChannels(context.Background(),&medialive.ListChannelsInput{})
	// if err != nil {
	// 	log.Fatal("Failed to fetch the list of channels")
	// }

	// var channels medialive.ListChannelsOutput
	// channels.Channels = ch.Channels

	// for _, channel := range(ch.Channels){
	// 	log.Println(*(channel.Name))
	// }

	
	// isg , _ := mlService.ListInputSecurityGroups()
	// var inputParams = services.IMediaLiveInput{
	// 	Name: "DynamicInputFromApp",
	// 	Type: "RTMP_PUSH",
	// 	InputSecurityGroupsId: []*string{isg.InputSecurityGroups[0].Id},
	// 	DestinationUrl: []string{"DynamicInpA/inpA","DynamicInpB/inpB"},
	// }	
	// createdInput, err := mlService.CreateInput(inputParams)

	// if err != nil{
	// 	log.Fatal("Failed to create input, Error: ",err)
	// }
	
	// log.Println("Created input = ", *(createdInput.Input.Id))

	// deletedInput, err := mlService.DeleteInput(*(createdInput.Input.Id))
	// if err != nil{
	// 	log.Fatal("Failed to create input, Error: ",err)
	// }

	// log.Println("Deleted Input = ",deletedInput)
}