package postgres

import (
	"comment-api/internal/domain"
	"fmt"
	"github.com/google/uuid"
)

func (r *CommentPostgres) Update(id uuid.UUID, list domain.UpdateCommentList) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	// 1. Проверяем, существует ли комментарий
	var exists bool
	checkQuery := "SELECT EXISTS (SELECT 1 FROM comment_lists WHERE id = $1)"
	err = tx.QueryRow(checkQuery, id).Scan(&exists)
	if err != nil {
		tx.Rollback()
		return err
	}
	if !exists {
		tx.Rollback()
		return fmt.Errorf("comment with id %s not found", id)
	}
	// 2. Обновляем комментарий
	updateQuery := `
		UPDATE comment_lists
		SET description = $1, updated_at = NOW()
		WHERE id = $2
	`
	_, err = tx.Exec(updateQuery, list.Description, id)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}
