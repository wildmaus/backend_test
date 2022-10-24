package server

import (
	"backend_test/internal/config"
	"backend_test/internal/storage/psql"
	"backend_test/pkg/client/postgressql"
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/rs/cors"
)

func Start(cfg *config.Config) error {

	pSQLClient, err := postgressql.NewClient(context.Background(), 3, cfg.Storage)
	if err != nil {
		log.Fatal(err)
	}
	storage := psql.NewStorage(pSQLClient)
	s := newServer(storage)
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.Server.Port))
	if err != nil {
		log.Fatal(err)
	}
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT"},
		AllowCredentials: true,
	})
	srv := &http.Server{
		Handler: c.Handler(s.router),
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Printf("server is listening port %s:%s", cfg.Server.Address, cfg.Server.Port)
	return srv.Serve(listener)
}
