// Package slack handles communicating with status pages following the Slack format
package slack

import (
	"encoding/json"
	"fmt"
	"time"

	goclient "github.com/sprak3000/go-client/client"
	"github.com/sprak3000/go-glitch/glitch"

	"github.com/sprak3000/go-whatsup-client/status"
)

// ServiceType is the name we use for various checks
const ServiceType = "slack"

// Note represents additional details about an incident
type Note struct {
	DateCreated time.Time `json:"date_created"`
	Body        string    `json:"body"`
}

// Incident represents a single reported status incident
type Incident struct {
	ID          int       `json:"id"`
	DateCreated time.Time `json:"date_created"`
	DateUpdated time.Time `json:"date_updated"`
	Title       string    `json:"title"`
	Type        string    `json:"type"`
	Status      string    `json:"status"`
	URL         string    `json:"url"`
	Services    []string  `json:"services"`
	Notes       []Note    `json:"notes"`
}

// ReadStatus handles communicating with the service to get its status details
func ReadStatus(c goclient.BaseClient, serviceName, slug string) (status.Details, glitch.DataError) {
	resp := Response{}

	respBytes, err := status.Get(c, slug)
	if err != nil {
		return resp, glitch.NewDataError(err, status.ErrorUnableToMakeClientRequest, fmt.Sprintf("unable to make client request for %s: %v", serviceName, err))
	}

	uErr := json.Unmarshal(respBytes, &resp)
	if uErr != nil {
		return resp, glitch.NewDataError(uErr, status.ErrorUnableToParseClientResponse, fmt.Sprintf("unable to parse client response for %s: %v", serviceName, uErr))
	}

	return resp, nil
}

// Response is the structure returned by Slack powered service status pages
type Response struct {
	Status          string     `json:"status"`
	DateCreated     time.Time  `json:"date_created"`
	DateUpdated     time.Time  `json:"date_updated"`
	ActiveIncidents []Incident `json:"active_incidents"`
}

// Indicator returns the current status for the service
func (r Response) Indicator() string {
	if r.Status == "active" {
		return "major"
	}

	return r.Status
}

// Name returns the name of the service
func (r Response) Name() string {
	return "Slack"
}

// UpdatedAt returns the date / time of the most recent status update for the service
func (r Response) UpdatedAt() time.Time {
	return r.DateUpdated
}

// URL returns the page the service hosts to display more details about current and past status updates
func (r Response) URL() string {
	return "https://status.slack.com/"
}
