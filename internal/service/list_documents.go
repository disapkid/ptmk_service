package service

import (
	"context"
	"ptmk-service/internal/model"
)

func (s *Service) ListDocuments(ctx context.Context ,input model.ListDocuments) (model.DocumentResponse, error) {
	docResponse, err := s.userRepository.GetListDocuments(ctx, input)
	if err != nil {
		return model.DocumentResponse{}, err
	}
	return docResponse, nil
}