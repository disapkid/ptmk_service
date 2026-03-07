package app

import (
	"context"
	"errors"
	"ptmk-service/internal/model"
	"ptmk-service/pkg/api"
)

func (h *Handler) DeleteDocumentByNumber(ctx context.Context, params api.DeleteDocumentByNumberParams) (api.DeleteDocumentByNumberRes, error) {
	err := h.Repo.DeleteDocumentByNumber(ctx, params.DocumentNumber)
	if err != nil {
		if errors.Is(err, model.ErrorDocumentNotFound) {
			return &api.DeleteDocumentByNumberNotFound{}, nil
		}

		return &api.DeleteDocumentByNumberInternalServerError{}, nil
	}
	
	return &api.DeleteDocumentByNumberNoContent{}, nil
}
