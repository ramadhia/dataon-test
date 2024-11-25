package http

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/ramadhia/dataon-test/internal/config"
	"github.com/ramadhia/dataon-test/internal/handler/http/handler"
	"github.com/ramadhia/dataon-test/internal/handler/http/middleware"
	"github.com/ramadhia/dataon-test/internal/provider"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type HttpServerImpl interface {
	Start() error
	Stop() error
	GetHandler() (http.Handler, error)
}

type DefaultHttpServer struct {
	config   config.Config
	engine   *gin.Engine
	handlers handlers
}

type handlers struct {
	organization handler.Organization
	user         handler.User
}

func NewHttpServer(p *provider.Provider) *DefaultHttpServer {
	gin.SetMode(gin.ReleaseMode)
	if strings.ToLower(p.Config().Log.Level) == gin.DebugMode {
		gin.SetMode(gin.DebugMode)
	}

	engine := gin.New()

	engine.Use(middleware.LogrusLogger(logrus.StandardLogger()), gin.Recovery())
	engine.Use(cors.New(middleware.CorsPolicy(p.Config())))

	// ini the handler
	handlers := handlers{
		organization: *handler.NewOrganization(p),
		user:         *handler.NewUser(p),
	}

	requestHandler := &DefaultHttpServer{p.Config(), engine, handlers}
	requestHandler.setupRouting()

	return requestHandler
}

func (h *DefaultHttpServer) Start() error {
	logger := logrus.WithField("method", "web.httpServer.Start")

	// wait for interrupt signal to gracefully shut down the server
	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", h.config.App.Host, h.config.App.Port),
		Handler: h.engine,
	}

	done := make(chan os.Signal)

	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	// Initializing the server in a go routine
	go func() {
		logger.Infof("About to listen at: %s", server.Addr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Errorf("Server closed caused by unhandled error: %s\n", err)
		}
	}()

	<-done

	logger.Println("Shutting down server...")

	// The context is used to inform the server it has 60 seconds to finish
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(5)*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Warningf("Error when shutting down: %s", err.Error())
		return err
	}

	logger.Info("Server exiting")
	return nil
}

func (h *DefaultHttpServer) GetHandler() (http.Handler, error) {
	return h.engine, nil
}
