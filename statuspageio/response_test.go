package statuspageio

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

	sampleOKResp := `{"page":{"id":"kctbh9vrtdwd","name":"GitHub","url":"https://www.githubstatus.com","time_zone":"Etc/UTC","updated_at":"2023-09-26T07:51:43.965Z"},"status":{"indicator":"none","description":"All Systems Operational"}}`
	sampleOKTime, _ := time.Parse(time.RFC3339, "2023-09-26T07:51:43.965Z")

	tests := map[string]struct {
		expectedResp    Response
		expectedErr     glitch.DataError
		setupBaseClient func(t *testing.T, expectedErr glitch.DataError) client.BaseClient
		validate        func(t *testing.T, expectedResp, actualResp status.Details, expectedErr, actualErr glitch.DataError)
	}{
		"base path": {
			expectedResp: Response{
				Status: Status{
					Indicator:   "none",
					Description: "All Systems Operational",
				},
				Page: Page{
					ID:        "kctbh9vrtdwd",
					Name:      "GitHub",
					URL:       "https://www.githubstatus.com",
					TimeZone:  "",
					UpdatedAt: sampleOKTime,
				},
			},
			setupBaseClient: func(t *testing.T, expectedErr glitch.DataError) client.BaseClient {
				c := clientmock.NewMockBaseClient(ctrl)
				c.EXPECT().MakeRequest(gomock.Any(), http.MethodGet, "test-slug", nil, gomock.Any(), nil).Return(http.StatusOK, []byte(sampleOKResp), expectedErr)
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
		"base path": {
			resp: Response{
				Status: Status{
					Indicator: "major",
				},
			},
			expectedIndicator: "major",
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
		resp         Response
		expectedName string
		validate     func(t *testing.T, expectedName, actualName string)
	}{
		"base path": {
			resp: Response{
				Page: Page{
					Name: "test-service",
				},
			},
			expectedName: "test-service",
			validate: func(t *testing.T, expectedName, actualName string) {
				require.Equal(t, expectedName, actualName)
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			n := tc.resp.Name()
			tc.validate(t, tc.expectedName, n)
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
				resp := Response{
					Page: Page{
						UpdatedAt: ua,
					},
				}
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
				resp := Response{
					Page: Page{
						URL: "https://foo.test/",
					},
				}
				require.Equal(t, "https://foo.test/", resp.URL())
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tc.validate(t)
		})
	}
}
