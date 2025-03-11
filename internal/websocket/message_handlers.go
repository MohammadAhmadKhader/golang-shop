package websocket

import (
	"encoding/json"
	"fmt"
	"log"

	"main.go/pkg/utils"
	"main.go/types"
)

func HandleMessageCreate(event Event, client *Client) error {
	var userMessage MessagePayload
	err := json.Unmarshal(event.Payload, &userMessage)
	if err != nil {
		log.Println("Error during converting from json to user message payload")
		return err
	}

	if userMessage.From != *client.id {
		return fmt.Errorf("sender id does not match the client id")
	}

	msg := userMessage.ToCreatePayload()
	msg.Status = string(types.Sent)

	err = utils.ValidateStruct(msg)
	if err != nil {
		log.Println("invalid message payload")
		return err
	}
	
	validMsg := msg.TrimStrs().ToModel()
	err = GlobalManager.store.CreateMessage(*validMsg)
	if err != nil {
		log.Println("Error during creating message")
		return err
	}

	return nil
}

func HandleMessageUpdate(event Event, client *Client) error {
	var userMessage MessagePayload
	err := json.Unmarshal(event.Payload, &userMessage)
	if err != nil {
		log.Println("Error during converting from json to user message payload")
		return err
	}

	if userMessage.From != *client.id {
		return fmt.Errorf("sender id does not match the client id")
	}

	msg := userMessage.ToUpdatePayload()

	err = utils.ValidateStruct(msg)
	if err != nil {
		log.Println("invalid message payload")
		return err
	}
	
	changes := msg.TrimStrs().ToModel()
	_,err = GlobalManager.store.UpdateMessage(userMessage.Id,*changes, &msg)
	if err != nil {
		log.Println("Error during creating message")
		return err
	}

	return nil
}

func HandleMessageDelete(event Event, client *Client) error {
	var delPayload DeleteMessagePayload
	err := json.Unmarshal(event.Payload, &delPayload)
	if err != nil {
		log.Println("Error during converting from json to user message payload")
		return err
	}

	payload := delPayload.ToDeletePayload()
	err = utils.ValidateStruct(payload)
	if err != nil {
		log.Println("invalid message payload")
		return err
	}
	
	err = GlobalManager.store.DeleteMessage(payload.Id)
	if err != nil {
		log.Println("Error during creating message")
		return err
	}

	return nil
}

func HandleMessageUpdateStatus(event Event, client *Client) error {
	var updateStatusPayload MessageStatusPayload
	err := json.Unmarshal(event.Payload, &updateStatusPayload)
	if err != nil {
		log.Println("Error during converting from json to user message status updating payload")
		return err
	}

	err = utils.ValidateStruct(updateStatusPayload)
	if err != nil {
		log.Println("invalid message update status payload")
		return err
	}
	
	err = GlobalManager.store.UpdateMessageStatus(updateStatusPayload.Id, updateStatusPayload.Status)
	if err != nil {
		log.Println("Error during updating message status")
		return err
	}

	return nil
}