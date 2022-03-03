package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/todo"
	"net/http"
)

type getLatestTodoItemResponse struct {
	Data    todo.TodoItem `json:"data"`
	Message *string       `json:"message"`
}

// @Summary Get latest TodoItem
// @Tags TodoItem
// @Description get latest TodoItem object
// @ID todo-latest
// @Accept  json
// @Produce  json
// @Success 200 {object} getLatestTodoItemResponse
// @Failure 404 {object} responseBag
// @Failure 500 {object} responseBag
// @Failure default {object} responseBag
// @Router /todo [get]
func (h *Handler) getLatestTodoItem(ctx *gin.Context) {
	item, err := h.services.TodoItem.GetLatest()

	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	if item == nil {
		newErrorResponse(ctx, http.StatusNotFound, "No entities found")
		return
	}

	newStatusResponse(ctx, http.StatusOK, item)
}
