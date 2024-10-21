package websocket

import (
	"encoding/json"
	"fmt"
	"log"

	"main.go/pkg/models"
)

func HandleMessageCreate(event Event, client *Client) error {
	var userMessage UserMessagePayload
	err := json.Unmarshal(event.Payload, &userMessage)
	if err != nil {
		log.Println("Error during converting from json to user message payload")
		return err
	}

	msg := models.Message{
		From:   userMessage.From,
		To:     userMessage.To,
		Status: "sent",
	}
	fmt.Println(msg)

	return nil
}

func HandleMessageUpdate(event Event, client *Client) error {
	var userMessage UserMessagePayload
	err := json.Unmarshal(event.Payload, &userMessage)
	if err != nil {
		log.Println("Error during converting from json to user message payload")
		return err
	}

	msg := models.Message{
		From:   userMessage.From,
		To:     userMessage.To,
		Status: "sent",
	}
	fmt.Println(msg)

	return nil
}

func HandleMessageDelete(event Event, client *Client) error {
	var userMessage UserMessagePayload
	err := json.Unmarshal(event.Payload, &userMessage)
	if err != nil {
		log.Println("Error during converting from json to user message payload")
		return err
	}

	msg := models.Message{
		From:   userMessage.From,
		To:     userMessage.To,
		Status: "sent",
	}
	fmt.Println(msg)

	return nil
}