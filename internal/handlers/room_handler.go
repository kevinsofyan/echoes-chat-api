package handlers

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/kevinsofyan/echoes-chat-api/internal/models"
	"github.com/kevinsofyan/echoes-chat-api/internal/services"
	"github.com/kevinsofyan/echoes-chat-api/internal/utils"
	"github.com/labstack/echo/v4"
)

type RoomHandler struct {
	roomService services.RoomService
}

type RoomResponse struct {
	ID          uuid.UUID    `json:"id"`
	Name        string       `json:"name"`
	Type        string       `json:"type"`
	Description string       `json:"description"`
	Avatar      string       `json:"avatar"`
	CreatedBy   uuid.UUID    `json:"created_by"`
	Members     []MemberInfo `json:"members,omitempty"`
	CreatedAt   interface{}  `json:"created_at"`
	UpdatedAt   interface{}  `json:"updated_at"`
}

type MemberInfo struct {
	UserID   uuid.UUID `json:"user_id"`
	Username string    `json:"username"`
	FullName string    `json:"full_name"`
	Avatar   string    `json:"avatar"`
	Role     string    `json:"role"`
	IsOnline bool      `json:"is_online"`
}

func NewRoomHandler(roomService services.RoomService) *RoomHandler {
	return &RoomHandler{
		roomService: roomService,
	}
}

func formatRoomResponse(room *models.Room, currentUserID uuid.UUID) RoomResponse {
	response := RoomResponse{
		ID:          room.ID,
		Name:        room.Name,
		Type:        string(room.Type),
		Description: room.Description,
		Avatar:      room.Avatar,
		CreatedBy:   room.CreatedBy,
		CreatedAt:   room.CreatedAt,
		UpdatedAt:   room.UpdatedAt,
	}

	// Format members info
	members := make([]MemberInfo, 0, len(room.Members))
	for _, member := range room.Members {
		members = append(members, MemberInfo{
			UserID:   member.User.ID,
			Username: member.User.Username,
			FullName: member.User.FullName,
			Avatar:   member.User.Avatar,
			Role:     string(member.Role),
			IsOnline: member.User.IsOnline,
		})
	}
	response.Members = members

	// For direct chats, use the other user's info
	if room.Type == models.RoomTypeDirect && len(room.Members) == 2 {
		for _, member := range room.Members {
			if member.UserID != currentUserID {
				response.Name = member.User.FullName
				if response.Name == "" {
					response.Name = member.User.Username
				}
				response.Avatar = member.User.Avatar
				response.Description = member.User.Email
				break
			}
		}
	}

	return response
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
		"data":    formatRoomResponse(room, userID),
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
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"error": "Unauthorized",
		})
	}

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
		"data": formatRoomResponse(room, userID),
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

	// Format each room response
	formattedRooms := make([]RoomResponse, len(rooms))
	for i, room := range rooms {
		formattedRooms[i] = formatRoomResponse(&room, userID)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": formattedRooms,
	})
}

// DeleteRoom godoc
// @Summary Delete a room
// @Tags rooms
// @Security BearerAuth
// @Produce json
// @Param id path string true "Room UUID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/v1/rooms/{id} [delete]
func (h *RoomHandler) DeleteRoom(c echo.Context) error {
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"error": "Unauthorized",
		})
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid room ID",
		})
	}

	// Check if user is the room creator
	room, err := h.roomService.GetRoomByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"error": "Room not found",
		})
	}

	if room.CreatedBy != userID {
		return c.JSON(http.StatusForbidden, map[string]interface{}{
			"error": "Only room creator can delete the room",
		})
	}

	if err := h.roomService.DeleteRoom(c.Request().Context(), id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Room deleted successfully",
	})
}

// GetOrCreateDirectChat godoc
// @Summary Get or create direct chat with another user
// @Tags rooms
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param user_id path string true "Other User UUID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /api/v1/rooms/direct/{user_id} [post]
func (h *RoomHandler) GetOrCreateDirectChat(c echo.Context) error {
	currentUserID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"error": "Unauthorized",
		})
	}

	otherUserID, err := uuid.Parse(c.Param("user_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid user ID",
		})
	}

	if currentUserID == otherUserID {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Cannot create direct chat with yourself",
		})
	}

	room, err := h.roomService.GetOrCreateDirectRoom(c.Request().Context(), currentUserID, otherUserID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Direct chat ready",
		"data":    formatRoomResponse(room, currentUserID),
	})
}
