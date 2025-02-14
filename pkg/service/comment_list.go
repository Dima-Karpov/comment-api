package service

import (
	"comment-api/internal/domain"
	"comment-api/pkg/repository"
	"github.com/google/uuid"
)

type CommentListService struct {
	repo repository.CommentList
}

func NewCommentListService(repo repository.CommentList) *CommentListService {
	return &CommentListService{repo: repo}
}
func (s *CommentListService) Create(list domain.CommentList) (uuid.UUID, error) {
	return s.repo.Create(list)
}
func (s *CommentListService) GetById(id uuid.UUID) (domain.Comment, error) {
	return s.repo.GetById(id)
}
func (s *CommentListService) Delete(id uuid.UUID) error {
	return s.repo.Delete(id)
}
func (s *CommentListService) Update(id uuid.UUID, list domain.UpdateCommentList) error {
	return s.repo.Update(id, list)
}
