package http

import (
	"github.com/gin-gonic/gin"
)

func (h *DefaultHttpServer) setupRouting() {
	router := h.engine

	// middleware groups
	var secureMiddlewares []gin.HandlerFunc

	router.GET("/ping", func(context *gin.Context) {
		context.String(200, "Ok")
	})

	router.GET("history", h.handlers.history.FetchHistory)
	router.PUT("transfer", h.handlers.history.Transfer)
	router.PUT("setor", h.handlers.history.Setor)
	router.PUT("tarik", h.handlers.history.Tarik)

	router.Use(secureMiddlewares...)
	{
		//router.PUT("traffics", h.handlers.traffic.UpsertTraffic)
		//router.DELETE("traffics/:id", h.handlers.traffic.DeleteTraffic)
	}

}
