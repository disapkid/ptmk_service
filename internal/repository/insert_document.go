package repository

import (
	"context"
	"ptmk-service/internal/model"
	"fmt"
	sq "github.com/Masterminds/squirrel"
)

func (r *UserRepository) InsertDocument(ctx context.Context, doc model.Document) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	var docID int64

	query := sq.
		Insert(Documents).
		Columns(DocColumns...). 
		Values(
			doc.UserID,
			doc.DocumentType,
			doc.DateCreated,
			doc.DateEnd,
			doc.DocumentNumber,
		).Suffix("RETURNING id").
		PlaceholderFormat(sq.Dollar)

	sqlStr, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("Build insert document error: %w", err)
	}

	if err = tx.QueryRow(ctx, sqlStr, args...).Scan(&docID); err != nil {
		return fmt.Errorf("Insert document error: %w", err)
	}

	if len(doc.Persons) > 0 {
		partieQuery := sq.
			Insert(DocParties).
			Columns(PartiesColumns...).
			PlaceholderFormat(sq.Dollar)

		for _, val := range doc.Persons {
			partieQuery = partieQuery.Values(
				docID,
				val.CompanyName,
				val.FirstName,
				val.LastName,
				val.Initials,
				val.PartyType,
				val.MiddleName,
			)
		}

		sqlStr, args, err := partieQuery.ToSql()
		if err != nil {
			return fmt.Errorf("Build insert parties query error: %w", err)
		}

		if _,err = tx.Exec(ctx, sqlStr, args...); err != nil {
			return fmt.Errorf("Insert parties error: %w", err)
		}
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("Commit tx error: %w", err)
	}

	return nil
}
