package postgres

import (
	"comment-api/internal/domain"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type CommentListPostgres interface {
	Create(list domain.CommentList) (uuid.UUID, error)
	GetById(entityId uuid.UUID) ([]domain.Comment, error)
	Delete(id uuid.UUID) error
	Update(id uuid.UUID, list domain.UpdateCommentList) error
}

type RepositoryPostgres struct {
	CommentListPostgres
}

func NewRepositoryPostgres(db *sqlx.DB) *RepositoryPostgres {
	return &RepositoryPostgres{
		CommentListPostgres: NewCommentListPostgres(db),
	}
}
