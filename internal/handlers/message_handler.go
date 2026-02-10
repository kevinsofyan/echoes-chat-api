package handlers

import (
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/kevinsofyan/echoes-chat-api/internal/services"
	"github.com/labstack/echo/v4"
)

type MessageHandler struct {
	messageService services.MessageService
}

func NewMessageHandler(messageService services.MessageService) *MessageHandler {
	return &MessageHandler{
		messageService: messageService,
	}
}

// GetRoomMessages godoc
// @Summary Get messages for a room
// @Tags messages
// @Security BearerAuth
// @Produce json
// @Param room_id path string true "Room UUID"
// @Param limit query int false "Limit" default(50)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /api/v1/rooms/{room_id}/messages [get]
func (h *MessageHandler) GetRoomMessages(c echo.Context) error {
	roomID, err := uuid.Parse(c.Param("room_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid room ID",
		})
	}

	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	offset, _ := strconv.Atoi(c.QueryParam("offset"))

	if limit == 0 {
		limit = 50
	}

	messages, err := h.messageService.GetMessagesByRoomID(c.Request().Context(), roomID, limit, offset)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": messages,
		"pagination": map[string]interface{}{
			"limit":  limit,
			"offset": offset,
			"count":  len(messages),
		},
	})
}
