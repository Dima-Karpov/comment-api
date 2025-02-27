package postgres

import (
	"comment-api/internal/domain"
	"database/sql/driver"
	"github.com/google/uuid"
	"sort"
)

type NullUUID struct {
	UUID  uuid.UUID
	Valid bool
}

func (n *NullUUID) Scan(value interface{}) error {
	if value == nil {
		n.Valid = false
		return nil
	}
	n.Valid = true
	return n.UUID.Scan(value)
}
func (n NullUUID) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.UUID, nil
}

func (r *CommentPostgres) GetById(entityId uuid.UUID) ([]domain.Comment, error) {
	commentMap := make(map[uuid.UUID]*domain.Comment) // Мапа для быстрого доступа по ID
	var allComments []*domain.Comment                 // Храним ссылки, чтобы не терять детей

	// 1. Получаем все комментарии по entity_id, сортируем по created_at
	getCommentsQuery := `
		SELECT cl.id, cl.description, cl.created_at, cl.updated_at, li.parent_id
		FROM comment_lists cl
		JOIN lists_items li ON cl.id = li.list_id
		WHERE li.entity_id = $1
		ORDER BY cl.created_at ASC
	`

	rows, err := r.db.Query(getCommentsQuery, entityId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// 2. Обрабатываем результаты запроса
	for rows.Next() {
		var comment domain.Comment
		var sqlParentId NullUUID

		err := rows.Scan(
			&comment.ID,
			&comment.Description,
			&comment.CreatedAt,
			&comment.UpdatedAt,
			&sqlParentId,
		)
		if err != nil {
			return nil, err
		}

		// Если parent_id НЕ NULL, то сохраняем его
		if sqlParentId.Valid {
			comment.ParentID = sqlParentId.UUID.String()
		} else {
			comment.ParentID = ""
		}

		// Добавляем комментарий в мапу и массив
		commentMap[comment.ID] = &comment
		allComments = append(allComments, &comment)
	}

	// 3. Строим дерево комментариев
	var rootComments []*domain.Comment
	for _, comment := range allComments {
		if comment.ParentID == "" {
			// Корневой комментарий
			rootComments = append(rootComments, comment)
		} else {
			// Дочерний комментарий
			parentUUID, err := uuid.Parse(comment.ParentID)
			if err == nil {
				if parent, exists := commentMap[parentUUID]; exists {
					parent.Children = append(parent.Children, *comment)
				}
			}
		}
	}

	// 4. Преобразуем дерево и сортируем детей
	result := make([]domain.Comment, len(rootComments))
	for i, root := range rootComments {
		result[i] = *root
		sortChildren(&result[i]) // Сортируем вложенные комментарии
	}

	return result, nil
}

func sortChildren(comment *domain.Comment) {
	if len(comment.Children) > 0 {
		sort.SliceStable(comment.Children, func(i, j int) bool {
			return comment.Children[i].CreatedAt.Before(comment.Children[j].CreatedAt)
		})
		// Рекурсивно сортируем детей
		for i := range comment.Children {
			sortChildren(&comment.Children[i])
		}
	}
}
