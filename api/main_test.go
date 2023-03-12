package api

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"github.com/thehaung/simplebank/config"
	db "github.com/thehaung/simplebank/db/sqlc"
	"github.com/thehaung/simplebank/util/randutil"
	"os"
	"testing"
	"time"
)

func newTestServer(t *testing.T, store db.Store) *Server {
	conf := &config.Config{
		TokenSymmetricKey:   randutil.StringWithQuantity(32),
		AccessTokenDuration: time.Minute,
	}

	server, err := NewHttpServer(conf, store)
	require.NoError(t, err)

	return server
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	os.Exit(m.Run())
}
