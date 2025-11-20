package handlers

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/kevinsofyan/echoes-chat-api/internal/services"
	"github.com/kevinsofyan/echoes-chat-api/internal/utils"
	"github.com/labstack/echo/v4"
)

type RoomHandler struct {
	roomService services.RoomService
}

func NewRoomHandler(roomService services.RoomService) *RoomHandler {
	return &RoomHandler{
		roomService: roomService,
	}
}

// CreateRoom godoc
// @Summary Create a new room
// @Tags rooms
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body services.CreateRoomRequest true "Room data"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /api/v1/rooms [post]
func (h *RoomHandler) CreateRoom(c echo.Context) error {
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"error": "Unauthorized",
		})
	}

	var req services.CreateRoomRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid request body",
		})
	}

	req.CreatedBy = userID

	room, err := h.roomService.CreateRoom(c.Request().Context(), req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "Room created successfully",
		"data":    room,
	})
}

// GetRoomByID godoc
// @Summary Get room by ID
// @Tags rooms
// @Security BearerAuth
// @Produce json
// @Param id path string true "Room UUID"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/v1/rooms/{id} [get]
func (h *RoomHandler) GetRoomByID(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid room ID",
		})
	}

	room, err := h.roomService.GetRoomByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"error": "Room not found",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": room,
	})
}

// GetMyRooms godoc
// @Summary Get rooms for authenticated user
// @Tags rooms
// @Security BearerAuth
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /api/v1/rooms/my [get]
func (h *RoomHandler) GetMyRooms(c echo.Context) error {
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"error": "Unauthorized",
		})
	}

	rooms, err := h.roomService.GetUserRooms(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": rooms,
	})
}
