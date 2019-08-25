package golexa

import (
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

type Request struct {
	Locale			string 	`json:"locale"`
	Type   			string 	`json:"type"`
	Time   			string 	`json:"timestamp"`
	RequestID   string 	`json:"requestId"`
	DialogState string 	`json:"dialogState"`
	Intent 			Intent `json:"intent"`
	Name 				string 	`json:"name"`
}

type Response struct {
	OutputSpeech     *OutputSpeech `json:"outputSpeech,omitempty"`
	Reprompt         *Reprompt     `json:"reprompt,omitempty"`
	Card             *Card         `json:"card,omitempty"`
	Directives       []interface{} `json:"directives,omitempty"`
	ShouldEndSession bool          `json:"shouldEndSession,omitempty"`
}

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