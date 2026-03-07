package repository

import (
	"context"
	"ptmk-service/internal/model"
	"time"

	sq "github.com/Masterminds/squirrel"
)

func (r *UserRepository) GetDocumentByID(ctx context.Context, ID int64) (model.Document, error) {
	selectReq := []string{
		"documents.id AS doc_id",	"documents.user_id",		"documents.doc_type",
		"documents.doc_number",		"documents.start_date",		"documents.expire_date",
		"doc_parties.party_type",	"doc_parties.company_name",	"doc_parties.first_name",
		"doc_parties.middle_name",	"doc_parties.last_name",	"doc_parties.initials",
	}

	query := sq.
		Select(selectReq...).
		From(Documents).
		LeftJoin("doc_parties ON doc_parties.doc_id = documents.id").
		Where(sq.Eq{ "documents.id" : ID}).
		PlaceholderFormat(sq.Dollar)
	
	sql, args, err := query.ToSql()
	if err != nil {
		return model.Document{}, err
	}

	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		return  model.Document{}, err
	} 
	defer rows.Close()

	var doc model.Document
	isFound := false

	for rows.Next() {
		var (
			docID 		int64
			userID 		int64
			docType 	string
			docNumber 	string
			startDate 	time.Time
			endDate 	*time.Time

			partyType 	*string
			companyName *string
			firstName 	*string
			middleName 	*string
			lastName 	*string
			initials 	*string
		)

		err := rows.Scan(
			&docID, 	&userID,
			&docType, 	&docNumber,
			&startDate, &endDate,
			&partyType, &companyName,
			&firstName, &middleName,
			&lastName, 	&initials,
		)

		if err != nil {
			return model.Document{}, err
		}

		if !isFound {
			doc = model.Document{
				UserID: uint64(userID),
				DocumentType: docType,
				DocumentNumber: docNumber,
				DateCreated: startDate,
				DateEnd: endDate,
			}

			isFound = true
		}

		if partyType != nil {
			doc.Persons = append(doc.Persons, model.Party{
				PartyType: *partyType,
				CompanyName: companyName,
				FirstName: firstName,
				MiddleName: middleName,
				LastName: lastName,
				Initials: initials,
			})
		}
	}

	err = rows.Err()
	if err != nil {
		return model.Document{}, err
	}

	if !isFound {
		return model.Document{}, model.ErrorDocumentNotFound
	}

	return doc , nil
}