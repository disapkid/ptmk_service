package service

import (
	"context"
	"ptmk-service/internal/model"
)

func (s *Service) GetDocumentByID(ctx context.Context, ID int64) (model.Document, error) {
	result, err := s.userRepository.GetDocumentByID(ctx, ID)
	if err != nil {
		return model.Document{}, err
	}
	return result, nil
}