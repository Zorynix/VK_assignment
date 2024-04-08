package tests

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"vk.com/m/services"
)

func TestActorAdd(t *testing.T) {

	ctx := context.Background()

	// Start a new PostgreSQL container
	containerReq := testcontainers.ContainerRequest{ // Renamed variable here
		Image:        "postgres:13",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_DB":       "vk",
			"POSTGRES_USER":     "user",
			"POSTGRES_PASSWORD": "password",
		},
		WaitingFor: wait.ForLog("database system is ready to accept connections").WithStartupTimeout(120 * time.Second),
	}
	postgresC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: containerReq,
		Started:          true,
	})
	require.NoError(t, err)
	defer postgresC.Terminate(ctx)

	actorToAdd := `{"name":"John Doe","date_of_birth":"2000-01-01"}`
	reqBody := strings.NewReader(actorToAdd)
	req, err := http.NewRequest("POST", "/v1/actor-add", reqBody)
	require.NoError(t, err)
	w := httptest.NewRecorder()

	pg, err := services.NewPostgreSQL(ctx)
	if err != nil {
		return
	}

	_, err = pg.ActorAdd(w, req)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, w.Code)

	// Optionally, retrieve and assert the added actor details from the database for completeness
}
