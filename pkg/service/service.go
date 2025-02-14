package service

import (
	"comment-api/internal/domain"
	"comment-api/pkg/repository"
	"github.com/google/uuid"
)

type CommentList interface {
	Create(list domain.CommentList) (uuid.UUID, error)
	GetById(id uuid.UUID) (domain.Comment, error)
	Delete(id uuid.UUID) error
	Update(id uuid.UUID, list domain.UpdateCommentList) error
}

type Service struct {
	CommentList
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		CommentList: NewCommentListService(repos),
	}
}
