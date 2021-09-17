package main

import (
	"aws-golang-proto/helper"
	"aws-golang-proto/model"
	"aws-golang-proto/services"
	"context"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/config"
	ivsTypes "github.com/aws/aws-sdk-go-v2/service/ivs/types"
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
	ivsService := services.NewInteractiveVideoService(cfg)

	httpRouter := gin.Default()
	httpRouter.GET("/ping",func(c *gin.Context) {
		c.JSON(http.StatusOK,gin.H{
			"message":"pong",
		})
	})

	// ------------------------		IVS BEGINS		---------------------------	//

	httpRouter.POST("/createIVSChannel", func(c *gin.Context){
		var channelInp model.IVSChannelInput

		c.ShouldBindJSON(&channelInp)
			
		channel, err := ivsService.CreateChannel(model.IVSChannel{
			Name: channelInp.StreamName,
			ChannelType: string(ivsTypes.ChannelTypeStandardChannelType),
			EnableAuthorization: false,
		})

		if err != nil {
			c.AbortWithError(http.StatusBadRequest,err)
			return
		}

		streamKeysList, err := ivsService.ListStreamKey(*channel.Channel.Arn)

		if err != nil {
			log.Println("Error getting list of stream keys\n ",err)
			c.AbortWithError(http.StatusBadRequest,err)
			return
		}

		streamKeyOut, err := ivsService.GetStreamKey(*streamKeysList.StreamKeys[0].Arn)

		if err != nil {
			log.Printf("Error getting stream key for arn %v \n %v",*streamKeysList.StreamKeys[0].Arn,err)
			c.AbortWithError(http.StatusBadRequest,err)
			return
		}

		// store channel arn in database and use uuid instead to uiniquly identify channels created		
		c.JSON(http.StatusOK,gin.H{
			"channelARN":&channel.Channel.Arn,
			"ingetsUrl": &channel.Channel.IngestEndpoint,
			"streamKeyARN": &streamKeyOut.StreamKey.Arn,
			"streamKey": &streamKeyOut.StreamKey.Value,
		})
	})

	httpRouter.GET("/allLiveStreams", func(c *gin.Context){
		nextPage, ok := c.GetQuery("nextToken")
		nextToken := "" 

		if ok {
			nextToken = nextPage
		}

		liveStreams, err := ivsService.ListLiveChannels(nextToken)

		if err != nil {
			c.AbortWithError(http.StatusBadRequest,err)
			return
		}
		
		c.JSON(http.StatusOK,gin.H{
			"liveStreams": liveStreams.Streams,
			"nextToken": liveStreams.NextToken,
		})
	})	

	httpRouter.GET("/streamUrl", func(c *gin.Context){
		channelARN, ok := c.GetQuery("channelARN")
		if !ok {
			c.JSON(http.StatusBadRequest,gin.H{
				"message":"Channel ARN cannot be empty",
			})
			return
		}

		streamUrl, err := ivsService.GetPlaybackURL(channelARN)

		if err != nil {
			c.AbortWithError(http.StatusBadRequest,err)
			return
		}
		
		c.JSON(http.StatusOK,gin.H{
			"streamUrl": streamUrl,
		})
	})

	httpRouter.DELETE("/deleteIVSChannel", func(c *gin.Context){
		channelARN, ok := c.GetQuery("channelARN")
		if !ok {
			c.JSON(http.StatusBadRequest,gin.H{
				"message":"Channel ARN cannot be empty",
			})
		}

		streamUrl, err := ivsService.DeleteChannel(channelARN)

		if err != nil {
			c.AbortWithError(http.StatusBadRequest,gin.Error{
				Err: err,
			})
			return
		}
		
		c.JSON(http.StatusOK,gin.H{
			"streamUrl": streamUrl,
		})
	})
	

	// ------------------------		IVS ENDS		-----------------------------	//



	httpRouter.GET("/startStream", func(c *gin.Context) {
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
					log.Println("Channel created : ",newChannel.Channel.Name)

					channelState := newChannel.Channel.State;

					// fetch channel description every 5 seconds
					clearInterval := helper.SetInterval(func() {
						channelDescription,err := mlService.DescribeChannel(*(newChannel.Channel.Id))
						if err != nil {
							log.Fatal("Failed to fetch the detail of the channel, error : ",err)
						} else {
							// log.Printf("Channel Info :\n\t\t\tName : %v\n\t\t\tState : %v",*(channelDescription.Name), channelDescription.State)
							log.Println("Desc, Channel State : ",channelDescription.State)
							channelState = channelDescription.State;
						}
					},5 * 1000,true);

					// run concurrent function that checks the channel state
					// if the cahnnel is running clear the interval and respond the rtmp input url
					go func () {
						for {
							if channelState == "CREATING" {
								// log.Println("Channel is being created...")
							} else if channelState == "IDLE" {
								_ ,err := mlService.StartChannel(*newChannel.Channel.Id)
								if err != nil {
									log.Println("Error starting the channel...",err)
								} else {
									// log.Println("Channel is starting...")
								}
							} else if channelState == "STARTING" {
								// log.Println("Channel is starting...")
							} else if channelState == "RUNNING" {
								clearInterval <- true;
								log.Println("Channel is running...")
								c.JSON(http.StatusOK,gin.H{
									"channelID": newChannel.Channel.Id,
									"inputID": createdInput.Input.Id,
									"streamUrls": []types.InputDestination{
										{
											Url: createdInput.Input.Destinations[0].Url,
										},
										{
											Url: createdInput.Input.Destinations[1].Url,
										},
									},
								})
								return												
							}
						}
					}()
				}
			}
		}
	})

	httpRouter.GET("/streamInfo", func(c *gin.Context) {
		channelID, ok := c.Params.Get("channelID")
		if !ok {
			c.JSON(http.StatusBadRequest,gin.H{
				"message":"Channel id cannot be empty",
			})
		} 

		inputID, ok := c.Params.Get("inputID")
		if !ok {
			c.JSON(http.StatusBadRequest,gin.H{
				"message":"Input id cannot be empty",
			})
		} 

		input, err := mlService.DescribeInput(inputID)
		channel, err := mlService.DescribeChannel(channelID)
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{
				"message":"Error getting channel information for the channel id provided",
			})
		} else {
			c.JSON(http.StatusOK,gin.H{
				"channelID": channel.Id,
				"inputID": input.Id,
				"streamUrls": []types.InputDestination{
					{
						Url: input.Destinations[0].Url,
					},
					{
						Url: input.Destinations[1].Url,
					},
				},
			})
		}
	})

	httpRouter.GET("/stopStream", func(c *gin.Context) {
		params := c.Params
		
		channelID, err := params.Get("channelID")
		if err {
			log.Println("ID of the Channel to be stopped & deleted not provided");
			c.JSON(http.StatusBadRequest, gin.H{
				"message" : "ID of the Channel to be stopped & deleted not provided",
			})
		}

		inputID, err := params.Get("inputID")
		if err {
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