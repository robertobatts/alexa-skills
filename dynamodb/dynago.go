package dynago

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type PlayerScore struct {
	PK      string     `json:"PK"`
	Name    string     `json:"NAME"`
	Score   string     `json:"SCORE"`
	UserId  string     `json:"USER_ID"`
	EndDate *time.Time `json:"END_DATE,omitempty"`
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

	playerScore := PlayerScore{
		PK:     "TEST_1",
		Name:   "Roberto",
		Score:  "12",
		UserId: "TEST",
	}

	InsertItem(svc, playerScore, "PLAYERSCORE")
}

func InsertItem(svc *dynamodb.DynamoDB, item interface{}, tableName string) error {
	av, err := dynamodbattribute.MarshalMap(item)

	if err == nil {
		input := &dynamodb.PutItemInput{
			Item:      av,
			TableName: aws.String(tableName),
		}

		_, err = svc.PutItem(input)
		if err != nil {
			fmt.Println("Error calling PutItem:")
			fmt.Println(err.Error())
			return err
		}
	}
	return err
}

func UpdateItem(svc *dynamodb.DynamoDB, playerScore PlayerScore, tableName string) error {
	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":score": {
				N: aws.String(playerScore.Score),
			},
		},
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"PK": {
				S: aws.String(playerScore.PK),
			},
		},
		ReturnValues:     aws.String("UPDATED_NEW"),
		UpdateExpression: aws.String("set SCORE = SCORE + :score"),
	}

	_, err := svc.UpdateItem(input)
	if err != nil {
		fmt.Println("Error calling UpdateItem:")
		fmt.Println(err.Error())
	}
	return err
}
