package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/davecgh/go-spew/spew"
)

type OutputSpeech struct {
	Type string `json:"type,omitempty"`
	Text string `json:"text,omitempty"`
}

type Reprompt struct {
	OutputSpeech OutputSpeech `json:"outputSpeech"`
}

type AlexaRequest struct {
	Version string `json:"version"`
	Request struct {
		Type   string `json:"type"`
		Time   string `json:"timestamp"`
		Intent struct {
			Name               string `json:"name"`
			ConfirmationStatus string `json:"confirmationstatus"`
		} `json:"intent"`
	} `json:"request"`
}

type AlexaResponse struct {
	Version  string `json:"version"`
	Response struct {
		OutputSpeech *OutputSpeech `json:"outputSpeech,omitempty"`
		Reprompt *Reprompt `json:"reprompt,omitempty"`
		ShouldEndSession bool `json:"shouldEndSession,omitempty"`
	} `json:"response"`
}

func CreateResponse() *AlexaResponse {
	var resp AlexaResponse
	resp.Version = "1.0"
	return &resp
}

func (resp *AlexaResponse) Say(text string) {
	resp.Response.OutputSpeech = &OutputSpeech {
		Type: "PlainText",
		Text: text,
	}
}

func (resp *AlexaResponse) Ask(text string) {
	resp.Response.Reprompt = &Reprompt {
		OutputSpeech: OutputSpeech{
			Type: "PlainText",
			Text: text,
		},
	}
	resp.Response.ShouldEndSession = false
}

func HandleRequest(ctx context.Context, i AlexaRequest) (AlexaResponse, error) {
	// Use Spew to output the request for debugging purposes:
	fmt.Println("---- Dumping Input Map: ----")
	spew.Dump(i)
	fmt.Println("---- Done. ----")

	// Example of accessing map value via index:
	log.Printf("Request type is ", i.Request.Intent.Name)

	// Create a response object
	resp := CreateResponse()

	// Customize the response for each Alexa Intent
	switch i.Request.Intent.Name {
	case "howmanyplayers":
		resp.Ask("What are the players names?")
	case "AMAZON.HelpIntent":
		resp.Say("")
		//TODO
	default:
		resp.Say("I'm sorry, the input does not look like something I understand.")
	}

	return *resp, nil
}

func main() {
	lambda.Start(HandleRequest)
}