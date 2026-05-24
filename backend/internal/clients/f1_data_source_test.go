package clients

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestF1DataClient_URLEscaping(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/results/2024/1/race type", r.URL.Path)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"year": 2024, "round": 1, "session_type": "race type", "results": []}`))
	}))
	defer server.Close()

	client := NewF1DataClient(server.URL)
	ctx := context.Background()

	_, err := client.GetSessionResults(ctx, 2024, 1, "race type")

	assert.NoError(t, err)
}

func TestF1DataClient_ErrorHandling(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"message": "not found"}`))
	}))
	defer server.Close()

	client := NewF1DataClient(server.URL)
	ctx := context.Background()

	res, err := client.GetScheduleByYear(ctx, 2024)

	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Contains(t, err.Error(), "remote resource not found")
}
