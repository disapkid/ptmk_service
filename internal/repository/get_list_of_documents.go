package repository

import (
	"context"
	"ptmk-service/internal/model"

	sq "github.com/Masterminds/squirrel"
	"time"
)

func (r *UserRepository) GetListDocuments(ctx context.Context, input model.ListDocuments) (model.DocumentResponse, error) {
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
		Where(
			sq.Eq{"documents.user_id" : input.UserID}).
		OrderBy(input.Sort, "documents.id").
		PlaceholderFormat(sq.Dollar)
	sql, args, err := query.ToSql()
	if err != nil {
		return model.DocumentResponse{}, err
	}

	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		return model.DocumentResponse{}, err
	}
	defer rows.Close()

	docs := make(map[int64]*model.Document)
	order := make([]int64, 0) 

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
			return model.DocumentResponse{}, err
		}

		document, ok := docs[docID]
		if !ok {
			docs[docID] = &model.Document{
				UserID: uint64(userID),
				DocumentType: docType,
				DocumentNumber: docNumber,
				DateCreated: startDate,
				DateEnd: endDate,
			}

			order = append(order, docID)
			document = docs[docID]
		}

		if partyType != nil {
			document.Persons = append(document.Persons, model.Party{
				PartyType: *partyType,
				CompanyName: companyName,
				FirstName: firstName,
				MiddleName: middleName,
				LastName: lastName,
				Initials: initials,
			})
		}
	}

	if err := rows.Err(); err != nil {
		return model.DocumentResponse{}, err
	}

	total := len(order)
	if *input.Limit < len(order) {
		order = order[:*input.Limit]
	}

	result := make([]model.Document, 0, len(order))

	for _, id := range order {
		result = append(result, *docs[id])
	}

	return model.DocumentResponse{
		Total: int64(total),
		Documents: result,
	}, nil
} 
