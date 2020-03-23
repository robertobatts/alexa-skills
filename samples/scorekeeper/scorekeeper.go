package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"time"

	"github.com/robertobatts/golexa/"
	"github.com/robertobatts/golexa/samples/scorekeeper/dynamodb"

	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

var gxa = &golexa.Golexa{Triggerable: &Scorekeeper{}, TranslationMap: GetjsonMap("./msg.json")}

//Scorekeeper is useful to override the Triggerable methods
type Scorekeeper struct {
}

//PlayerScore maps PLAYERSCORE table's items
type PlayerScore struct {
	PK      string     `json:"PK,omitempty"`
	Name    string     `json:"NAME,omitempty"`
	Score   int        `json:"SCORE,omitempty"`
	UserID  string     `json:"USER_ID,omitempty"`
	EndDate *time.Time `json:"END_DATE,omitempty"`
}

/*func (resp *golexa.AlexaResponse) SavePlayerNumbers(req golexa.AlexaRequest) {
	resp.SessionAttributes = map[string]string {
		"number": req.Request.Intent.Slots["number"].Value,
	}
}*/

//GetJsonMap transforms the json file into a
func GetJsonMap(path string) map[string]map[string]string {
	jsonMap := map[string]map[string]string{}
	raw, _ := ioutil.ReadFile(path)
	json.Unmarshal(raw, &jsonMap)
	return jsonMap
}

//SaveNewPlayer saves the player name on dynamodb, then ask again the user for other players
func SaveNewPlayer(req golexa.AlexaRequest, resp *golexa.AlexaResponse) {
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
func UpdatePlayerScore(req golexa.AlexaRequest, resp *golexa.AlexaResponse) {
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
		text = gxa.Translate("errorTryAgain")
	}
	resp.Ask(text)
}

//ReadScore reads the scores of all the players associated to the alexa's userID
func ReadScore(req golexa.AlexaRequest, resp *golexa.AlexaResponse) {
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
				text += item.Name + " " + gxa.Translate("has") + " " + strconv.Itoa(item.Score) + " " + gxa.Translate("points") + ", "
			}

		}
	}
	resp.Say(text)
}

//OnLaunch overrides Triggerable.OnLaunch
func (scorekeeper *Scorekeeper) OnLaunch(ctx context.Context, req golexa.AlexaRequest, resp *golexa.AlexaResponse) error {
	resp.Ask(gxa.Translate("playerNames"))

	return nil
}

//OnIntent overrides Tribgerable.OnIntent
func (scorekeeper *Scorekeeper) OnIntent(ctx context.Context, req golexa.AlexaRequest, resp *golexa.AlexaResponse) error {
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
		resp.Say(gxa.Translate("dontUnderstand"))
	}

	return nil
}

func main() {
	gxa.LambdaStart()
}
