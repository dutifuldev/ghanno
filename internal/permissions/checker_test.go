package permissions

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGitHubCheckerMissingTokenDeniesWithoutError(t *testing.T) {
	checker := NewGitHubChecker(0)

	allowed, err := checker.CanWrite(context.Background(), Actor{Type: "github", ID: "user-without-token"}, "acme", "widgets")
	require.NoError(t, err)
	require.False(t, allowed)
}

func TestGitHubCheckerTreatsNotFoundAndUnauthorizedAsDenied(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/repos/acme/widgets") {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer server.Close()

	t.Setenv("GITHUB_API_URL", server.URL+"/")

	checker := NewGitHubChecker(0)
	allowed, err := checker.CanWrite(context.Background(), Actor{Type: "github", ID: "tester", Token: "token"}, "acme", "widgets")
	require.NoError(t, err)
	require.False(t, allowed)
}

func TestGitHubCheckerResolvesIdentity(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/user") {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"id":7937614,"login":"dutifulbob"}`))
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	t.Setenv("GITHUB_API_URL", server.URL+"/")

	checker := NewGitHubChecker(0)
	identity, err := checker.ResolveIdentity(context.Background(), Actor{Type: "github", ID: "tester", Token: "token"})
	require.NoError(t, err)
	require.EqualValues(t, 7937614, identity.GitHubUserID)
	require.Equal(t, "dutifulbob", identity.GitHubLogin)
}
