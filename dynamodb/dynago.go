package dynago 

import (
	"fmt"
	"time"
	"os"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type PlayerScore struct {
	PK string 					`json:"PK"`
	Name string 				`json:"NAME"`
	Score int 					`json:"SCORE"`
	UserId string 			`json:"USER_ID"`
	EndDate *time.Time	`json:"END_DATE,omitempty"`
}

func CreateNewSession() *session.Session {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	return sess
}

func GetDynamoInstance() *dynamodb.DynamoDB {
	sess := CreateNewSession()
	fmt.Print((sess.Config.Credentials.Get()))
	return dynamodb.New(sess)
}


func main() {

	svc := GetDynamoInstance()
	fmt.Println(svc)

	playerScore := PlayerScore {
		PK: "TEST_1",
		Name: "Roberto",
		Score: 12,
		UserId: "TEST",
	}

	InsertItem(svc, playerScore)
}

func InsertItem(svc *dynamodb.DynamoDB, item PlayerScore) {
	av, err := dynamodbattribute.MarshalMap(item)

	if err == nil {
		input := &dynamodb.PutItemInput{
			Item:      av,
			TableName: aws.String("PLAYERSCORE"),
		}

		_, err = svc.PutItem(input)

		if err != nil {
			fmt.Println("Error calling PutItem:")
			fmt.Println(err.Error())
			os.Exit(1)
		}
	}
}