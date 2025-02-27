package handler

import (
	"comment-api/internal/domain"
	"comment-api/pkg/service"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

// @Summary Create comment
// @Tags comments
// @Description create comment
// @ID create-comment
// @Accept  json
// @Produce  json
// @Param input body domain.CommentList true "list info"
// @Success 200 {object} domain.CreateCommentResponse
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /v1/comment [post]
func (h *Handler) createComment(c *gin.Context) {
	var input domain.CommentList

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid list param")
		return
	}

	traceID := c.GetHeader("X-Request-ID")

	// Проверяем, что entity_id не пустой
	if input.EntityID == "" {
		newErrorResponse(c, http.StatusForbidden, "entity_id cannot be empty")
		return
	}

	id, err := h.services.CommentList.Create(input, traceID)
	if err != nil {
		// Если это ошибка валидации (403)
		if errors.Is(err, service.ErrProfaneText) {
			c.JSON(http.StatusForbidden, gin.H{"message": err.Error()})
			return
		}

		// Если это другая ошибка, возвращаем 500
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

// @Summary Get Comment By EntityId
// @Tags comments
// @Description get comment by id
// @ID get-comment-by-entityId
// @Accept  json
// @Produce  json
// @Param entityId path string true "ID list (UUID)"
// @Success 200 {object} domain.Comment
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /v1/comment/{entityId} [get]
func (h *Handler) getComment(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	comment, err := h.services.CommentList.GetById(id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, comment)
}

// @Summary Delete comment
// @Tags comments
// @Description delete comment
// @ID delete-comment
// @Accept  json
// @Produce  json
// @Param commentId path string true "ID list (UUID)"
// @Success 200 {object} statusResponse
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /v1/comment/{commentId} [delete]
func (h *Handler) deleteComment(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	err = h.services.CommentList.Delete(id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{"ok"})
}

// @Summary Update comment
// @Tags comments
// @Description comment post
// @ID comment-post
// @Accept  json
// @Produce  json
// @Param commentId path string true "ID list (UUID)"
// @Param data body domain.UpdateCommentList true "Data for list"
// @Success 200 {object} statusResponse
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /v1/comment/{commentId} [put]
func (h *Handler) updateComment(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}
	var input domain.UpdateCommentList
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid list param")
		return
	}
	traceID := c.GetHeader("X-Request-ID")

	err = h.services.CommentList.Update(id, input, traceID)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{"ok"})

}
