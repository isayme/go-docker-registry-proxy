package src

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// version info
var (
	Name        = "@isayme/go-docker-registry-proxy"
	Version     = "unknown"
	BuildTime   = "unknown"
	GitRevision = "unknown"
)

// PrintVersion print version
func PrintVersion(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"name":        Name,
		"version":     Version,
		"buildTime":   BuildTime,
		"gitRevision": GitRevision,
	})
}
