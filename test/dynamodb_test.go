package main

import (
	"fmt"
	"time"
	"testing"
	"../dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type PlayerScore struct {
	PK      string     `json:"PK,omitempty"`
	Name    string     `json:"NAME,omitempty"`
	Score   int        `json:"SCORE,omitempty"`
	UserID  string     `json:"USER_ID,omitempty"`
	EndDate *time.Time `json:"END_DATE,omitempty"`
}

func main() {

}

func TestQuery(t *testing.T) {
	userID := "amzn1.ask.account.AHPJZMV4MAXNKAXO7MWZAGW7Y6HEOUV2UC6UWPOYNTUTRNQCJTKX7O6PME3ECF23PDIIAEZ7YQ2H4HGCH6B6QTZEONRPDHP3V7RXFMCWP2TP6FLDRXV6OR34TTQV4TL42AHUI5M4QSTA5YGXPERE33WRLGLIZI5Y45O6SEH237MALBKA2PPH7OS7IR6AEAN65UB2HNISOAEX6CA"

	values := map[string]string{":userID": userID}

	queryExp := "USER_ID = :userID"

	svc := dynamodb.GetDynamoInstance()
	results, err := dynamodb.Query(svc, values, "PLAYERSCORE", queryExp, "userId-index")

	if err != nil {
		t.Errorf("Error calling SelectItems:" + err.Error())
	}

	for _, i := range results {
    item := PlayerScore{}

    err = dynamodbattribute.UnmarshalMap(i, &item)

    if err != nil {
        fmt.Println("Got error unmarshalling:")
        fmt.Println(err.Error())
		}
		fmt.Println(item)
	}

}