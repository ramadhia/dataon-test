package middleware

import (
	"github.com/ramadhia/bosnet/be/internal/config"

	"github.com/gin-contrib/cors"
)

func CorsPolicy(config config.Config) cors.Config {
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowCredentials = true
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Authorization", "Content-Type"}

	return corsConfig
}
