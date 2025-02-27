package service

import (
	"comment-api/internal/domain"
	"comment-api/pkg/repository/postgres"
	"github.com/google/uuid"
)

type CommentList interface {
	Create(list domain.CommentList, traceID string) (uuid.UUID, error)
	GetById(entityId uuid.UUID) ([]domain.Comment, error)
	Delete(id uuid.UUID) error
	Update(id uuid.UUID, list domain.UpdateCommentList, traceID string) error
}

type Service struct {
	CommentList
}

func NewService(repos *postgres.RepositoryPostgres, filter *ProfaneFilterService) *Service {
	return &Service{
		CommentList: NewCommentListService(repos, filter),
	}
}
