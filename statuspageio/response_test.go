package statuspageio

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
