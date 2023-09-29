// Package whatsup is the client interface for getting status pages
package whatsup

import (
	"net/url"
	"time"

	"github.com/sprak3000/go-client/client"
	"github.com/sprak3000/go-glitch/glitch"

	"github.com/sprak3000/go-whatsup-client/slack"
	"github.com/sprak3000/go-whatsup-client/status"
	"github.com/sprak3000/go-whatsup-client/statuspageio"
)

//go:generate go run -mod=mod github.com/golang/mock/mockgen -package clientmock -destination=./clientmock/client-mock.go -source=../whatsup/client.go -build_flags=-mod=mod

// StatusPageClient handles the requests for status pages
type StatusPageClient interface {
	StatuspageIoService(serviceName, pageURL string) (status.Details, glitch.DataError)
	Slack() (status.Details, glitch.DataError)
}

type statusPageClient struct {
}

// NewStatusPageClient creates a new StatusPageClient
func NewStatusPageClient() StatusPageClient {
	return &statusPageClient{}
}

// StatuspageIoService handles fetching statuspage.io style status pages
func (spc *statusPageClient) StatuspageIoService(serviceName, pageURL string) (status.Details, glitch.DataError) {
	u, err := url.Parse(pageURL)
	sf := func(serviceName string, useTLS bool) (url.URL, error) {
		return *u, err
	}
	c := client.NewBaseClient(sf, serviceName, true, 10*time.Second, nil)
	return statuspageio.ReadStatus(c, serviceName, u.Path)
}

// Slack handles fetching the Slack status page
func (spc *statusPageClient) Slack() (status.Details, glitch.DataError) {
	sn := slack.ServiceType
	u, err := url.Parse("https://status.slack.com/api/v2.0.0/current")
	sf := func(serviceName string, useTLS bool) (url.URL, error) {
		return *u, err
	}
	c := client.NewBaseClient(sf, sn, true, 10*time.Second, nil)
	return slack.ReadStatus(c, sn, u.Path)
}
