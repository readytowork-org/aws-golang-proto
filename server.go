package main

import (
	"aws-golang-proto/services"
	"aws-golang-proto/utils"
	"log"
)

func main() {
	utils.LoadEnv()

	// cfg,err := config.LoadDefaultConfig(context.TODO(),config.WithSharedConfigProfile("default"))	// loads from ~/.aws/credentials
	
	// if err != nil {
	// 	log.Fatal("Failed to load aws configuration")
	// }

	// mlService := services.NewMediaLiveService(cfg)
	msService := services.NewMediaStoreService()

	container, err := msService.DescribeContainer("ProgrammaticContainer")

	if err != nil{
		log.Println("Error fetching container info : ",err)
	} else {
		log.Println(*(container.Container.Endpoint))
	}
	
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