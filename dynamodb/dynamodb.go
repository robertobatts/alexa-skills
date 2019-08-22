package dynamodb

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

func CreateNewSession() *session.Session {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	return sess
}

func GetDynamoInstance() *dynamodb.DynamoDB {
	return dynamodb.New(CreateNewSession())
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

func UpdateItem(svc *dynamodb.DynamoDB, values interface{}, keys interface{}, tableName string, 
	updateExp string) (*dynamodb.UpdateItemOutput, error) {
	marshalledValues, _ := dynamodbattribute.MarshalMap(values)
	marshalledKeys, _ := dynamodbattribute.MarshalMap(keys)

	input := &dynamodb.UpdateItemInput {
		ExpressionAttributeValues: 	marshalledValues,
		TableName:                 	aws.String(tableName),
		Key:                       	marshalledKeys,
		ReturnValues:              	aws.String("UPDATED_NEW"),
		UpdateExpression:          	aws.String(updateExp),
	}

	newItem, err := svc.UpdateItem(input)

	if err != nil {
		fmt.Println("Error calling UpdateItem:")
		fmt.Println(err.Error())
	}
	return newItem, err
}

func SelectItems(svc *dynamodb.DynamoDB, values interface{}, tableName string, 
	queryExp string, indexName string) ([]map[string]*dynamodb.AttributeValue, error) {
	marshalledValues, _ := dynamodbattribute.MarshalMap(values)

	params := &dynamodb.QueryInput {
		ExpressionAttributeValues:	marshalledValues,
		KeyConditionExpression:			aws.String(queryExp),
		TableName: 									aws.String(tableName),
		IndexName:									aws.String(indexName),
	}

	result, err := svc.Query(params)

	if err != nil {
		fmt.Println("Error calling GetItem:")
		fmt.Println(err.Error())
		return nil, err
	}


	return result.Items, err
}