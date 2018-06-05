package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// custom structures to format the message for slack
type SlackMessage struct {
	Text        string       `json:"text"`
	Attachments []Attachment `json:"attachments"`
}

type Attachment struct {
	Text  string `json:"text"`
	Color string `json:"color"`
	Title string `json:"title"`
}

// reads in SNS events and executes functions that send message to slack
func handler(ctx context.Context, snsEvent events.SNSEvent) {
	for _, record := range snsEvent.Records {
		snsRecord := record.SNS

		fmt.Printf("[%s %s] Message = %s \n", record.EventSource, snsRecord.Timestamp, snsRecord.Message)
		slackMessage := buildSlackMessage(snsRecord)
		postToSlack(slackMessage)
		log.Println("Notification has been sent")
	}
}

// formats the message for slack
func buildSlackMessage(message events.SNSEntity) SlackMessage {
	return SlackMessage{
		Text: fmt.Sprintf("*%s*", message.Subject),
		Attachments: []Attachment{
			Attachment{
				Text:  message.Message,
				Color: "danger",
				Title: "Reason",
			},
		},
	}
}

// sends the formatted message to the slack incoming webhook
func postToSlack(message SlackMessage) error {
	client := &http.Client{}
	data, err := json.Marshal(message)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", os.Getenv("SLACK_WEBHOOK"), bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		fmt.Println(resp.StatusCode)
		return err
	}

	return nil
}

func main() {
	lambda.Start(handler)
}
