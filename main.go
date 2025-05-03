package main

import (
	"fmt"
	"log"
	"time"
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

func main() {
	//load aws config
	cfg,err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("ap-south-1"))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}
	client := sqs.NewFromConfig(cfg)
	queueUrl := "https://sqs.ap-south-1.amazonaws.com/216989136485/test-sqs"

	// Send a message to the queue
	messageBody := "Hello, this is a test message"
	err = sendMessageToSQS(client, queueUrl, messageBody)
	if err != nil {
		log.Fatalf("failed to send message, %v", err)
	}

	// Wait for a few seconds to ensure the message is sent
	time.Sleep(5 * time.Second)

	// Receive messages from the queue
	receiveMessageInput := &sqs.ReceiveMessageInput{
		QueueUrl: aws.String(queueUrl),
		MaxNumberOfMessages: 1,
		WaitTimeSeconds: 10,
	}
	msgResult, err := client.ReceiveMessage(context.TODO(), receiveMessageInput)
	if err != nil {
		log.Fatalf("failed to receive message, %v", err)
	}
	if len(msgResult.Messages) > 0 {
		for _, message := range msgResult.Messages {
			fmt.Printf("Message ID: %s\n", *message.MessageId)
			fmt.Printf("Message Body: %s\n", *message.Body)
			// Delete the message after processing
			deleteMessageInput := &sqs.DeleteMessageInput{
				QueueUrl: aws.String(queueUrl),
				ReceiptHandle: message.ReceiptHandle,
			}
			_, err := client.DeleteMessage(context.TODO(), deleteMessageInput)
			if err != nil {
				log.Fatalf("failed to delete message, %v", err)
			}
			fmt.Printf("Deleted message ID: %s\n", *message.MessageId)
		}
		fmt.Printf("Received message: %s\n", *msgResult.Messages[0].Body)
	} else {
		fmt.Println("No messages received")
	}	


}
// refactor the code for sending message to SQS queue

func sendMessageToSQS(client *sqs.Client, queueUrl string, messageBody string) error {
	sendMessageInput := &sqs.SendMessageInput{
		MessageBody: aws.String(messageBody),
		QueueUrl:    aws.String(queueUrl),
	}
	
	_, err := client.SendMessage(context.TODO(), sendMessageInput)
	if err != nil {
		return fmt.Errorf("failed to send message, %v", err)
	}
	fmt.Println("Message sent to SQS queue")
	return nil
}
