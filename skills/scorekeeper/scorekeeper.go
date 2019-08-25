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
	"golexa"
)


type PlayerScore struct {
	PK      string     `json:"PK,omitempty"`
	Name    string     `json:"NAME,omitempty"`
	Score   int        `json:"SCORE,omitempty"`
	UserId  string     `json:"USER_ID,omitempty"`
	EndDate *time.Time `json:"END_DATE,omitempty"`
}


/*func (resp *golexa.Response) SavePlayerNumbers(req golexa.Request) {
	resp.SessionAttributes = map[string]string {
		"number": req.Request.Intent.Slots["number"].Value,
	}
}*/

func SaveNewPlayer(req golexa.Request, resp *golexa.Response) {
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

func UpdatePlayerScore(req golexa.Request, resp *golexa.Response) {
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

func ReadScore(req golexa.Request, resp *golexa.Response) {
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

func HandleRequest(ctx context.Context, req golexa.Request) (golexa.Response, error) {
	// Use Spew to output the request for debugging purposes:
	fmt.Println("---- Dumping Input Map: ----")
	spew.Dump(req)

	resp := golexa.CreateResponse()

	if req.Request.Type == "LaunchRequest" {
		resp.Ask("What are the players names?")

		return *resp, nil
	}

	switch req.Request.Intent.Name {
	case "addplayer":
		SaveNewPlayer(req, resp)
	case "addscore":
		UpdatePlayerScore(req, resp)
	case "readscore":
		ReadScore(req, resp)
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