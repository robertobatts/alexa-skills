package main

import (
	"context"
	"dynago"
	"fmt"

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

type AlexaRequest struct {
	Version string `json:"version"`
	Session struct {
		User struct {
			UserId string `json:"userId"`
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

type AlexaResponse struct {
	Version           string            `json:"version"`
	SessionAttributes map[string]string `json:"sessionAttributes,omitempty"`
	Response          struct {
		OutputSpeech     *OutputSpeech `json:"outputSpeech,omitempty"`
		Reprompt         *Reprompt     `json:"reprompt,omitempty"`
		ShouldEndSession bool          `json:"shouldEndSession,omitempty"`
	} `json:"response"`
}

func CreateResponse() *AlexaResponse {
	var resp AlexaResponse
	resp.Version = "1.0"
	return &resp
}

func (resp *AlexaResponse) Say(text string) {
	resp.Response.OutputSpeech = &OutputSpeech{
		Type: "PlainText",
		Text: text,
	}
}

func (resp *AlexaResponse) Ask(text string) {
	resp.Response.Reprompt = &Reprompt{
		OutputSpeech: OutputSpeech{
			Type: "PlainText",
			Text: text,
		},
	}
	resp.Response.ShouldEndSession = false
}

/*func (resp *AlexaResponse) SavePlayerNumbers(req AlexaRequest) {
	resp.SessionAttributes = map[string]string {
		"number": req.Request.Intent.Slots["number"].Value,
	}
}*/

func (resp *AlexaResponse) SaveNewPlayer(req AlexaRequest) {
	name := req.Request.Intent.Slots["name"].Value
	userId := req.Session.User.UserId

	playerScore := dynago.PlayerScore{
		PK:     userId + "_" + name,
		Name:   name,
		UserId: userId,
	}

	svc := dynago.GetDynamoInstance()
	err := dynago.InsertItem(svc, playerScore, "PLAYERSCORE")

	text := ""
	if err != nil {
		text = "There has been an error! Try again"
	}
	resp.Response.Reprompt = &Reprompt{
		OutputSpeech: OutputSpeech{
			Type: "PlainText",
			Text: text,
		},
	}
	resp.Response.ShouldEndSession = false
}

func (resp *AlexaResponse) UpdatePlayerScore(req AlexaRequest) {
	score := req.Request.Intent.Slots["score"].Value
	name := req.Request.Intent.Slots["name"].Value
	userId := req.Session.User.UserId

	playerScore := dynago.PlayerScore{
		PK:     userId + "_" + name,
		Name:   name,
		Score:  score,
		UserId: userId,
	}

	svc := dynago.GetDynamoInstance()
	err := dynago.UpdateItem(svc, playerScore, "PLAYERSCORE")

	text := ""
	if err != nil {
		text = "There has been an error! Try again"
	}
	resp.Response.Reprompt = &Reprompt{
		OutputSpeech: OutputSpeech{
			Type: "PlainText",
			Text: text,
		},
	}
	resp.Response.ShouldEndSession = false
}

func HandleRequest(ctx context.Context, req AlexaRequest) (AlexaResponse, error) {
	// Use Spew to output the request for debugging purposes:
	fmt.Println("---- Dumping Input Map: ----")
	spew.Dump(req)

	resp := CreateResponse()

	if req.Request.Type == "LaunchRequest" {
		resp.Ask("What are the players names?")

		return *resp, nil
	}

	switch req.Request.Intent.Name {
	/*case "howmanyplayers":
	resp.SavePlayerNumbers(req);
	resp.Ask("What are the players names?")*/
	case "playername":
		resp.SaveNewPlayer(req)
	case "playerscore":
		resp.UpdatePlayerScore(req)
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
