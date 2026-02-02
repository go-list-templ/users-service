package diagnostic

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// nolint:goconst
func TestHealthz(t *testing.T) {
	host := "app"
	httpURL := host + ":8080"
	url := httpURL + "/healthz"
	requestTimeout := 5 * time.Second

	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer func() {
		err = resp.Body.Close()
		require.NoError(t, err)
	}()

	require.Equal(t, http.StatusOK, resp.StatusCode)
}
