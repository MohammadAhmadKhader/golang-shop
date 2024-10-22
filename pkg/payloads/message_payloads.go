package payloads

import (
	"slices"
	"strings"

	"main.go/pkg/models"
)

type CreateMessage struct {
	From     uint   `json:"from" validate:"required,min=1"`
	To       uint   `json:"to" validate:"required,min=1"`
	Content  string `json:"content" validate:"required,min=1,max=256"`
	Status   string `json:"status" validate:"omitempty,oneof=Sent"`
}

type UpdateMessage struct {
	Id uint `json:"id" validate:"required,min=1"`
	Content  string `json:"content" validate:"min=1,max=256"`
	Status   string `json:"status" validate:"omitempty,oneof=Sent Delivered Seen"`
}

type DeleteMessage struct {
	Id  uint `json:"id" validate:"required,min=1"`
}

func (m *UpdateMessage) TrimStrs() *UpdateMessage {
	if m != nil {
		m.Content = strings.Trim(m.Content, " ")
		m.Status = strings.Trim(m.Status, " ")
	}
	return m
}

// this used for websocket payload, id was set to 0 so it does not throw error during db update process.
func (um *UpdateMessage) ToModel() *models.Message {
	return &models.Message{
		Identifier: models.Identifier{ID: 0},
		Content: um.Content,
		Status: um.Status,
	}
}

func (m *UpdateMessage) Exclude(selectedFields []string) []string {
	removedCols := map[string]any{}
	if m != nil {
		if m.Content == "" {
			removedCols["Content"] = 1
		}
		
		if m.Status == "" {
			removedCols["Status"] = 1
		}
	}

	selectedFields = slices.DeleteFunc(selectedFields, func(element string) bool {
		_, exists := removedCols[element]
		return exists
	})

	return selectedFields
}

func (m *CreateMessage) TrimStrs() *CreateMessage {
	if m != nil {
		m.Content = strings.Trim(m.Content, " ")
		m.Status = strings.Trim(m.Status, " ")
	}
	return m
}

func (cm *CreateMessage) ToModel() *models.Message {
	return &models.Message{
		From: cm.From,
		To: cm.To,
		Content: cm.Content,
		Status: cm.Status,
	}
}

