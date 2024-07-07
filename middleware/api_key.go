package middleware

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rangganovsky/go-billing-engine/config"
)

func StaticAPIKey() gin.HandlerFunc {
	conf, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("Could not load environment variable", err)
	}

	apiKey := conf.APIKEY

	return func(c *gin.Context) {
		reqKey := c.Request.Header.Get("X-API-Key")

		if reqKey != apiKey {
			log.Println("api key empty or not valid")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": 401, "message": "Authentication failed"})
			return
		}
	}
}
