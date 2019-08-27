package golexa

import (
	"fmt"
	"context"
	"errors"
	"github.com/aws/aws-lambda-go/lambda"
)

//AlexaRequest defines the json passed from Alexa when the skill is called
type AlexaRequest struct {
	Version string 		`json:"version"`
	Session *Session	`json:"session"`
	Request *Request `json:"request"`
}

//AlexaResponse defines the json passed to Alexa 
type AlexaResponse struct {
	Version           string            `json:"version"`
	SessionAttributes map[string]string `json:"sessionAttributes,omitempty"`
	Response 					*Response 				`json:"response"`
}

//Session contains the session data of the request
type Session struct {
	New        bool   `json:"new"`
	SessionID  string `json:"sessionId"`
	Attributes struct {
		String map[string]interface{} `json:"string"`
	} `json:"attributes"`
	User struct {
		UserID      string `json:"userId"`
		AccessToken string `json:"accessToken"`
	} `json:"user"`
	Application struct {
		ApplicationID string `json:"applicationId"`
	} `json:"application"`
}

//Request contains main data of the AlexaRequest
type Request struct {
	Locale			string 	`json:"locale"`
	Type   			string 	`json:"type"`
	Time   			string 	`json:"timestamp"`
	RequestID   string 	`json:"requestId"`
	DialogState string 	`json:"dialogState"`
	Intent 			Intent `json:"intent"`
	Name 				string 	`json:"name"`
}

//Response contains main data of the AlexaResponse
type Response struct {
	OutputSpeech     *OutputSpeech `json:"outputSpeech,omitempty"`
	Reprompt         *Reprompt     `json:"reprompt,omitempty"`
	Card             *Card         `json:"card,omitempty"`
	Directives       []interface{} `json:"directives,omitempty"`
	ShouldEndSession bool          `json:"shouldEndSession,omitempty"`
}

//Intent contained the intent detected by Alexa
type Intent struct {
	Name               string                `json:"name"`
	ConfirmationStatus string                `json:"confirmationstatus"`
	Slots              map[string]IntentSlot `json:"slots"`
}

//OutputSpeech contains what Alexa says to the user
type OutputSpeech struct {
	Type string `json:"type,omitempty"`
	Text string `json:"text,omitempty"`
	SSML string `json:"ssml,omitempty"`
}

//Reprompt containes what Alexa asks to the user
type Reprompt struct {
	OutputSpeech OutputSpeech `json:"outputSpeech"`
}

//IntentSlot contains an Alexa Slot
type IntentSlot struct {
	Name               string       `json:"name"`
	ConfirmationStatus string       `json:"confirmationStatus,omitempty"`
	Value              string       `json:"value"`
	Resolutions        *Resolutions `json:"resolutions,omitempty"`
}

//Resolutions contain	extra properties of a slot
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

//Card contains data displayed by Alexa
type Card struct {
	Type    string `json:"type"`
	Title   string `json:"title,omitempty"`
	Content string `json:"content,omitempty"`
	Text    string `json:"text,omitempty"`
	Image   struct {
		SmallImageURL string `json:"smallImageUrl,omitempty"`
		LargeImageURL string `json:"largeImageUrl,omitempty"`
	} `json:"image,omitempty"`
}

//Triggerable contains the definition of the functions triggered by Handle()
type Triggerable interface {
	OnLaunch(ctx context.Context, req *AlexaRequest, resp *AlexaResponse) error
	OnIntent(ctx context.Context, req *AlexaRequest, resp *AlexaResponse) error
}

//Golexa is the container that implements the triggerable functions
type Golexa struct {
	Triggerable Triggerable
}

//Handle triggers the triggerable function by looking at the request
func (golexa *Golexa) Handle(ctx context.Context, req *AlexaRequest) (*AlexaResponse, error) {

	resp := CreateResponse()

	switch req.Request.Type {
		case "LaunchRequest":
			err := golexa.Triggerable.OnLaunch(ctx, req, resp)
			if err != nil {
				fmt.Println("Error on launch: " + err.Error())
			}
			return resp, err
		case "IntentRequest":
			err := golexa.Triggerable.OnIntent(ctx, req, resp)
			if err != nil {
				fmt.Println("Error on launch: " + err.Error())
			}
			return resp, err
		default:
			return resp, errors.New("Request type not recognized")
	}
}

//CreateResponse initialize the Alexa Response
func CreateResponse() *AlexaResponse {
	var resp AlexaResponse
	resp.Version = "2.0"
	resp.Response = &Response{}
	return &resp
}

//Say make Alexa say something to the user
func (resp *AlexaResponse) Say(text string) {
	resp.Response.OutputSpeech = &OutputSpeech{
		Type: "PlainText",
		Text: text,
	}
}

//Ask make Alexa ask something to the user
func (resp *AlexaResponse) Ask(text string) {
	resp.Response.Reprompt = &Reprompt{
		OutputSpeech: OutputSpeech{
			Type: "PlainText",
			Text: text,
		},
	}
	resp.Response.ShouldEndSession = false
}

//LambdaStart pass the function to be triggered to lambda
func LambdaStart(handlerFunc interface{}) {
	lambda.Start(handlerFunc)
}