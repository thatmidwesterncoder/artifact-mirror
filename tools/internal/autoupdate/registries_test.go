package autoupdate

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGitHubRegistryFetchTags(t *testing.T) {
	tests := []struct {
		name                   string
		mockResponseStatus     int
		mockResponseHeaderLink string
		mockResponseBody       string
		expectedTags           []string
		expectedNext           string
		wantErr                bool
		errorMessage           string
	}{
		{
			name:                   "should return tags list for successful Github fetchTags request",
			mockResponseStatus:     http.StatusOK,
			mockResponseHeaderLink: "",
			mockResponseBody:       `{"tags": ["v1.0.0","v1.1.0","v2.0.0","v2.0.1"]}`,
			expectedTags:           []string{"v1.0.0", "v1.1.0", "v2.0.0", "v2.0.1"},
			expectedNext:           "",
			wantErr:                false,
			errorMessage:           "",
		},
		{
			name:                   "should return correct next link for successful Github fetchTags request",
			mockResponseStatus:     http.StatusOK,
			mockResponseHeaderLink: `</v2/owner/repo/tags/list?last=v1.10.0-rc1-5-g0d7c9121&n=200>; rel="next"`,
			mockResponseBody:       `{"tags": ["v1.0.0"]}`,
			expectedTags:           []string{"v1.0.0"},
			expectedNext:           "https://ghcr.io/v2/owner/repo/tags/list?last=v1.10.0-rc1-5-g0d7c9121&n=200",
			wantErr:                false,
			errorMessage:           "",
		},
		{
			name:                   "should return error for unsuccessful Github fetchTags request",
			mockResponseStatus:     http.StatusNotFound,
			mockResponseHeaderLink: "",
			mockResponseBody:       "not found",
			expectedTags:           nil,
			expectedNext:           "",
			wantErr:                true,
			errorMessage:           "failed with status 404 Not Found and body not found",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Link", tt.mockResponseHeaderLink)
				w.WriteHeader(tt.mockResponseStatus)
				_, err := w.Write([]byte(tt.mockResponseBody))
				if err != nil {
					t.Errorf("fetchTags() failed: %v", err)
				}
			}))
			defer server.Close()

			g := &GitHubRegistry{}
			tags, next, gotErr := g.fetchTags("<token>", server.URL)
			if tt.wantErr {
				assert.ErrorContains(t, gotErr, tt.errorMessage)
			}
			assert.Equal(t, tt.expectedTags, tags, "fetchTags() = %v, want %v", tags, tt.expectedTags)
			assert.Equal(t, tt.expectedNext, next, "fetchTags() = %v, want %v", next, tt.expectedNext)
		})
	}
}
