package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"fmt"
	"log"
	"os"
)

type Item struct {
	MemberId    string `dynamodbav:"MemberId"`
	Resolution  string `dynamodbav:"Resolution"`
	MemberState string `dynamodbav:"MemberState"`
}

var tableName string = os.Getenv("TABLE_NAME")

func initClient() *dynamodb.DynamoDB {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create DynamoDB client
	client := dynamodb.New(sess)

	return client
}

func getMemberData(memberId string) (string, string) {

	client := initClient()

	result, err := client.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"MemberId": {
				S: aws.String(memberId),
			},
		},
	})
	if err != nil {
		log.Fatalf("Got error calling GetItem: %s", err)
	}
	if result.Item == nil {
		log.Fatal("Got error result is nil")
	}

	log.Print(result.Item)

	item := Item{}

	err = dynamodbattribute.UnmarshalMap(result.Item, &item)
	if err != nil {
		panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
	}
	log.Print(item)
	log.Print(item.Resolution, item.MemberState)

	return item.Resolution, item.MemberState
}

func updateCounter(memberId string, postState string) {
	log.Print("Start Updating")
	client := initClient()

	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":r": {
				S: aws.String(postState),
			},
		},
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"MemberId": {
				S: aws.String(memberId),
			},
		},
		ReturnValues:     aws.String("UPDATED_NEW"),
		UpdateExpression: aws.String("set MemberState = :r"),
	}

	_, err := client.UpdateItem(input)
	if err != nil {
		log.Fatalf("Got error calling UpdateItem: %s", err)
	}
	log.Print("Finish Updating")
}
