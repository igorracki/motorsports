package clients

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestF1DataClient_Validation(t *testing.T) {
	client := NewF1DataClient("http://localhost:8080")
	ctx := context.Background()

	t.Run("GetScheduleByYear - Invalid Year", func(tt *testing.T) {
		// When
		res, err := client.GetScheduleByYear(ctx, 1949)
		// Then
		assert.Error(tt, err)
		assert.Nil(tt, res)
		assert.Contains(tt, err.Error(), "invalid year parameter")

		// When
		res, err = client.GetScheduleByYear(ctx, 2101)
		// Then
		assert.Error(tt, err)
		assert.Nil(tt, res)
		assert.Contains(tt, err.Error(), "invalid year parameter")
	})

	t.Run("GetSessionResults - Invalid Year/Round", func(tt *testing.T) {
		// When
		res, err := client.GetSessionResults(ctx, 1949, 1, "race")
		// Then
		assert.Error(tt, err)
		assert.Nil(tt, res)

		// When
		res, err = client.GetSessionResults(ctx, 2024, 0, "race")
		// Then
		assert.Error(tt, err)
		assert.Nil(tt, res)

		// When
		res, err = client.GetSessionResults(ctx, 2024, 51, "race")
		// Then
		assert.Error(tt, err)
		assert.Nil(tt, res)
	})

	t.Run("GetCircuit - Invalid Year/Round", func(tt *testing.T) {
		// When
		res, err := client.GetCircuit(ctx, 1949, 1)
		// Then
		assert.Error(tt, err)
		assert.Nil(tt, res)
	})
}

func TestF1DataClient_URLEscaping(t *testing.T) {
	// Given
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/results/2024/1/race type", r.URL.Path)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"year": 2024, "round": 1, "session_type": "race type", "results": []}`))
	}))
	defer server.Close()

	client := NewF1DataClient(server.URL)
	ctx := context.Background()

	// When
	_, err := client.GetSessionResults(ctx, 2024, 1, "race type")

	// Then
	assert.NoError(t, err)
}
