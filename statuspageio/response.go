// Package statuspageio handles communicating with status pages following the statuspage.io format
package statuspageio

import (
	"encoding/json"
	"fmt"
	"time"

	goclient "github.com/sprak3000/go-client/client"
	"github.com/sprak3000/go-glitch/glitch"

	"github.com/sprak3000/go-whatsup-client/status"
)

// ServiceType is the name we use for various checks
const ServiceType = "statuspage.io"

// Page represents the page information
type Page struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	URL       string    `json:"url"`
	TimeZone  string    `json:"string"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Status represents the status information
type Status struct {
	Indicator   string `json:"indicator"`
	Description string `json:"description"`
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

// Response is the structure returned by statuspage.io powered service status pages
type Response struct {
	Page   Page   `json:"page"`
	Status Status `json:"status"`
}

// Indicator returns the current status for the service
func (r Response) Indicator() string {
	return r.Status.Indicator
}

// Name returns the name of the service
func (r Response) Name() string {
	return r.Page.Name
}

// UpdatedAt returns the date / time of the most recent status update for the service
func (r Response) UpdatedAt() time.Time {
	return r.Page.UpdatedAt
}

// URL returns the page the service hosts to display more details about current and past status updates
func (r Response) URL() string {
	return r.Page.URL
}
