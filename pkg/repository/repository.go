package repository

import (
	"comment-api/internal/domain"
	"github.com/google/uuid"
)

type CommentList interface {
	Create(list domain.CommentList) (uuid.UUID, error)
	GetById(id uuid.UUID) (domain.Comment, error)
	Delete(id uuid.UUID) error
	Update(id uuid.UUID, list domain.UpdateCommentList) error
}

type Repository struct {
	CommentList
}

func NewRepository(stage *LocalFile) *Repository {
	return &Repository{
		CommentList: NewCommentListStage(stage),
	}
}
