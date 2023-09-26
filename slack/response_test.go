package slack

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

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
