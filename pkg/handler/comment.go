package handler

import (
	"comment-api/internal/domain"
	"comment-api/pkg/service"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

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
