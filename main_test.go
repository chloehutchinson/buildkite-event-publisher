package main

import (
	"errors"
	"os"
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

func TestGetEnvironmentSetsConstants(t *testing.T) {
	os.Setenv("BUILDKITE_WEBHOOK_TOKEN", "test1")
	os.Setenv("NEWRELIC_APP_NAME", "anything")
	os.Setenv("NEWRELIC_LICENSE_KEY", "test2")
	getEnvironment()

	if BuildkiteWebhookToken != "test1" {
		t.Fatalf("BUILDKITE_WEBHOOK_TOKEN should be set")
	}
	if NewRelicAppName != "anything" {
		t.Fatalf("NEWRELIC_APP_NAME should be set")
	}
	if NewRelicLicenseKey != "test2" {
		t.Fatalf("NEWRELIC_LICENSE_KEY should be set")
	}
}

func TestGetEnvironmentDefaultAppName(t *testing.T) {
	os.Setenv("BUILDKITE_WEBHOOK_TOKEN", "test1") // required value
	os.Setenv("NEWRELIC_LICENSE_KEY", "test2")    // required value
	os.Unsetenv("NEWRELIC_APP_NAME")              // ensure unset
	getEnvironment()

	if NewRelicAppName != "Buildkite" {
		t.Fatalf("NewRelicAppName should default to Buildkite (got %s)\n", NewRelicAppName)
	}
}

type mockBackend struct {
	err *error
	t   *testing.T
}

func (nr *mockBackend) Publish(eventName string, data map[string]interface{}) error {
	if nr.err != nil {
		return *nr.err
	}

	nr.t.Logf("Mock publish %s : %+v\n", eventName, data)
	return nil
}
func (nr *mockBackend) Close() error {
	nr.t.Logf("Close mockBackend\n")
	return nil
}

func mockPublisher(t *testing.T, err *error) {
	publisher = &mockBackend{
		t:   t,
		err: err,
	}
}

func TestHandleProxyRequestCheckTokenHeader(t *testing.T) {
	BuildkiteWebhookToken = "test"
	mockPublisher(t, nil)

	resp, _ := handleProxyRequest(events.APIGatewayProxyRequest{
		Headers: map[string]string{
			"X-Buildkite-Token": "fake",
		},
	})

	if resp.StatusCode != 401 {
		t.Fatalf("Expected handleProxyRequest to return 401, got %d\n", resp.StatusCode)
	}
}

func TestHandleProxyRequestReturn400OnEmptyBody(t *testing.T) {
	BuildkiteWebhookToken = "test"
	mockPublisher(t, nil)

	resp, _ := handleProxyRequest(events.APIGatewayProxyRequest{
		Headers: map[string]string{
			"X-Buildkite-Token": BuildkiteWebhookToken,
		},
		Body: "",
	})

	if resp.StatusCode != 400 {
		t.Fatalf("Expected handleProxyRequest to return 400, got %d\n", resp.StatusCode)
	}
}

func TestHandleProxyRequestReturn400OnInvalidJsonBody(t *testing.T) {
	BuildkiteWebhookToken = "test"
	mockPublisher(t, nil)

	resp, _ := handleProxyRequest(events.APIGatewayProxyRequest{
		Headers: map[string]string{
			"X-Buildkite-Token": BuildkiteWebhookToken,
		},
		Body: "never gonna give you up", // invalid json
	})

	if resp.StatusCode != 400 {
		t.Fatalf("Expected handleProxyRequest to return 400, got %d\n", resp.StatusCode)
	}
}

func TestHandleProxyRequestReturn500OnPublishError(t *testing.T) {
	BuildkiteWebhookToken = "test"
	fakeError := errors.New("fake error")
	mockPublisher(t, &fakeError)

	resp, _ := handleProxyRequest(events.APIGatewayProxyRequest{
		Headers: map[string]string{
			"X-Buildkite-Token": BuildkiteWebhookToken,
			"X-Buildkite-Event": "build.running",
		},
		Body: "{}", // valid json
	})

	if resp.StatusCode != 500 {
		t.Fatalf("Expected handleProxyRequest to return 500, got %d\n", resp.StatusCode)
	}
}

func TestHandleProxyRequestReturn200OnHappyPath(t *testing.T) {
	BuildkiteWebhookToken = "test"
	mockPublisher(t, nil)

	resp, err := handleProxyRequest(events.APIGatewayProxyRequest{
		Headers: map[string]string{
			"X-Buildkite-Token": BuildkiteWebhookToken,
			"X-Buildkite-Event": "build.running",
		},
		Body: "{}", // valid json
	})

	if resp.StatusCode != 200 {
		t.Fatalf("Expected handleProxyRequest to return 200, got %d (%s)\n", resp.StatusCode, err)
	}
}
