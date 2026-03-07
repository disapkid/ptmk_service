package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"ptmk-service/internal/app"
	"ptmk-service/internal/repository"
	"ptmk-service/internal/service"
	"ptmk-service/pkg/api"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	address := os.Getenv("DATABASE_URL")
	if address == "" {
		address = "postgres://ptmk:123@localhost:5432/db?sslmode=disable"
	}
	ctx := context.Background()

	db, err := pgxpool.New(ctx, address)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	userRepo := repository.NewRepository(db)
	svc := service.NewService(userRepo)
	h := app.NewHandler(svc, userRepo)
	server, err := api.NewServer(h)
	
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(http.ListenAndServe(":8080", server))
}
