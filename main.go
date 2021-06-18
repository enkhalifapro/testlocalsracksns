package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"log"
)

// CustomEvent for lambda
type CustomEvent struct {
	ID   string
	Name string
}

var (
	awsRegion   string
	awsEndpoint string

	snsClient *sns.Client
)

// We use this interface to test the function using a mocked service.
type SNSSubscribeAPI interface {
	Subscribe(ctx context.Context,
		params *sns.SubscribeInput,
		optFns ...func(*sns.Options)) (*sns.SubscribeOutput, error)
}

func SubscribeTopic(c context.Context, api SNSSubscribeAPI, input *sns.SubscribeInput) (*sns.SubscribeOutput, error) {
	return api.Subscribe(c, input)
}


func init() {
	awsRegion ="us-east-1"// os.Getenv("AWS_REGION")
	awsEndpoint = "http://localhost:4566/"//os.Getenv("AWS_ENDPOINT")

	customResolver := aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
		if awsEndpoint != "" {
			return aws.Endpoint{
				PartitionID:   "aws",
				URL:           awsEndpoint,
				SigningRegion: awsRegion,
			}, nil
		}

		// returning EndpointNotFoundError will allow the service to fallback to it's default resolution
		return aws.Endpoint{}, &aws.EndpointNotFoundError{}
	})

	awsCfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(awsRegion),
		config.WithEndpointResolver(customResolver),
	)
	if err != nil {
		log.Fatalf("Cannot load the AWS configs: %s", err)
	}


	snsClient = sns.NewFromConfig(awsCfg)
}

func subscribeToSNSTopic(email string,topicARN string)  {
	input := &sns.SubscribeInput{
		Endpoint:              &email,
		Protocol:              aws.String("email"),
		ReturnSubscriptionArn: true, // Return the ARN, even if user has yet to confirm
		TopicArn:              &topicARN,
	}

	result, err := SubscribeTopic(context.TODO(), snsClient, input)
	if err != nil {
		fmt.Println("Got an error subscribing to the topic:")
		fmt.Println(err)
		return
	}

	fmt.Println(*result.SubscriptionArn)
}
func main() {

	subscribeToSNSTopic("smartsoft07@gmail.com","arn:aws:sns:us-east-1:000000000000:myTopic")
}

