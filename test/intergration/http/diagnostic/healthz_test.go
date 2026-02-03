package diagnostic

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

// nolint:goconst
func TestHealthz(t *testing.T) {
	host := "app"
	httpURL := "http://" + host + ":8081"
	url := httpURL + "/healthz"

	req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, url, http.NoBody)
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer func() {
		err = resp.Body.Close()
		require.NoError(t, err)
	}()

	require.Equal(t, http.StatusOK, resp.StatusCode)
}
