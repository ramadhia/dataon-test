package http

import (
	"github.com/gin-gonic/gin"
)

func (h *DefaultHttpServer) setupRouting() {
	router := h.engine
	//cfg := h.config

	router.GET("/ping", func(context *gin.Context) {
		context.String(200, "Ok")
	})

	router.GET("users", h.handlers.user.FetchUser)
	router.POST("users", h.handlers.user.RegisterUser)
	router.PUT("users", h.handlers.user.Update)

	router.GET("organizations", h.handlers.organization.FetchComplete)

}
