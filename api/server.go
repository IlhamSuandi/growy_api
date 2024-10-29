package api

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ilhamSuandi/business_assistant/api/middleware"
	"github.com/ilhamSuandi/business_assistant/api/route"
	"github.com/ilhamSuandi/business_assistant/utils"
	"golang.org/x/time/rate"
	"gorm.io/gorm"
)

type APIServer struct {
	db   *gorm.DB
	addr string
}

func NewApiServer(addr string, db *gorm.DB) *APIServer {
	return &APIServer{
		addr: addr,
		db:   db,
	}
}

var server http.Server

func (s *APIServer) Start() error {
	router := route.RegisterRoutes(s.db)

	// enable global cors
	handler := middleware.Cors(router)

	// enable global rate limiter
	limit := rate.Every(time.Second / 3) // 3 requests per second
	handler = middleware.Limiter(limit, nil, handler)

	// enable global logger
	handler = middleware.Log(s.db, handler)

	// configure global server
	server = http.Server{
		Addr:         s.addr,
		Handler:      handler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// log server start
	utils.Log.Infof("server is running on localhost%s", s.addr)

	// start server
	return server.ListenAndServe()
}

func (s *APIServer) GracefulShutdown() {
	// Create a channel to listen for interrupt signals
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// Wait for a signal
	<-sigs
	utils.Log.Info("Shutting down server...")

	// Create a context with a timeout for the shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Shutdown the server gracefully
	if err := server.Shutdown(ctx); err != nil {
		utils.Log.Errorf("Server forced to shutdown: %v", err)
	}

	utils.Log.Info("Server exited gracefully")
}
