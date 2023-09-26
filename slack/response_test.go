package slack

import (
	"net/http"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/sprak3000/go-client/client"
	"github.com/sprak3000/go-client/client/clientmock"
	"github.com/sprak3000/go-glitch/glitch"
	"github.com/sprak3000/go-whatsup-client/status"
)

func TestUnit_ReadStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sampleActiveResp := `{"status":"active","date_created":"2023-09-26T11:55:36-07:00","date_updated":"2023-09-26T11:55:36-07:00","active_incidents":[{"id":1269,"date_created":"2023-09-26T11:34:10-07:00","date_updated":"2023-09-26T11:55:36-07:00","title":"Can't dismiss Now you've got Later","type":"incident","status":"active","url":"https:\/\/status.slack.com\/2023-09\/badc5543a21e1fa7","services":["Connections"],"notes":[{"date_created":"2023-09-26T11:55:36-07:00","body":"sample note 2"},{"date_created":"2023-09-26T11:34:10-07:00","body":"sample note"}]}]}`
	sampleActiveStatusTime, _ := time.Parse(time.RFC3339, "2023-09-26T11:55:36-07:00")
	sampleActiveIncidentCreatedTime, _ := time.Parse(time.RFC3339, "2023-09-26T11:34:10-07:00")
	sampleActiveIncidentUpdatedTime, _ := time.Parse(time.RFC3339, "2023-09-26T11:55:36-07:00")

	tests := map[string]struct {
		expectedResp    Response
		expectedErr     glitch.DataError
		setupBaseClient func(t *testing.T, expectedErr glitch.DataError) client.BaseClient
		validate        func(t *testing.T, expectedResp, actualResp status.Details, expectedErr, actualErr glitch.DataError)
	}{
		"base path": {
			expectedResp: Response{
				Status:      "active",
				DateCreated: sampleActiveStatusTime,
				DateUpdated: sampleActiveStatusTime,
				ActiveIncidents: []Incident{
					{
						ID:          1269,
						DateCreated: sampleActiveIncidentCreatedTime,
						DateUpdated: sampleActiveIncidentUpdatedTime,
						Title:       "Can't dismiss Now you've got Later",
						Type:        "incident",
						Status:      "active",
						URL:         "https://status.slack.com/2023-09/badc5543a21e1fa7",
						Services:    []string{"Connections"},
						Notes: []Note{
							{
								DateCreated: sampleActiveIncidentUpdatedTime,
								Body:        "sample note 2",
							},
							{
								DateCreated: sampleActiveIncidentCreatedTime,
								Body:        "sample note",
							},
						},
					},
				},
			},
			setupBaseClient: func(t *testing.T, expectedErr glitch.DataError) client.BaseClient {
				c := clientmock.NewMockBaseClient(ctrl)
				c.EXPECT().MakeRequest(gomock.Any(), http.MethodGet, "test-slug", nil, gomock.Any(), nil).Return(http.StatusOK, []byte(sampleActiveResp), expectedErr)
				return c
			},
			validate: func(t *testing.T, expectedResp, actualResp status.Details, expectedErr, actualErr glitch.DataError) {
				require.Equal(t, expectedResp, actualResp)
				require.NoError(t, actualErr)
			},
		},
		"exceptional path- unable to make client request": {
			expectedResp: Response{},
			expectedErr:  glitch.NewDataError(nil, "test-err-code", "test-err"),
			setupBaseClient: func(t *testing.T, expectedErr glitch.DataError) client.BaseClient {
				c := clientmock.NewMockBaseClient(ctrl)
				c.EXPECT().MakeRequest(gomock.Any(), http.MethodGet, "test-slug", nil, gomock.Any(), nil).Return(http.StatusInternalServerError, nil, expectedErr)
				return c
			},
			validate: func(t *testing.T, expectedResp, actualResp status.Details, expectedErr, actualErr glitch.DataError) {
				require.Equal(t, expectedResp, actualResp)
				require.Error(t, actualErr)
				require.Equal(t, status.ErrorUnableToMakeClientRequest, actualErr.Code())
			},
		},
		"exceptional path- unable to parse response": {
			expectedResp: Response{},
			setupBaseClient: func(t *testing.T, expectedErr glitch.DataError) client.BaseClient {
				c := clientmock.NewMockBaseClient(ctrl)
				c.EXPECT().MakeRequest(gomock.Any(), http.MethodGet, "test-slug", nil, gomock.Any(), nil).Return(http.StatusOK, []byte(`{bad: ^JSON`), expectedErr)
				return c
			},
			validate: func(t *testing.T, expectedResp, actualResp status.Details, expectedErr, actualErr glitch.DataError) {
				require.Equal(t, expectedResp, actualResp)
				require.Error(t, actualErr)
				require.Equal(t, status.ErrorUnableToParseClientResponse, actualErr.Code())
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			resp, err := ReadStatus(tc.setupBaseClient(t, tc.expectedErr), "test-service", "test-slug")
			tc.validate(t, tc.expectedResp, resp, tc.expectedErr, err)
		})
	}
}

func TestUnit_Response_Indicator(t *testing.T) {
	tests := map[string]struct {
		resp              Response
		expectedIndicator string
		validate          func(t *testing.T, expectedIndicator, actualIndicator string)
	}{
		"base path- indicator for active status": {
			resp: Response{
				Status: "active",
			},
			expectedIndicator: "major",
			validate: func(t *testing.T, expectedIndicator, actualIndicator string) {
				require.Equal(t, expectedIndicator, actualIndicator)
			},
		},
		"base path- indicator for all other statuses": {
			resp: Response{
				Status: "copacetic",
			},
			expectedIndicator: "copacetic",
			validate: func(t *testing.T, expectedIndicator, actualIndicator string) {
				require.Equal(t, expectedIndicator, actualIndicator)
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			i := tc.resp.Indicator()
			tc.validate(t, tc.expectedIndicator, i)
		})
	}
}

func TestUnit_Response_Name(t *testing.T) {
	tests := map[string]struct {
		validate func(t *testing.T)
	}{
		"base path": {
			validate: func(t *testing.T) {
				resp := Response{}
				require.Equal(t, "Slack", resp.Name())
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tc.validate(t)
		})
	}
}

func TestUnit_Response_UpdatedAt(t *testing.T) {
	tests := map[string]struct {
		validate func(t *testing.T)
	}{
		"base path": {
			validate: func(t *testing.T) {
				ua := time.Now()
				resp := Response{DateUpdated: ua}
				require.Equal(t, ua, resp.UpdatedAt())
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tc.validate(t)
		})
	}
}
func TestUnit_Response_URL(t *testing.T) {
	tests := map[string]struct {
		validate func(t *testing.T)
	}{
		"base path": {
			validate: func(t *testing.T) {
				resp := Response{}
				require.Equal(t, "https://status.slack.com/", resp.URL())
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tc.validate(t)
		})
	}
}
