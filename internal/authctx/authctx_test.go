package authctx

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestWithPrincipalAndPrincipalFrom(t *testing.T) {
	principal := Principal{
		UserUUID: uuid.New(),
		Login:    "max",
	}

	ctx := WithPrincipal(context.Background(), principal)
	result, ok := PrincipalFrom(ctx)

	require.True(t, ok)
	require.Equal(t, principal, result)
}

func TestPrincipalFromMissing(t *testing.T) {
	_, ok := PrincipalFrom(context.Background())

	require.False(t, ok)
}
