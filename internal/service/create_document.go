package service

import (
	"context"
	"ptmk-service/internal/model"
	"fmt"
)

func (s *Service) CreateDocument(ctx context.Context, document CreateDocument) error {
	doc := model.Document{
		UserID: uint64(document.UserID),
		DocumentType: document.DocumentType,
		DocumentNumber: document.DocumentNumber,
		DateCreated: document.DateCreated,
		DateEnd: document.DateEnd,
	}

	persons := make([]model.Party, 0, len(document.Persons))

	for _,val := range document.Persons {
		switch val.Type {
		case "LEGAL_ENTITY":
			if val.LegalEntity == nil {
    		    return fmt.Errorf("LEGAL_ENTITY: LegalEntity is nil")
    		}
			companyName := val.LegalEntity.Name

			persons = append(persons, model.Party{
				PartyType: "legal",
				CompanyName: &companyName,
				FirstName: nil,
				MiddleName: nil,
				LastName: nil,
				Initials: nil,
			})

		case "NATURAL_PERSON":
			if val.NaturalPerson == nil {
				return fmt.Errorf("NATURAL_PERSON: Natural person is nil")
			}

			persons = append(persons, model.Party{
				PartyType: "natural",
				CompanyName: nil,
				FirstName: val.NaturalPerson.FirstName,
				MiddleName: val.NaturalPerson.MiddleName,
				LastName: val.NaturalPerson.LastName,
				Initials: val.NaturalPerson.Initials,
			})

		default:
    		return fmt.Errorf("unknown person type: %q", val.Type)
		}
	}

	doc.Persons = persons
	err := s.userRepository.InsertDocument(ctx, doc)

	if err != nil {
		return err
	}

	return nil
}