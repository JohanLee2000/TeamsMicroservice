// Created by Johan Lee - June 2023

/* This file contains the implementation of the Message Card, a legacy actionable message card
used by Office 365 or Microsoft Teams connectors. The basic methods here are to support a
simple message with a title and text. Future implementations can include adding images to the
card or a potential clickable action that opens up a uri */

package messageCard

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

// MessageCard struct and fields
type MessageCard struct {
	//Required fields
	Type string `json:"@type" yaml:"@context"`

	Context string `json:"@context" yaml:"@context"`

	//Title of message, at top of the card
	Title string `json:"title,omitempty" yaml:"title,omitempty"`

	//Text of message, displays content of message
	Text string `json:"text,omitempty" yaml:"text,omitempty"`

	//Color of message card
	Color string `json:"color,omitempty" yaml:"color,omitempty"`

	//Payload, JSON format
	payload *bytes.Buffer `json:"-" yaml:"-"`
}

// Validate performs validation for MessageCard, checks for Text field
func (card *MessageCard) Validate() error {
	if card.Title == "" {
		return fmt.Errorf("invalid message card: title required")
	}
	if card.Text == "" {
		return fmt.Errorf("invalid message card: text required")
	}
	return nil
}

// Prepare handles the task to construct payload
func (card *MessageCard) Prepare() error {
	jsonMessage, err := json.Marshal(card)
	if err != nil {
		return fmt.Errorf("error marshalling MessageCard to JSON: %w", err)
	}
	if card.payload == nil {
		card.payload = &bytes.Buffer{}
	} else {
		card.payload.Reset()
	}

	_, err = card.payload.Write(jsonMessage)
	if err != nil {
		return fmt.Errorf("errorr writing JSON for MessageCard: %w", err)
	}

	return nil
}

// Payload returns the payload field, Prepare() should be called before this method
func (card *MessageCard) Payload() io.Reader {
	return card.payload
}

// CreateMessageCard creates a new MessageCard with required fields
func CreateMessageCard() *MessageCard {
	return &MessageCard{
		Type:    "MessageCard",
		Context: "https://schema.org/extensions",
	}
}
