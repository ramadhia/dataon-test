package http

import (
	"github.com/gin-gonic/gin"
	"github.com/ramadhia/mnc-test/internal/handler/http/middleware"
)

func (h *DefaultHttpServer) setupRouting() {
	router := h.engine
	cfg := h.config

	// middleware groups
	secureMiddlewares := []gin.HandlerFunc{
		middleware.NewHmacJwtMiddleware([]byte(cfg.App.JwtSecret)),
	}

	router.GET("/ping", func(context *gin.Context) {
		context.String(200, "Ok")
	})

	router.POST("register", h.handlers.user.RegisterUser)
	router.POST("login", h.handlers.user.Login)

	router.Use(secureMiddlewares...)
	{
		router.POST("topup", h.handlers.transaction.Topup)
		router.POST("pay", h.handlers.transaction.Pay)
		router.POST("transfer", h.handlers.transaction.Transfer)

		router.POST("transafer-test", h.handlers.transaction.TransferTest)

		router.GET("transactions", h.handlers.transaction.FetchTransaction)
		router.GET("profile", h.handlers.user.GetProfile)
		router.PUT("profile", h.handlers.user.Update)
	}

}
