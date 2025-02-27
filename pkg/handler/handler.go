package handler

import (
	"comment-api/pkg/middleware"
	"comment-api/pkg/service"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"

	_ "comment-api/docs"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	logger := logrus.New()
	router.Use(middleware.RequestLoggerMiddleware(logger))

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	comment := router.Group("/v1/comment")
	{
		comment.POST("/", h.createComment)
		comment.GET("/:id", h.getComment)
		comment.DELETE("/:id", h.deleteComment)
		comment.PUT("/:id", h.updateComment)
	}

	return router
}
