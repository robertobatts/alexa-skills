package main

import (
	"context"
	"fmt"
	"time"
	"strconv"

	"github.com/davecgh/go-spew/spew"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"dynamodb"
	"golexa"
)


//PlayerScore maps PLAYERSCORE table's items
type PlayerScore struct {
	PK      string     `json:"PK,omitempty"`
	Name    string     `json:"NAME,omitempty"`
	Score   int        `json:"SCORE,omitempty"`
	UserID  string     `json:"USER_ID,omitempty"`
	EndDate *time.Time `json:"END_DATE,omitempty"`
}


/*func (resp *golexa.Response) SavePlayerNumbers(req golexa.Request) {
	resp.SessionAttributes = map[string]string {
		"number": req.Request.Intent.Slots["number"].Value,
	}
}*/

//SaveNewPlayer saves the player name on dynamodb, then ask again the user for other players
func SaveNewPlayer(req golexa.Request, resp *golexa.Response) {
	name := req.Request.Intent.Slots["name"].Value
	userID := req.Session.User.UserID

	playerScore := PlayerScore{
		PK:     userID + "_" + name,
		Name:   name,
		UserID: userID,
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

//UpdatePlayerScore updates the player's score, then wait for the user to take other scores
func UpdatePlayerScore(req golexa.Request, resp *golexa.Response) {
	score, err := strconv.Atoi(req.Request.Intent.Slots["score"].Value)
	name := req.Request.Intent.Slots["name"].Value
	userID := req.Session.User.UserID

	keys := PlayerScore{PK: userID + "_" + name}
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

//ReadScore reads the scores of all the players associated to the alexa's userID
func ReadScore(req golexa.Request, resp *golexa.Response) {
	userID := req.Session.User.UserID
	values := map[string]string{":userID": userID}

	queryExp := "USER_ID = :userID"

	svc := dynamodb.GetDynamoInstance()
	results, err := dynamodb.Query(svc, values, "PLAYERSCORE", queryExp, "userID-index")

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

//HandlerRequest is the handler function of lambda
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
	golexa.LambdaStart(HandleRequest)
}