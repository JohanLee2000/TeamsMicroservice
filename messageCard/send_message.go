// Created by Johan Lee - June 2023

/* This file implements the sending of a message as well as validation and error checking
for the URL, request processing and message construction. */

package messageCard

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

// Regex constant to validate the pattern of incoming webhook url provided by user
const (
	WebhookURLValidPattern = `^https:\/\/(?:.*\.webhook|outlook)\.office(?:365)?\.com`

	//Other URL regex constants can go here

)

// ExpectedEndpointResponseText is the expected success response text when submitting messages, given by webhook endpoint
const ExpectedEndpointResponseText string = "1"

// WebhookSendTimeout is the amount of time a message takes before it times out and is cancelled
const WebhookSendTimeout = 5 * time.Second

// ErrWebHookURLPattern returns when the URL does not match any specified pattern
var ErrWebHookURLPattern = errors.New("the webhook URL does not match the expected pattern")

// ErrInvalidResponseText returns when the message is unsuccessful via a response text
var ErrInvalidResponseText = errors.New("message unsuccessful, invalid webhook URL response text")

//Interface & Structs--------------------------------------------

// MessageSender functions as a client
type MessageSender interface {
	HTTPClient() *http.Client
	//UserAgent() string
	ValidateWebhook(webhookURL string) error
}

// messagePreparer prepares messages via marshaling
type messagePreparer interface {
	Prepare() error
}

// messageValidator provides validation for the format of a message
type messageValidator interface {
	Validate() error
}

// teamsMessage supports message formats to submit to a Microsoft Teams channel
type teamsMessage interface {
	messagePreparer
	messageValidator

	Payload() io.Reader
}

// TeamsClient submits messages to a Microsoft Teams channel
type TeamsClient struct {
	httpClient *http.Client
	//userAgent                    string
	//webhookURLValidationPatterns []string <- for multiple patterns
	//skipWebhookURLValidation bool
}

//Functions-------------------------------------------------------

// createTeamsClient creates a client to submit messages to the teams channel
func CreateTeamsClient() *TeamsClient {
	client := TeamsClient{
		httpClient: &http.Client{},
		//skipWebhookURLValidation: false,
	}
	return &client
}

// HTTPClient returns the http.Client field in the TeamsClient struct
func (client *TeamsClient) HTTPClient() *http.Client {
	return client.httpClient
}

// Unused for now

// setHTTPClient sets a new http.Client value to replace the old one
// func (client *TeamsClient) setHTTPClient(httpClient *http.Client) *TeamsClient {
// 	client.httpClient = httpClient
// 	return client
// }

// ValidateWebhook uses the constant WebhookURLValidPattern to ensure the URL is valid, can check for multiple patterns with patterns param
func (client *TeamsClient) ValidateWebhook(webhookURL string) error {
	urlLink, err := url.Parse(webhookURL)
	if err != nil {
		return fmt.Errorf("could not parse webhook URL %q: %w", webhookURL, err)
	}

	patterns := []string{WebhookURLValidPattern}
	//For loop here for multiple patterns
	for _, thisPattern := range patterns {
		match, err := regexp.MatchString(thisPattern, webhookURL)
		if err != nil {
			return err
		}
		if match {
			return nil
		}
	}

	return fmt.Errorf("%w; got: %q", ErrWebHookURLPattern, urlLink.String())
}

// prepareRequest prepares a http.Request with the required headers to send a message
func prepareRequest(ctxt context.Context, webhookURL string, message io.Reader) (*http.Request, error) {
	request, err := http.NewRequestWithContext(ctxt, http.MethodPost, webhookURL, message)
	if err != nil {
		return nil, err
	}
	request.Header.Add("Content-Type", "application/json;charset=utf-8")

	return request, nil
}

// processResponse validates the response from the endpoint after sending a message
func processResponse(response *http.Response) (string, error) {
	data, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	responseString := string(data)

	//Checks status code and endpoint response
	if response.StatusCode >= 299 {
		err := fmt.Errorf("error on code: %v, %q", response.Status, responseString)
		return "", err
	} else if responseString != strings.TrimSpace(ExpectedEndpointResponseText) {
		err := fmt.Errorf("got %q, expected %q: %w", responseString, ExpectedEndpointResponseText, ErrInvalidResponseText)
		return "", err
	} else {
		return responseString, nil
	}
}

// Send uses a function sendWithContext to send a message with the ability to timeout
func (client *TeamsClient) Send(webhookURL string, message teamsMessage) error {
	// For timeout
	ctxt, cancel := context.WithTimeout(context.Background(), WebhookSendTimeout)
	defer cancel()

	return sendWithContext(ctxt, client, webhookURL, message)
}

// sendWithContext sends a message to the Teams channel using the given webhookURL and client
func sendWithContext(ctxt context.Context, client MessageSender, webhookURL string, message teamsMessage) error {
	if err := client.ValidateWebhook(webhookURL); err != nil {
		return fmt.Errorf("webhook URL validation failed: %w", err)
	}

	if err := message.Validate(); err != nil {
		return fmt.Errorf("failed to validate message: %w", err)
	}

	if err := message.Prepare(); err != nil {
		return fmt.Errorf("failed to prepare message: %w", err)
	}

	// Prepare request with payload
	request, err := prepareRequest(ctxt, webhookURL, message.Payload())
	if err != nil {
		return fmt.Errorf("failed to prepare request: %w", err)
	}

	// Send message
	response, err := client.HTTPClient().Do(request)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	//Close response body
	defer func() {
		if err := response.Body.Close(); err != nil {
			log.Printf("Error closiing response body: %v", err)
		}
	}()

	//Process the response, check status code to ensure success
	responseString, err := processResponse(response)
	if err != nil {
		return fmt.Errorf("failed to process response: %w", err)
	}
	log.Printf("Response string: %v\n", responseString)

	return nil
}
