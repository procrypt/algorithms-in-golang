package main

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

var (
	content string
	data    []string
	qURL    = "https://sqs.us-east-1.amazonaws.com/263331787521/procrypt"
	svcSqs  = sqs.New(session.New())
)

func sendsSQS(msg []string) {
	svcSqs.SendMessage(&sqs.SendMessageInput{
		DelaySeconds: aws.Int64(10),
		MessageAttributes: map[string]*sqs.MessageAttributeValue{
			"Location": &sqs.MessageAttributeValue{
				DataType:    aws.String("String"),
				StringValue: aws.String(data[0]),
			},
			"Time": &sqs.MessageAttributeValue{
				DataType:    aws.String("String"),
				StringValue: aws.String(data[1]),
			},
			"Phone": &sqs.MessageAttributeValue{
				DataType:    aws.String("String"),
				StringValue: aws.String(data[2]),
			},
		},
		MessageBody: aws.String("Information about reservation."),
		QueueUrl:    &qURL,
	})
	data = []string{}
}

func handleRequest(req events.LexEvent) (events.LexEvent, error) {
	// GreetingIntent
	if req.CurrentIntent.Name == "GreetingIntent" {
		dialogAction := events.LexEvent{
			DialogAction: &events.LexDialogAction{
				Type:             "Close",
				FulfillmentState: "Fulfilled",
				Message: map[string]string{
					"contentType": "PlainText",
					"content":     "Hi there,how can I help?",
				},
			},
		}
		json.Marshal(dialogAction)
		return dialogAction, nil

		// ThankYouIntent
	} else if req.CurrentIntent.Name == "ThankYouIntent" {
		dialogAction := events.LexEvent{
			DialogAction: &events.LexDialogAction{
				Type:             "Close",
				FulfillmentState: "Fulfilled",
				Message: map[string]string{
					"contentType": "PlainText",
					"content":     "You’re welcome.",
				},
			},
		}
		json.Marshal(dialogAction)
		return dialogAction, nil

		// DiningSuggestionsIntent
	} else if req.CurrentIntent.Name == "DiningSuggestionsIntent" {
		if req.InputTranscript == "I need some restaurant suggestions." {
			content = " Great! I can help you with that. What city or city area are you looking to dine in?"
			dialogAction := events.LexEvent{
				DialogAction: &events.LexDialogAction{
					Type: "ElicitSlot",
					Message: map[string]string{
						"contentType": "PlainText",
						"content":     content,
					},
					IntentName: "DiningSuggestionsIntent",
					Slots: map[string]string{
						"Time":     req.CurrentIntent.Slots["Time"],
						"phone":    req.CurrentIntent.Slots["phone"],
						"location": req.CurrentIntent.Slots["location"],
					},
					SlotToElicit: "location",
				},
			}
			json.Marshal(dialogAction)
			return dialogAction, nil

		} else if req.CurrentIntent.Slots["location"] != "null" && req.InputTranscript == req.CurrentIntent.Slots["location"] {
			content = "Got it " + req.CurrentIntent.Slots["location"] + "." + " What cuisine would you like to try?"
			data = append(data, req.CurrentIntent.Slots["location"])
		} else if req.InputTranscript == "Japanese" {
			content = "Ok, how many people are in your party?"
		} else if req.InputTranscript == "Two people" {
			content = "A few more to go. What date?"
		} else if req.InputTranscript == "Today" {
			content = "What time?"
			dialogAction := events.LexEvent{
				DialogAction: &events.LexDialogAction{
					Type: "ElicitSlot",
					Message: map[string]string{
						"contentType": "PlainText",
						"content":     content,
					},
					IntentName: "DiningSuggestionsIntent",
					Slots: map[string]string{
						"Time":     req.CurrentIntent.Slots["Time"],
						"phone":    req.CurrentIntent.Slots["phone"],
						"location": req.CurrentIntent.Slots["location"],
					},
					SlotToElicit: "Time",
				},
			}
			json.Marshal(dialogAction)
			return dialogAction, nil

		} else if req.CurrentIntent.Slots["Time"] != "null" && req.CurrentIntent.Slots["phone"] == "" {
			data = append(data, req.CurrentIntent.Slots["Time"])
			content = "Awesome! Lastly, I need your phone number so I can send you my findings."
			dialogAction := events.LexEvent{
				DialogAction: &events.LexDialogAction{
					Type: "ElicitSlot",
					Message: map[string]string{
						"contentType": "PlainText",
						"content":     content,
					},
					IntentName: "DiningSuggestionsIntent",
					Slots: map[string]string{
						"Time":     req.CurrentIntent.Slots["Time"],
						"phone":    req.CurrentIntent.Slots["phone"],
						"location": req.CurrentIntent.Slots["location"],
					},
					SlotToElicit: "phone",
				},
			}
			json.Marshal(dialogAction)
			return dialogAction, nil

		} else if req.CurrentIntent.Slots["phone"] != "null" {
			content = "You’re all set. Expect my recommendations shortly! Have a good day."
			data = append(data, req.CurrentIntent.Slots["phone"])
			sendsSQS(data)
		}
		dialogAction := events.LexEvent{
			DialogAction: &events.LexDialogAction{
				Type:             "Close",
				FulfillmentState: "Fulfilled",
				Message: map[string]string{
					"contentType": "PlainText",
					"content":     content,
				},
			},
		}
		json.Marshal(dialogAction)
		return dialogAction, nil
	} else {
		return events.LexEvent{}, nil
	}
}

func main() {
	lambda.Start(handleRequest)
}
