package main

import (
	"context"
	"fmt"
	"time"
	"strconv"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/davecgh/go-spew/spew"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"dynamodb"
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

type PlayerScore struct {
	PK      string     `json:"PK,omitempty"`
	Name    string     `json:"NAME,omitempty"`
	Score   int        `json:"SCORE,omitempty"`
	UserId  string     `json:"USER_ID,omitempty"`
	EndDate *time.Time `json:"END_DATE,omitempty"`
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

	playerScore := PlayerScore{
		PK:     userId + "_" + name,
		Name:   name,
		UserId: userId,
	}

	svc := dynamodb.GetDynamoInstance()
	err := dynamodb.InsertItem(svc, playerScore, "PLAYERSCORE")

	text := ""
	if err != nil {
		fmt.Println(err.Error())
		text = "There has been an error! Try again"
	}
	resp.Ask(text)
}

func (resp *AlexaResponse) UpdatePlayerScore(req AlexaRequest) {
	score, err := strconv.Atoi(req.Request.Intent.Slots["score"].Value)
	name := req.Request.Intent.Slots["name"].Value
	userId := req.Session.User.UserId

	keys := PlayerScore{PK: userId + "_" + name}
	values := map[string]int{":score": score}

	svc := dynamodb.GetDynamoInstance()
	_, err = dynamodb.UpdateItem(svc, values, keys, "PLAYERSCORE", "set SCORE = SCORE + :score")

	text := ""
	if err != nil {
		fmt.Println(err.Error())
		text = "There has been an error! Try again"
	}
	resp.Ask(text)
}

func (resp *AlexaResponse) ReadScore(req AlexaRequest) {
	userId := req.Session.User.UserId
	values := map[string]string{":userId": userId}

	queryExp := "USER_ID = :userId"

	svc := dynamodb.GetDynamoInstance()
	results, err := dynamodb.Query(svc, values, "PLAYERSCORE", queryExp, "userId-index")

	text := ""
	if err != nil {
		fmt.Println(err.Error())
		text = "There has been an error! Try again"
	} else {
		for _, i := range results {
			item := PlayerScore{}

			err = dynamodbattribute.UnmarshalMap(i, &item)

			if err != nil {
					fmt.Println("Got error unmarshalling:")
					fmt.Println(err.Error())
			} else {
				text += item.Name + " has " + strconv.Itoa(item.Score) + " points, "
			}	
			
		}
	}
	resp.Say(text)
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
	case "addplayer":
		resp.SaveNewPlayer(req)
	case "addscore":
		resp.UpdatePlayerScore(req)
	case "readscore":
		resp.ReadScore(req)
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