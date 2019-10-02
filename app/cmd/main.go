package main

import (
	"fmt"
	"log"
	"runtime"

	"user-service/app/config"
	"user-service/app/model"
	"user-service/app/platform/nats"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	envFlag    = "env"
	defaultEnv = "dev"
	serverPort = "server.port"
	natsServer = "nats.server"
)

func main() {

	// Get environment flag
	env := pflag.String(envFlag, defaultEnv, "environment config value to use")
	pflag.Parse()

	if err := config.LoadConfiguration(*env); err != nil {
		checkErr(err)
	}

	// Create new NATS server connection
	nc, err := natsclient.NewNATSServerConnection(viper.GetString(natsServer))
	checkErr(err)

	// Subscribe to user.create
	nc.Subscribe("user.create", func(u *model.User) {
		fmt.Printf("Received user.create: %+v\n", u)
		saveUser(u)
		nc.Publish("user.create.completed", u)
	})

	// Subscribe to user.list
	nc.Subscribe("user.list", func(ul *[]model.User) {
		fmt.Println("Received user.list")
		getUsers(ul)
		nc.Publish("user.list.completed", ul)
	})

	runtime.Goexit()
}

func saveUser(u *model.User) {

	config := &aws.Config{
		Region:   aws.String("us-east-1"),
		Endpoint: aws.String("http://localhost:8000"),
	}

	sess := session.Must(session.NewSession(config))

	svc := dynamodb.New(sess)

	// Auto increment id
	input := &dynamodb.UpdateItemInput{
		TableName: aws.String("NextIdTable"),
		Key: map[string]*dynamodb.AttributeValue{
			"NextKey": {
				S: aws.String("Users"),
			},
		},
		UpdateExpression: aws.String("ADD NextId :x"),

		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":x": {
				N: aws.String("1"),
			},
		},
		ReturnValues: aws.String("UPDATED_NEW"),
	}

	resp, err := svc.UpdateItem(input)
	var nextIdTable model.NextIdTable
	err = dynamodbattribute.UnmarshalMap(resp.Attributes, &nextIdTable)
	u.Id = nextIdTable.NextId
	userMap, err := dynamodbattribute.MarshalMap(u)

	// Create a new User
	params := &dynamodb.PutItemInput{
		TableName: aws.String("Users"),
		Item:      userMap,
	}
	fmt.Printf("\nparams: %+v", params)

	res, err := svc.PutItem(params)

	if err != nil {
		fmt.Printf("Error while save user: %v\n", err.Error())
		return
	}
	fmt.Println(res)
	fmt.Println("User saved successfully")

}

func getUsers(ul *[]model.User){

	fmt.Println("getUsers")
	config := &aws.Config{
		Region:   aws.String("us-east-1"),
		Endpoint: aws.String("http://localhost:8000"),
	}

	sess := session.Must(session.NewSession(config))

	svc := dynamodb.New(sess)

	params := &dynamodb.ScanInput{
		TableName: aws.String("Users"),
	}
	result, err := svc.Scan(params)

	fmt.Println(result)
	if err != nil {
		fmt.Errorf("failed to make Query API call, %v", err)

	}

	
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &ul)
	if err != nil {
		fmt.Errorf("failed to unmarshal Query result items, %v", err)

	}
	fmt.Printf("Users: %+v", ul)
	
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
