package repository

import (
	"context"
	"fmt"
	"ptmk-service/internal/model"

	sq "github.com/Masterminds/squirrel"
)

func (r *UserRepository) DeleteDocumentByNumber(ctx context.Context, documentNumber string) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	query := sq.
		Delete(Documents).
		Where(sq.Eq{"doc_number": documentNumber}).
		PlaceholderFormat(sq.Dollar)

	sqlStr, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("build delete document query: %w", err)
	}

	tag, err := tx.Exec(ctx, sqlStr, args...)
	if err != nil {
		return fmt.Errorf("delete document: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return model.ErrorDocumentNotFound
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}

	return nil
}
