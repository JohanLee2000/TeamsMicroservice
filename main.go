// Created by Johan Lee - June 2023
/* This file takes an assumed json input file from the command line argument 'go run main.go <inputfile.json>', parses
it's data and sets the fields accordingly to a Microsoft Teams Message Card to be sent to a Teams channel. Specifically,
the 3 required fields are the webhookURL(which channel to send), title and text of message to send. */
//POSTMAN TEST
//webhookURL := "https://compassion.webhook.office.com/webhookb2/bb521093-fc98-488f-8876-af892e9649b4@d3c08eb3-7c4c-4a01-bc01-6c67b56c8da9/IncomingWebhook/34c22a9c8b004f74b8fc50cfe3754e89/93c9fa28-01e3-42ea-bab7-9264721a3fad"
//DEVOPS MONITOR
//webhookURL := "https://compassion.webhook.office.com/webhookb2/bb521093-fc98-488f-8876-af892e9649b4@d3c08eb3-7c4c-4a01-bc01-6c67b56c8da9/IncomingWebhook/ed37e28fcaec4e6cbbbe51a537e3877c/93c9fa28-01e3-42ea-bab7-9264721a3fad"

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	messagecard "teams_listening_service/messageCard"
)

type jsonMsg struct {
	WebhookURL string `json:"webhookURL"`
	Title      string `json:"title"`
	Text       string `json:"text"`
}

func main() {

	if len(os.Args) < 2 {
		log.Fatal("Usage: go run main.go <inputfile.json>")
	}

	//Take a JSON input file from command line argument
	inputFile := os.Args[1]
	jsonFile, err := os.Open(inputFile)
	if err != nil {
		log.Fatal("Failed to open json\n", err)
	}

	defer jsonFile.Close()

	byteValue, _ := io.ReadAll(jsonFile)

	var jsonMessage jsonMsg

	err = json.Unmarshal(byteValue, &jsonMessage)
	if err != nil {
		log.Fatal("Failed to parse json\n", err)
	}

	webhookURL := jsonMessage.WebhookURL
	if webhookURL == "" {
		log.Fatal("error: Missing webhook URL\n")
	}

	//There is already a check in message_card.go under Validate() to check for empty text and title fields during send process
	titleData := jsonMessage.Title

	textData := jsonMessage.Text

	//Initialize a Microsoft Teams client
	client := messagecard.CreateTeamsClient()

	//Set message card with contents here
	card := messagecard.CreateMessageCard()
	card.Title = titleData
	card.Text = textData

	//Send the message
	if err := client.Send(webhookURL, card); err != nil {
		log.Printf("failed to send message: %v", err)
		os.Exit(1)
	}
	fmt.Printf("\nThis is the webhookURL: %s\nThis is the titleData: %s\nThis is the textData: %s\n", webhookURL, titleData, textData)
}
