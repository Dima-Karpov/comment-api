package domain

import (
	"github.com/google/uuid"
	"time"
)

type CommentList struct {
	Description string `json:"description" db:"description"`
	ParentID    string `json:"parent_id" db:"parent_id"`
}

type UpdateCommentList struct {
	Description string `json:"description" db:"description"`
}

type Comment struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Description string    `json:"description" db:"description"`
	ParentID    string    `json:"parent_id" db:"parent_id"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}
