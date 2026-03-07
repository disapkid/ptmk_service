package app

import (
	"ptmk-service/internal/repository"
	"ptmk-service/internal/service"
)

type Handler struct{
	Svc *service.Service
	Repo *repository.UserRepository
}

func NewHandler(svc *service.Service, repo *repository.UserRepository) *Handler {
	return &Handler{
		Svc: svc,
		Repo: repo,
	}
}
