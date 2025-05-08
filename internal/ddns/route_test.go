package ddns

import (
	"github.com/gin-gonic/gin"
	"testing"
)

func TestRouter(t *testing.T) {
	e := gin.Default()
	g := e.Group("/")
	LoadRouter(g)
	_ = e.Run("127.0.0.1:5801")
}
