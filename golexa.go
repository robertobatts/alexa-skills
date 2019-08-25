package golexa

import (
	"github.com/aws/aws-lambda-go/lambda"
)


type OutputSpeech struct {
	Type string `json:"type,omitempty"`
	Text string `json:"text,omitempty"`
}

type Reprompt struct {
	OutputSpeech OutputSpeech `json:"outputSpeech"`
}

type IntentSlot struct {
	Name               string       `json:"name"`
	ConfirmationStatus string       `json:"confirmationStatus,omitempty"`
	Value              string       `json:"value"`
	Resolutions        *Resolutions `json:"resolutions,omitempty"`
}

type Resolutions struct {
	ResolutionsPerAuthority []struct {
		Authority string `json:"authority"`
		Status    struct {
			Code string `json:"code"`
		} `json:"status"`
		Values []struct {
			Value struct {
				Name string `json:"name"`
				ID   string `json:"id"`
			} `json:"value"`
		} `json:"values"`
	} `json:"resolutionsPerAuthority"`
}

type Request struct {
	Version string `json:"version"`
	Session struct {
		User struct {
			UserID string `json:"userId"`
		} `json:"userId"`
	} `json:"session"`
	Request struct {
		Type   string `json:"type"`
		Time   string `json:"timestamp"`
		Intent struct {
			Name               string                `json:"name"`
			ConfirmationStatus string                `json:"confirmationstatus"`
			Slots              map[string]IntentSlot `json:"slots"`
		} `json:"intent"`
	} `json:"request"`
}

type Response struct {
	Version           string            `json:"version"`
	SessionAttributes map[string]string `json:"sessionAttributes,omitempty"`
	Response          struct {
		OutputSpeech     *OutputSpeech `json:"outputSpeech,omitempty"`
		Reprompt         *Reprompt     `json:"reprompt,omitempty"`
		ShouldEndSession bool          `json:"shouldEndSession,omitempty"`
	} `json:"response"`
}

func CreateResponse() *Response {
	var resp Response
	resp.Version = "1.0"
	return &resp
}

func (resp *Response) Say(text string) {
	resp.Response.OutputSpeech = &OutputSpeech{
		Type: "PlainText",
		Text: text,
	}
}

func (resp *Response) Ask(text string) {
	resp.Response.Reprompt = &Reprompt{
		OutputSpeech: OutputSpeech{
			Type: "PlainText",
			Text: text,
		},
	}
	resp.Response.ShouldEndSession = false
}

func LambdaStart(handlerFunc interface{}) {
	lambda.Start(handlerFunc)
}