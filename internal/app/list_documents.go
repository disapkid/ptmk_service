package app

import (
	"context"
	"fmt"
	"ptmk-service/internal/model"
	"ptmk-service/pkg/api"
	"time"
)

func (h *Handler) ListDocuments(ctx context.Context, params api.ListDocumentsParams) (api.ListDocumentsRes, error) {	
	input := model.ListDocuments{
		UserID: uint64(params.UserID),
		Limit: &params.Limit.Value,
	}

	var isCorrectSortType bool

	for _, sortType := range model.SortType {
		if params.Sort == sortType {
			input.Sort = sortType
			isCorrectSortType = true
		}
	}

	if !isCorrectSortType {
		return &api.ListDocumentsBadRequest{},nil
	}

	docResponse, err := h.Svc.ListDocuments(ctx, input)
	if err != nil {
		return &api.ListDocumentsInternalServerError{}, nil
	}

	docs, err := convertType(docResponse)
	if err != nil {
		return &api.ListDocumentsInternalServerError{}, nil
	}

	return &api.ListDocumentsOK{
		Total: int(docResponse.Total),
		Documents: docs,
	}, nil
}

func convertType(input model.DocumentResponse) ([]api.DocumentResponse, error) {
	result := make([]api.DocumentResponse, 0)

	for _ , val := range input.Documents {
		docResponse := api.DocumentResponse{
			UserID: 		int64(val.UserID),
			DocumentType: 	val.DocumentType,
			DocumentNumber: val.DocumentNumber,
			DateCreated: 	val.DateCreated,
			DateEnd: 		optNilDateConv(val.DateEnd),
		}

		for _, person := range val.Persons {
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
				return nil, fmt.Errorf("Unknown type: %q", person.PartyType)
			}

			docResponse.Persons = append(docResponse.Persons, persons)
		}
		
		result = append(result, docResponse)
	}

	return result, nil
}

func optNilDateConv(date *time.Time) api.OptNilDate {
	if date == nil {
		return api.OptNilDate{
			Value: time.Time{},
			Set: false,
			Null: true,
		}
	}

	return api.NewOptNilDate(*date)
}

func optLegalEntityConv(input model.Party) api.OptLegalEntity {
	if input.CompanyName == nil {
		return api.OptLegalEntity{}
	}

	v := api.LegalEntity{
		Name:       *input.CompanyName,
		EntityType: api.NewOptString("legal"),
	}

	return api.NewOptLegalEntity(v)
}


func optNaturalPersonConv(input model.Party) api.OptNaturalPerson {
	v := api.NaturalPerson{}

	if input.FirstName != nil {
		v.FirstName = api.NewOptNilString(*input.FirstName)
	} else {
		v.FirstName = api.OptNilString{}
	}
	if input.MiddleName != nil {
		v.MiddleName = api.NewOptNilString(*input.MiddleName)
	} else {
		v.MiddleName = api.OptNilString{}
	}
	if input.LastName != nil {
		v.LastName = api.NewOptNilString(*input.LastName)
	} else {
		v.LastName = api.OptNilString{}
	}
	if input.Initials != nil {
		v.Initials = api.NewOptNilString(*input.Initials)
	} else {
		v.Initials = api.OptNilString{}
	}

	return api.NewOptNaturalPerson(v)
}