package api

import (
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	db "github.com/lfvm/simplebank/db/sqlc"
	"github.com/lfvm/simplebank/utils"
	"github.com/stretchr/testify/require"
)

func newTestServer(t *testing.T, store db.Store) *Server {
	config := utils.Config{
		TokenSymetricKey: utils.RandomString(32),
		ServerAddress:    "0.0.0.0:8080",
	}

	server, err := NewServer(config, store)

	require.NoError(t, err)
	return server
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}
