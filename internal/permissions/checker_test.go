package permissions

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGitHubCheckerMissingTokenDeniesWithoutError(t *testing.T) {
	checker := NewGitHubChecker(0)

	allowed, err := checker.CanWrite(context.Background(), Actor{Type: "github", ID: "user-without-token"}, "acme", "widgets")
	require.NoError(t, err)
	require.False(t, allowed)
}
