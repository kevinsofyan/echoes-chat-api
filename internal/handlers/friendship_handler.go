package handlers

import (
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/kevinsofyan/echoes-chat-api/internal/services"
	"github.com/kevinsofyan/echoes-chat-api/internal/utils"
	"github.com/labstack/echo/v4"
)

type FriendshipHandler struct {
	friendshipService services.FriendshipService
}

func NewFriendshipHandler(friendshipService services.FriendshipService) *FriendshipHandler {
	return &FriendshipHandler{
		friendshipService: friendshipService,
	}
}

// SendFriendRequest godoc
// @Summary Send friend request
// @Tags friendships
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param user_id path string true "User ID to send request to"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /api/v1/friends/request/{user_id} [post]
func (h *FriendshipHandler) SendFriendRequest(c echo.Context) error {
	senderID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"error": "Unauthorized",
		})
	}

	receiverID, err := uuid.Parse(c.Param("user_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid user ID",
		})
	}

	friendship, err := h.friendshipService.SendFriendRequest(c.Request().Context(), senderID, receiverID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Friend request sent successfully",
		"data": map[string]interface{}{
			"id":             friendship.ID,
			"user_id":        friendship.UserID,
			"friend_id":      friendship.FriendID,
			"status":         friendship.Status,
			"action_user_id": friendship.ActionUserID,
			"created_at":     friendship.CreatedAt,
			"updated_at":     friendship.UpdatedAt,
		},
	})
}

// ManageFriendRequest godoc
// @Summary Manage friend request (accept/reject)
// @Tags friendships
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Friendship ID"
// @Param body body object true "Action body" example({"action": "accept"})
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /api/v1/friends/request/{id} [put]
func (h *FriendshipHandler) ManageFriendRequest(c echo.Context) error {
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"error": "Unauthorized",
		})
	}

	friendshipID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid friendship ID",
		})
	}

	var req struct {
		Action string `json:"action"` // "accept" or "reject"
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid request body",
		})
	}

	var message string
	switch req.Action {
	case "accept":
		if err := h.friendshipService.AcceptFriendRequest(c.Request().Context(), userID, friendshipID); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"error": err.Error(),
			})
		}
		message = "Friend request accepted"
	case "reject":
		if err := h.friendshipService.RejectFriendRequest(c.Request().Context(), userID, friendshipID); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"error": err.Error(),
			})
		}
		message = "Friend request rejected"
	default:
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid action. Use 'accept' or 'reject'",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": message,
	})
}

// ManageBlockUser godoc
// @Summary Manage user blocking (block/unblock)
// @Tags friendships
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param user_id path string true "User ID to block/unblock"
// @Param body body object true "Action body" example({"action": "block"})
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /api/v1/friends/block/{user_id} [put]
func (h *FriendshipHandler) ManageBlockUser(c echo.Context) error {
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"error": "Unauthorized",
		})
	}

	targetUserID, err := uuid.Parse(c.Param("user_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid user ID",
		})
	}

	var req struct {
		Action string `json:"action"` // "block" or "unblock"
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid request body",
		})
	}

	var message string
	switch req.Action {
	case "block":
		if err := h.friendshipService.BlockUser(c.Request().Context(), userID, targetUserID); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"error": err.Error(),
			})
		}
		message = "User blocked successfully"
	case "unblock":
		if err := h.friendshipService.UnblockUser(c.Request().Context(), userID, targetUserID); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"error": err.Error(),
			})
		}
		message = "User unblocked successfully"
	default:
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid action. Use 'block' or 'unblock'",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": message,
	})
}

// Unfriend godoc
// @Summary Unfriend a user
// @Tags friendships
// @Security BearerAuth
// @Produce json
// @Param user_id path string true "Friend user ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /api/v1/friends/{user_id} [delete]
func (h *FriendshipHandler) Unfriend(c echo.Context) error {
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"error": "Unauthorized",
		})
	}

	friendID, err := uuid.Parse(c.Param("user_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid user ID",
		})
	}

	if err := h.friendshipService.UnfriendUser(c.Request().Context(), userID, friendID); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Unfriended successfully",
	})
}

// GetFriendsList godoc
// @Summary Get friends list with optional status filter
// @Tags friendships
// @Security BearerAuth
// @Produce json
// @Param status query string false "Status filter: pending, sent, accepted (leave empty for all)"
// @Param limit query int false "Limit" default(20)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /api/v1/friends [get]
func (h *FriendshipHandler) GetFriendsList(c echo.Context) error {
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"error": "Unauthorized",
		})
	}

	status := c.QueryParam("status")
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	offset, _ := strconv.Atoi(c.QueryParam("offset"))

	friends, err := h.friendshipService.GetFriendsList(c.Request().Context(), userID, status, limit, offset)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": friends,
		"pagination": map[string]interface{}{
			"limit":  limit,
			"offset": offset,
			"count":  len(friends),
		},
	})
}
