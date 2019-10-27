package main

import (
	"log"
	"time"

	newrelic "github.com/newrelic/go-agent"
)

const newRelicConnectionTimeout = time.Second * 30

// NewRelicBackend sends metrics to New Relic Insights
type NewRelicBackend struct {
	client newrelic.Application
}

// NewNewRelicBackend returns a backend for New Relic
// Where appName is your desired application name in New Relic
//   and licenseKey is your New Relic license key
func NewNewRelicBackend(appName string, licenseKey string) (*NewRelicBackend, error) {
	config := newrelic.NewConfig(appName, licenseKey)
	app, err := newrelic.NewApplication(config)
	if err != nil {
		return nil, err
	}

	// Waiting for connection is essential or no data will make it during short-lived execution (e.g. Lambda)
	err = app.WaitForConnection(newRelicConnectionTimeout)
	if err != nil {
		return nil, err
	}

	return &NewRelicBackend{
		client: app,
	}, nil
}

// Publish metrics to New Relic
func (nr *NewRelicBackend) Publish(eventName string, data map[string]interface{}) error {
	err := nr.client.RecordCustomEvent(eventName, data)

	return err
}

// Close by shutting down NR client
func (nr *NewRelicBackend) Close() error {
	nr.client.Shutdown(newRelicConnectionTimeout)
	log.Printf("Disposed New Relic client")

	return nil
}
