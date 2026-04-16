package auth

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestDeviceFlowRoundTrip(t *testing.T) {
	var polls int
	var server *httptest.Server
	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/login/device/code":
			require.Equal(t, http.MethodPost, r.Method)
			_ = json.NewEncoder(w).Encode(DeviceCodeResponse{
				DeviceCode:      "device-123",
				UserCode:        "ABCD-EFGH",
				VerificationURI: server.URL + "/login/device",
				ExpiresIn:       60,
				Interval:        1,
			})
		case "/login/oauth/access_token":
			polls++
			if polls == 1 {
				_ = json.NewEncoder(w).Encode(AccessTokenResponse{
					Error:            "authorization_pending",
					ErrorDescription: "still waiting",
				})
				return
			}
			_ = json.NewEncoder(w).Encode(AccessTokenResponse{
				AccessToken: "gho_test",
				TokenType:   "bearer",
				Scope:       "read:org,repo",
			})
		case "/user":
			require.Equal(t, "Bearer gho_test", r.Header.Get("Authorization"))
			_ = json.NewEncoder(w).Encode(Viewer{
				Login: "bob",
				ID:    42,
			})
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	cfg := Config{
		ClientID:     "client-123",
		Scope:        "read:org repo",
		OAuthBaseURL: server.URL,
		APIBaseURL:   server.URL,
		HTTPClient:   server.Client(),
	}

	device, err := cfg.StartDeviceFlow(context.Background())
	require.NoError(t, err)
	require.Equal(t, "device-123", device.DeviceCode)

	token, err := cfg.PollAccessToken(context.Background(), device.DeviceCode, 0, time.Minute)
	require.NoError(t, err)
	require.Equal(t, "gho_test", token.AccessToken)

	viewer, err := cfg.GetViewer(context.Background(), token.AccessToken)
	require.NoError(t, err)
	require.Equal(t, "bob", viewer.Login)
	require.EqualValues(t, 42, viewer.ID)
}

func TestStoredTokenRoundTrip(t *testing.T) {
	tempDir := t.TempDir()
	t.Setenv("PRTAGS_CONFIG_DIR", tempDir)

	path, err := SaveStoredToken(StoredToken{
		AccessToken: "gho_saved",
		TokenType:   "bearer",
		Scope:       "read:org,repo",
		UserLogin:   "bob",
		UserID:      42,
		ClientID:    "client-123",
	})
	require.NoError(t, err)
	require.Equal(t, filepath.Join(tempDir, "auth.json"), path)

	info, err := os.Stat(path)
	require.NoError(t, err)
	require.Equal(t, os.FileMode(0o600), info.Mode().Perm())

	token, err := LoadStoredToken()
	require.NoError(t, err)
	require.Equal(t, "gho_saved", token.AccessToken)
	require.Equal(t, "bob", token.UserLogin)

	require.NoError(t, DeleteStoredToken())
	_, err = LoadStoredToken()
	require.Error(t, err)
}
