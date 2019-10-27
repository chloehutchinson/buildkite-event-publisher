package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type backend interface {
	Publish(eventName string, data map[string]interface{}) error
	Close() error
}

var (
	// NewRelicAppName is the name of the application in New Relic.
	NewRelicAppName string
	// NewRelicLicenseKey is the license key used to publish events to New Relic.
	NewRelicLicenseKey string
	// BuildkiteWebhookToken is the token used to identify webhook requests originating from Buildkite
	BuildkiteWebhookToken string
	publisher             backend
)

func main() {
	getEnvironment()
	lambda.Start(publisherWrapper)
}

func getEnvironment() {
	var ok bool

	NewRelicAppName, ok = os.LookupEnv("NEWRELIC_APP_NAME")
	if !ok || NewRelicAppName == "" {
		NewRelicAppName = "Buildkite"
	}

	NewRelicLicenseKey, ok = os.LookupEnv("NEWRELIC_LICENSE_KEY")
	if !ok || NewRelicLicenseKey == "" {
		log.Fatalln("NEWRELIC_LICENSE_KEY environment var missing!")
	}

	BuildkiteWebhookToken, ok = os.LookupEnv("BUILDKITE_WEBHOOK_TOKEN")
	if !ok || BuildkiteWebhookToken == "" {
		log.Fatalln("BUILDKITE_WEBHOOK_TOKEN environment var missing!")
	}
}

func handleBuildkiteEvent(eventType string, eventBody buildkiteEvent) error {
	var data map[string]interface{}
	log.Printf("Handling %s event...\n", eventType)
	var eventName string

	switch eventType {
	case "build.running", "build.finished":
		eventName = "BuildEvent"
		data = handleBuildEvent(eventBody)
	case "job.started", "job.finished", "job.activated":
		eventName = "JobEvent"
		data = handleJobEvent(eventBody)
	default:
		log.Printf("Unrecognised event type: %s, event body: %+v\n", eventType, eventBody)
		return errors.New(fmt.Sprintf("Unrecognised event type: %s", eventType))
	}

	log.Printf("Publishing %s event to New Relic\n", eventName)

	err := publisher.Publish(eventName, data)
	return err
}

func parseBuildkiteEvent(jsonBody string) (buildkiteEvent, error) {
	var event buildkiteEvent
	err := json.Unmarshal([]byte(jsonBody), &event)

	return event, err
}

func getPublisher() backend {
	var err error

	log.Printf("Connecting to New Relic backend...\n")
	publisher, err = NewNewRelicBackend(NewRelicAppName, NewRelicLicenseKey)
	if err != nil {
		log.Fatalf("ERROR: Unable to start New Relic backend: %+v\n", err)
	}

	return publisher
}

func publisherWrapper(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	publisher = getPublisher()
	defer publisher.Close()
	return handleProxyRequest(request)
}

func handleProxyRequest(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if request.Headers["X-Buildkite-Token"] != BuildkiteWebhookToken {
		return events.APIGatewayProxyResponse{
			StatusCode: 401, // unauthorised
			Body:       "ERROR: invalid token",
		}, nil
	}
	eventType := request.Headers["X-Buildkite-Event"]

	if request.Body == "" {
		return events.APIGatewayProxyResponse{
			StatusCode: 400, // bad request
			Body:       "ERROR: empty event body",
		}, errors.New("ERROR: fired with no request body")
	}

	eventBody, err := parseBuildkiteEvent(request.Body)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 400, // bad request
			Body:       "ERROR: Malformed event data",
		}, err
	}

	err = handleBuildkiteEvent(eventType, eventBody)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500, // internal server error
			Body:       "ERROR: Something went wrong handling event",
		}, err
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200, // ok
		Body:       "Hello Buildkite!",
	}, nil
}
