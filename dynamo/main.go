package main

import (
	_ "context"
	"fmt"
	"math/rand"
	"strconv"
	_ "strconv"
	_ "strings"
	"time"

	"github.com/arienmalec/alexa-go"

	_ "math/rand"

	_ "github.com/arienmalec/alexa-go"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
)

type Alexa struct {
	Question string            `dynamo:"question" json:"question"`
	Options  map[string]string `dynamo:"set" json:"options"`
	Id       int               `dynamo:"omitempty" json:"id"`
	Correct  []string          `dynamo:"set" json:"correct"`
	Count    int               `dynamo:"omitempty" json:"count"`
}

//creating connection to the database
func id() int {
	var dbItem Alexa
	db := dynamo.New(session.New(), &aws.Config{Region: aws.String("us-east-1")})
	table := db.Table("alexadeveloperquiz")
	err := table.Get("id", 0).One(dynamo.AWSEncoding(&dbItem))

	if err != nil {
		fmt.Println(err)
	}
	//get count from the table
	Y := dbItem.Count
	//randomize a number between 1 and Count
	rand.Seed(time.Now().UnixNano())
	min := 1
	max := Y
	randNum := rand.Intn(max-min) + min
	return randNum
}

func question() (alexa.Response, error) {
	var dbItem2 Alexa

	db := dynamo.New(session.New(), &aws.Config{Region: aws.String("us-east-1")})
	table := db.Table("alexadeveloperquiz")
	err := table.Get("id", id()).One(dynamo.AWSEncoding(&dbItem2))

	if err != nil {
		fmt.Println(err)
	}
	question := dbItem2.Question
	options := dbItem2.Options
	fmt.Println(question)

	var optionsTogether string
	for opt, opt2 := range options {
		optionsTogether += fmt.Sprintf(opt + " - " + opt2 + "\n")
		fmt.Print(optionsTogether)
	}

	return alexa.NewSimpleResponse("question", question+optionsTogether), nil
}

func answer(request alexa.Request) alexa.Response {
	var dbItem3 Alexa

	db := dynamo.New(session.New(), &aws.Config{Region: aws.String("us-east-1")})
	table := db.Table("alexadeveloperquiz")
	err := table.Get("id", id()).One(dynamo.AWSEncoding(&dbItem3))

	if err != nil {
		fmt.Println(err)
	}
	answers := dbItem3.Correct

	var answerString string
	for ans := range answers {
		answerString += fmt.Sprintf(strconv.Itoa(ans) + " ")
	}
	return alexa.NewSimpleResponse("correct", answerString)
}

//func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
//	question()
//	answer()
//	return events.APIGatewayProxyResponse{
//		StatusCode: 200}, nil
//}

func main() {
	lambda.Start(handler)
}
