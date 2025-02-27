package postgres

import (
	"github.com/jmoiron/sqlx"
)

type CommentPostgres struct {
	db *sqlx.DB
}

func NewCommentListPostgres(db *sqlx.DB) *CommentPostgres {
	return &CommentPostgres{db: db}
}
