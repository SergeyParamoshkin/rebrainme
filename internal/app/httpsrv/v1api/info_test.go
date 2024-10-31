package v1api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/SergeyParamoshkin/alerts/pkg/info"
)

func TestInfo_RenderCorrectResponseFormat(t *testing.T) {
	api := API{}
	req := httptest.NewRequest(http.MethodGet, "/info", nil)
	rec := httptest.NewRecorder()

	api.info(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))

	var infoResp infoResponse
	err := json.NewDecoder(resp.Body).Decode(&infoResp)
	require.NoError(t, err)

	assert.Equal(t, info.CommitSHA, infoResp.CommitSHA)
}
