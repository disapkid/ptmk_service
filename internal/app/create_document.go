package app

import (
	"context"
	"fmt"
	service "ptmk-service/internal/service"
	"ptmk-service/pkg/api"
	"time"
)

func (h *Handler) CreateDocument(ctx context.Context, request *api.DocumentCreateRequest) (api.CreateDocumentRes, error) {
	input := service.CreateDocument{
		UserID: request.UserID,
		DocumentType: request.DocumentType,
		DocumentNumber: request.DocumentNumber,
		DateCreated: request.DateCreated,
	}
	persons,err := parsePerson(request)
	if err != nil {
		return &api.CreateDocumentBadRequest{}, nil
	}

	input.Persons = persons
	input.DateEnd = optNilDatePtr(request.DateEnd)

	err = h.Svc.CreateDocument(ctx, input)
	if err != nil {
		return &api.CreateDocumentInternalServerError{}, nil
	}

	return &api.DocumentResponse{
			UserID: request.UserID,
			DocumentType: request.DocumentType,
			DocumentNumber: request.DocumentNumber,
			DateCreated: request.DateCreated,
			DateEnd: request.DateEnd,
			Persons: request.Persons,
		}, nil
}

func parsePerson (req *api.DocumentCreateRequest) ([]service.Person, error) {
	res := make([]service.Person, 0, len(req.Persons))

	for _, val := range req.Persons {
		valType := val.GetType()
		switch valType {

		case api.PersonTypeLEGALENTITY:
			legalType := string(api.PersonTypeLEGALENTITY)

			var companyName string
			if !val.LegalEntity.Set {
				return nil, fmt.Errorf("legal entity data is missing")
			}
			companyName = val.LegalEntity.Value.Name

			res = append(res, service.Person{
				Type: legalType,
				LegalEntity: &service.LegalEntity{ Name: companyName },
				NaturalPerson: nil,
			})

		case api.PersonTypeNATURALPERSON:
			naturalPersonType := string(api.PersonTypeNATURALPERSON)
			if !val.NaturalPerson.Set {
				return nil, fmt.Errorf("natural person data is missing")
			}

			firstName 	:= optNilStringPtr(val.NaturalPerson.Value.FirstName)
			middleName 	:= optNilStringPtr(val.NaturalPerson.Value.MiddleName)
			lastName 	:= optNilStringPtr(val.NaturalPerson.Value.LastName)
			initials 	:= optNilStringPtr(val.NaturalPerson.Value.Initials)

			res = append(res, service.Person{
				Type: naturalPersonType,
				NaturalPerson: &service.NaturalPerson{
					FirstName: 	firstName,
					MiddleName: middleName,
					LastName: 	lastName,
					Initials: 	initials,
				},
			})
		default:
			return nil, fmt.Errorf("Unexpected type of person: %v", valType)
		}
	}

	return res, nil
}

func optNilDatePtr(v api.OptNilDate) *time.Time {
	value, ok := v.Get()
	if !ok {
		return nil
	}

	valueCopy := value
	return &valueCopy
}

func optNilStringPtr(v api.OptNilString) *string {
	value, ok := v.Get()
	if !ok {
		return nil
	}

	valueCopy := value
	return &valueCopy
}