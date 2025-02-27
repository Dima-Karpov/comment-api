package service

import (
	"comment-api/internal/domain"
	"comment-api/pkg/repository/postgres"
	"github.com/google/uuid"
)

type CommentListService struct {
	repo   postgres.CommentListPostgres
	filter *ProfaneFilterService
}

func NewCommentListService(repo postgres.CommentListPostgres, filter *ProfaneFilterService) *CommentListService {
	return &CommentListService{
		repo:   repo,
		filter: filter,
	}
}
func (s *CommentListService) Create(list domain.CommentList, traceID string) (uuid.UUID, error) {
	// Проверяем текст на запрезенные слова
	if err := s.filter.Check(list.Description, traceID); err != nil {
		return uuid.Nil, err
	}

	return s.repo.Create(list)
}
func (s *CommentListService) GetById(entityId uuid.UUID) ([]domain.Comment, error) {
	return s.repo.GetById(entityId)
}
func (s *CommentListService) Delete(id uuid.UUID) error {
	return s.repo.Delete(id)
}
func (s *CommentListService) Update(id uuid.UUID, list domain.UpdateCommentList, traceID string) error {
	// Проверяем текст на запрезенные слова
	if err := s.filter.Check(list.Description, traceID); err != nil {
		return err
	}

	return s.repo.Update(id, list)
}
