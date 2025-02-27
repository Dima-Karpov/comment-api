package postgres

import (
	"comment-api/internal/domain"
	"fmt"
	"github.com/google/uuid"
	"time"
)

func (r *CommentPostgres) Create(input domain.CommentList) (uuid.UUID, error) {
	var comment domain.Comment
	// Парсим entity_id (он обязательный)
	entityID, err := uuid.Parse(input.EntityID)
	if err != nil {
		return uuid.Nil, err
	}
	// Обрабатываем parent_id (может быть пустым)
	var parentId *uuid.UUID
	if input.ParentID != "" {
		parsedUUID, err := uuid.Parse(input.ParentID)
		if err != nil {
			return uuid.Nil, err
		}
		var exists bool
		checkParentQuery := fmt.Sprintf("SELECT EXISTS (SELECT 1 FROM %s WHERE id = $1)", commentTable)
		err = r.db.QueryRow(checkParentQuery, parsedUUID).Scan(&exists)
		if err != nil {
			return uuid.Nil, err
		}
		if !exists {
			return uuid.Nil, fmt.Errorf("parent_id %s not found", parsedUUID.String())
		}
		parentId = &parsedUUID
	}

	comment.CreatedAt = time.Now()
	comment.UpdatedAt = time.Now()
	comment.Description = input.Description
	comment.ID = uuid.New()

	tx, err := r.db.Begin()
	if err != nil {
		return uuid.Nil, err
	}

	var id uuid.UUID
	// Вставляем данные в comment_lists
	createListQuery := fmt.Sprintf(
		"INSERT INTO %s (id, description, created_at, updated_at) VALUES ($1, $2, $3, $4) RETURNING id",
		commentTable,
	)
	row := tx.QueryRow(createListQuery, comment.ID, comment.Description, comment.CreatedAt, comment.UpdatedAt)
	if err := row.Scan(&id); err != nil {
		tx.Rollback()
		return uuid.Nil, err
	}

	// Вставляем данные в lists_items (с parent_id, который может быть NULL)
	createUsersListQuery := fmt.Sprintf(
		"INSERT INTO %s (entity_id, parent_id, list_id) VALUES ($1, $2, $3)", listsItemsTable,
	)
	_, err = tx.Exec(createUsersListQuery, entityID, parentId, comment.ID)
	if err != nil {
		tx.Rollback()
		return uuid.Nil, err
	}

	// Завершаем транзакцию
	if err := tx.Commit(); err != nil {
		return uuid.Nil, err
	}

	return comment.ID, nil
}
