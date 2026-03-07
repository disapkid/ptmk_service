package app

import (
	"context"
	"errors"
	"fmt"
	"ptmk-service/internal/model"
	"ptmk-service/pkg/api"
)

func (h *Handler) GetDocumentByID(ctx context.Context, params api.GetDocumentByIDParams) (api.GetDocumentByIDRes, error) {
	document, err := h.Svc.GetDocumentByID(ctx, params.ID)

	if err != nil {
		if errors.Is(err, model.ErrorDocumentNotFound) {
			return &api.GetDocumentByIDNotFound{}, nil
		}

		return &api.GetDocumentByIDInternalServerError{}, nil
	}

	result, err := convertDocument(document)

	if err != nil {
		return &api.GetDocumentByIDInternalServerError{}, nil
	}

	return &result, nil
}

func convertDocument(input model.Document) (api.DocumentResponse, error) {
	result := api.DocumentResponse{
		UserID: int64(input.UserID),
		DocumentType: input.DocumentType,
		DocumentNumber: input.DocumentNumber,
		DateCreated: input.DateCreated,
		DateEnd: optNilDateConv(input.DateEnd),
	}

	for _, person := range input.Persons {
		var personType api.PersonType
		var persons api.Person

		switch person.PartyType {
		case "legal":
			personType = api.PersonTypeLEGALENTITY
			persons = api.Person{
				Type: personType,
				LegalEntity: optLegalEntityConv(person),
				NaturalPerson: api.OptNaturalPerson{},
			}
		case "natural":
			personType = api.PersonTypeNATURALPERSON
			persons = api.Person{
				Type: personType,
				LegalEntity: api.OptLegalEntity{},
				NaturalPerson: optNaturalPersonConv(person),
			}
		default:
			return api.DocumentResponse{}, fmt.Errorf("unknown party type: %q", person.PartyType)
		}

		result.Persons = append(result.Persons, persons)
	}

	return result, nil
}