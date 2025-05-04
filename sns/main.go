package main
import (
	"fmt"
	"log"
	"time"
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sns"	
)
func main() {
	// load aws config
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("ap-south-1"))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}
	client := sns.NewFromConfig(cfg)
	topicArn := "arn:aws:sns:ap-south-1:216321913216485:myttest"

	// create a subscription to the topic
	subscriptionsInput := &sns.SubscribeInput{
		TopicArn: aws.String(topicArn),
		Protocol: aws.String("email"),
		Endpoint: aws.String("testumesh@yopmail.com"),
	}
	subs_result, err := client.Subscribe(context.TODO(), subscriptionsInput)
	if err != nil {	
		log.Fatalf("failed to subscribe to topic, %v", err)
	}
	fmt.Printf("Subscription ARN: %s\n", *subs_result.SubscriptionArn)
	fmt.Printf("Subscription ID: %s\n", *subs_result.SubscriptionArn)
	fmt.Println("=======================================")


	// Get the topic subscriptions
	list_subscriptionsInput := &sns.ListSubscriptionsByTopicInput{
		TopicArn: aws.String(topicArn),
	}
	subscriptionsOutput, err := client.ListSubscriptionsByTopic(context.TODO(), list_subscriptionsInput)
	if err != nil {
		log.Fatalf("failed to list subscriptions, %v", err)
	}
	for _, subscription := range subscriptionsOutput.Subscriptions {
		fmt.Printf("Subscription ARN: %s\n", *subscription.SubscriptionArn)
		fmt.Printf("Protocol: %s\n", subscription.Protocol)
		fmt.Printf("Endpoint: %s\n", *subscription.Endpoint)
		fmt.Println("=======================================")
	}



	// send a message to the topic
	messageBody := "Hello, this is a test message v2 from Go SDK! at " + time.Now().String()
	err = sendMessageToSNS(client, topicArn, messageBody)
	if err != nil {
		log.Fatalf("failed to send message, %v", err)
	}	

}
func sendMessageToSNS(client *sns.Client, topicArn string, messageBody string) error {
	// Create the input for the Publish API call
	input := &sns.PublishInput{
		TopicArn: aws.String(topicArn),
		Message:  aws.String(messageBody),
	}
	resp , err := client.Publish(context.TODO(), input)
	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}
	fmt.Printf("Message ID: %s\n", *resp.MessageId)
	return nil
}