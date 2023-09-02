package api

import (
	"github.com/gin-gonic/gin"
)

func Setup() *gin.Engine {
	r := gin.Default()
	r.SetTrustedProxies([]string{"127.0.0.1", "::1"})

	r.GET("/health_check", healthCheck)
	r.PATCH("/providers/:providerType/zones/:zoneId/records/:id", updateZoneRecord)

	return r
}

func healthCheck(c *gin.Context) {
	c.String(200, "OK")
}

func updateZoneRecord(c *gin.Context) {
	zoneId := c.Param("zoneId")
	id := c.Param("id")
	ip := c.ClientIP()

	// NOTE only a single DNS provider is currently supported
	switch provider := c.Param("providerType"); provider {
	case "cloudflare":
		err := cloudflareUpdateRecord(c, zoneId, id, ip)
		if err != nil {
			c.String(400, err.Error())
			return
		}

	default:
		c.String(400, "Unsupported provider")
		return
	}

	c.Status(204)
}
