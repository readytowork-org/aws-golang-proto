package main

import (
	"aws-golang-proto/model"
	"aws-golang-proto/services"
	"context"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/medialive/types"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg,err := config.LoadDefaultConfig(context.TODO(),config.WithSharedConfigProfile("default"))	// loads from ~/.aws/credentials
	
	if err != nil {
		log.Fatal("Failed to load aws configuration")
	}

	mlService := services.NewMediaLiveService(cfg)
	msService := services.NewMediaStoreService(cfg)

	httpRouter := gin.Default()
	httpRouter.GET("/ping",func(c *gin.Context) {
		c.JSON(http.StatusOK,gin.H{
			"message":"pong",
		})
	})

	httpRouter.GET("/createStreamResources", func(c *gin.Context) {
		container, err := msService.DescribeContainer("ProgrammaticContainer") // TODO : container name should come from .env file

		if err != nil{
			log.Println("Error fetching container info : ",err)
		} else {
			log.Println(*(container.Container.Endpoint))
	
			isg , _ := mlService.ListInputSecurityGroups()	// In Production: Create a security group from console and security group id should come from .env
			var inputParams = model.Input{
				Name: "DynamicInputFromApp",
				Type: "RTMP_PUSH",
				InputSecurityGroupsId: []*string{isg.InputSecurityGroups[0].Id},
				DestinationUrl: []string{"DynamicInpA/inpA","DynamicInpB/inpB"},
			}	

			createdInput, err := mlService.CreateInput(inputParams)
			if err != nil{
				log.Fatal("Failed to create input, Error: ",err)
			} else {
				log.Println("Created input = ", *(createdInput.Input.Id))
			
				newChannel,err := mlService.CreateChannel("Golang Channel","STANDARD",*createdInput.Input,*container.Container)
	
				if err != nil {
					log.Fatal("Failed to create channel, error : ",err)
				} else {
					log.Println("Created channel : ",*(newChannel.Channel.Name))

					c.JSON(http.StatusOK,gin.H{
						"channelID": newChannel.Channel.Id,
						"inputID": createdInput.Input.Id,
						"state": newChannel.Channel.State,
					})
				}
			}
		}
	})

	httpRouter.GET("/startChannel",func(c *gin.Context) {
		queryParams := c.Request.URL.Query()
		channelID := queryParams.Get("channelID")
		if channelID == "" {
			c.JSON(http.StatusBadRequest,gin.H{
				"message":"Channel id cannot be empty",
			})
			return
		}

		channel, err := mlService.DescribeChannel(channelID)
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{
				"message":"Error getting channel information for the channel id provided",
			})
			return
		} else {
			if channel.State == "IDLE" {
				startingChannel, err := mlService.StartChannel(*channel.Id)
				if err != nil {
					c.JSON(http.StatusInternalServerError,gin.H{
						"message":"Error starting channel for the channel id provided",
					})
					return
				} else {
					c.JSON(http.StatusOK,gin.H{
						"channelID": startingChannel.Id,
						"inputID": startingChannel.InputAttachments[0].InputId,
						"state": startingChannel.State,
					})
				}
			}
		}
	})

	httpRouter.GET("/streamInfo", func(c *gin.Context) {
		queryParams := c.Request.URL.Query()
		channelID := queryParams.Get("channelID")
		log.Println(channelID)
		if channelID == "" {
			c.JSON(http.StatusBadRequest,gin.H{
				"message":"Channel id cannot be empty",
			})
			return
		} 

		inputID := queryParams.Get("inputID")
		log.Println(inputID)
		if inputID == "" {
			c.JSON(http.StatusBadRequest,gin.H{
				"message":"Input id cannot be empty",
			})
			return
		} 

		input, err := mlService.DescribeInput(inputID)

		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{
				"message":"Error getting input information for the input id provided",
			})
			return
		}
		
		channel, err := mlService.DescribeChannel(channelID)
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{
				"message":"Error getting channel information for the channel id provided",
			})
			return
		} else {
			if(channel.State != "RUNNING"){

				c.JSON(http.StatusOK,gin.H{
					"channelID": channel.Id,
					"inputID": input.Id,
					"state": channel.State,
				})
			} else {
			c.JSON(http.StatusOK,gin.H{
				"channelID": channel.Id,
				"inputID": input.Id,
				"state": channel.State,
				"streamUrls": []types.InputDestination{
					{
						Url: input.Destinations[0].Url,
					},
					{
						Url: input.Destinations[1].Url,
					},
				},
			})}
		}
	})

	httpRouter.GET("/stopStream", func(c *gin.Context) {
		params := c.Request.URL.Query()
		
		channelID := params.Get("channelID")
		if channelID == "" {
			log.Println("ID of the Channel to be stopped & deleted not provided");
			c.JSON(http.StatusBadRequest, gin.H{
				"message" : "ID of the Channel to be stopped & deleted not provided",
			})
		}

		inputID := params.Get("inputID")
		if inputID == "" {
			log.Println("ID of the Channel to be stopped & deleted not provided");
			c.JSON(http.StatusBadRequest, gin.H{
				"message" : "ID of the Input to be deleted not provided",
			})
		}

		// Stop Channel
		stoppingChannel, error := mlService.StopChannel(channelID)
		if error != nil {
			log.Fatal("Failed to stop the channel, error : ",error)
		} else {
			log.Println("Channel State : ",stoppingChannel.State)

			// Delete Channel
			deletedChannel, error := mlService.DeleteChannel(channelID);
			if error != nil {
				log.Fatal("Failed to delete the channel, error : ",error)
			} else {
				log.Println("Channel State : ",deletedChannel.State)
			}

			// Delete Input
			// rather than deleting the input, let the input exist and create only the channel with the existing input
			deletedInput, error := mlService.DeleteInput(inputID)
			if error != nil{
				log.Fatal("Failed to create input, Error: ",error)
			}

			log.Println("Deleted Input = ",deletedInput)
			return
		}
	})

	httpRouter.Run(":8000")
}