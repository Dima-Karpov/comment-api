package postgres

import (
	"database/sql"
	"fmt"
	"github.com/google/uuid"
)

func (r *CommentPostgres) Delete(id uuid.UUID) error {
	// Начинаем транзакцию
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

	// 2. Рекурсивно удаляем все дочерние комментарии
	err = deleteChildren(tx, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	// 3. Удаляем комментарий
	deleteQuery := "DELETE FROM comment_lists WHERE id = $1"
	_, err = tx.Exec(deleteQuery, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	// 4. Завершаем транзакцию
	return tx.Commit()
}

// Рекурсивное удаление дочерних комментариев
func deleteChildren(tx *sql.Tx, parentId uuid.UUID) error {
	// Получаем всех детей
	getChildrenQuery := "SELECT list_id FROM lists_items WHERE parent_id = $1"
	rows, err := tx.Query(getChildrenQuery, parentId)
	if err != nil {
		return err
	}
	defer rows.Close()

	var childIds []uuid.UUID
	for rows.Next() {
		var childId uuid.UUID
		if err := rows.Scan(&childId); err != nil {
			return err
		}
		childIds = append(childIds, childId)
	}

	// Удаляем каждого ребенка
	for _, childId := range childIds {
		err = deleteChildren(tx, childId)
		if err != nil {
			return err
		}
		_, err = tx.Exec("DELETE FROM comment_lists WHERE id = $1", childId)
		if err != nil {
			return err
		}
	}
	return nil
}
