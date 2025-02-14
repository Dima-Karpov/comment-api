package handler

import (
	"comment-api/pkg/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()
	comment := router.Group("/comment")
	{
		comment.POST("/", h.createComment)
		comment.GET("/:id", h.getComment)
		comment.DELETE("/:id", h.deleteComment)
		comment.PUT("/:id", h.updateComment)
	}

	return router
}
