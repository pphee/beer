package servers

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/peedans/beerleo/config"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type IServer interface {
	Start()
}

type server struct {
	app *gin.Engine
	cfg config.IConfig
	db  *sqlx.DB
}

func NewServer(cfg config.IConfig, db *sqlx.DB) IServer {
	gin.SetMode(gin.ReleaseMode)

	app := gin.Default()
	return &server{
		cfg: cfg,
		db:  db,
		app: app,
	}
}

func (s *server) Start() {

	v1 := s.app.Group("v1")
	modules := InitModule(v1, s)

	modules.monitorModule()
	modules.beersleoModule()
	// Graceful Shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	srv := &http.Server{
		Addr:    s.cfg.App().Url(),
		Handler: s.app,
	}
	go func() {
		<-c // Wait for an interrupt signal.

		log.Println("Received interrupt. Shutting down servers...")

		// Create a context with a 5-second timeout to allow for graceful shutdown.
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Attempt to gracefully shutdown the servers.
		if err := srv.Shutdown(ctx); err != nil {
			log.Fatalf("Server shutdown failed: %v", err)
		} else {
			log.Println("Server shutdown successfully.")
		}
	}()

	// Listen to host:port
	log.Printf("servers is starting on %v", s.cfg.App().Url())
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("listen: %s\n", err)
	}

}
