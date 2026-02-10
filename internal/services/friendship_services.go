package services

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/kevinsofyan/echoes-chat-api/internal/models"
	"github.com/kevinsofyan/echoes-chat-api/internal/repositories"
)

type FriendshipService interface {
	SendFriendRequest(ctx context.Context, senderID, receiverID uuid.UUID) (*models.Friendship, error)
	AcceptFriendRequest(ctx context.Context, userID, friendshipID uuid.UUID) error
	RejectFriendRequest(ctx context.Context, userID, friendshipID uuid.UUID) error
	BlockUser(ctx context.Context, userID, targetUserID uuid.UUID) error
	UnblockUser(ctx context.Context, userID, targetUserID uuid.UUID) error
	UnfriendUser(ctx context.Context, userID, friendID uuid.UUID) error
	GetFriendsList(ctx context.Context, userID uuid.UUID, status string, limit, offset int) ([]FriendResponse, error)
}

type friendshipService struct {
	friendshipRepo repositories.FriendshipRepository
	userRepo       repositories.UserRepository
}

func NewFriendshipService(friendshipRepo repositories.FriendshipRepository, userRepo repositories.UserRepository) FriendshipService {
	return &friendshipService{
		friendshipRepo: friendshipRepo,
		userRepo:       userRepo,
	}
}

type FriendResponse struct {
	FriendshipID uuid.UUID `json:"friendship_id"`
	UserID       uuid.UUID `json:"user_id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	FullName     string    `json:"full_name"`
	Avatar       string    `json:"avatar"`
	IsOnline     bool      `json:"is_online"`
	Status       string    `json:"status"`
	CreatedAt    string    `json:"created_at"`
}

type FriendRequestResponse struct {
	ID        uuid.UUID      `json:"id"`
	Sender    FriendResponse `json:"sender"`
	Receiver  FriendResponse `json:"receiver"`
	Status    string         `json:"status"`
	CreatedAt string         `json:"created_at"`
}

func (s *friendshipService) SendFriendRequest(ctx context.Context, senderID, receiverID uuid.UUID) (*models.Friendship, error) {
	// Can't send request to yourself
	if senderID == receiverID {
		return nil, errors.New("cannot send friend request to yourself")
	}

	// Check if receiver exists
	_, err := s.userRepo.FindByID(ctx, receiverID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Check if friendship already exists
	existing, err := s.friendshipRepo.FindByUsers(ctx, senderID, receiverID)
	if err != nil {
		return nil, err
	}

	if existing != nil {
		switch existing.Status {
		case models.FriendshipPending:
			return nil, errors.New("friend request already sent")
		case models.FriendshipAccepted:
			return nil, errors.New("already friends")
		case models.FriendshipBlocked:
			return nil, errors.New("cannot send friend request")
		case models.FriendshipRejected:
			// Allow resending after rejection
			existing.Status = models.FriendshipPending
			existing.ActionUserID = senderID
			if err := s.friendshipRepo.Update(ctx, existing); err != nil {
				return nil, err
			}
			return existing, nil
		}
	}

	// Create new friend request
	friendship := &models.Friendship{
		UserID:       senderID,
		FriendID:     receiverID,
		Status:       models.FriendshipPending,
		ActionUserID: senderID,
	}

	if err := s.friendshipRepo.Create(ctx, friendship); err != nil {
		return nil, err
	}

	return friendship, nil
}

func (s *friendshipService) AcceptFriendRequest(ctx context.Context, userID, friendshipID uuid.UUID) error {
	friendship, err := s.friendshipRepo.FindByID(ctx, friendshipID)
	if err != nil {
		return err
	}

	// Only receiver can accept
	if friendship.FriendID != userID {
		return errors.New("unauthorized to accept this request")
	}

	if friendship.Status != models.FriendshipPending {
		return errors.New("friend request is not pending")
	}

	friendship.Status = models.FriendshipAccepted
	friendship.ActionUserID = userID

	return s.friendshipRepo.Update(ctx, friendship)
}

func (s *friendshipService) RejectFriendRequest(ctx context.Context, userID, friendshipID uuid.UUID) error {
	friendship, err := s.friendshipRepo.FindByID(ctx, friendshipID)
	if err != nil {
		return err
	}

	// Only receiver can reject
	if friendship.FriendID != userID {
		return errors.New("unauthorized to reject this request")
	}

	if friendship.Status != models.FriendshipPending {
		return errors.New("friend request is not pending")
	}

	friendship.Status = models.FriendshipRejected
	friendship.ActionUserID = userID

	return s.friendshipRepo.Update(ctx, friendship)
}

func (s *friendshipService) BlockUser(ctx context.Context, userID, targetUserID uuid.UUID) error {
	if userID == targetUserID {
		return errors.New("cannot block yourself")
	}

	friendship, err := s.friendshipRepo.FindByUsers(ctx, userID, targetUserID)
	if err != nil {
		return err
	}

	if friendship == nil {
		// Create new block entry
		friendship = &models.Friendship{
			UserID:       userID,
			FriendID:     targetUserID,
			Status:       models.FriendshipBlocked,
			ActionUserID: userID,
		}
		return s.friendshipRepo.Create(ctx, friendship)
	}

	friendship.Status = models.FriendshipBlocked
	friendship.ActionUserID = userID

	return s.friendshipRepo.Update(ctx, friendship)
}

func (s *friendshipService) UnblockUser(ctx context.Context, userID, targetUserID uuid.UUID) error {
	if userID == targetUserID {
		return errors.New("cannot unblock yourself")
	}

	friendship, err := s.friendshipRepo.FindByUsers(ctx, userID, targetUserID)
	if err != nil {
		return err
	}

	if friendship == nil {
		return errors.New("no relationship found")
	}

	if friendship.Status != models.FriendshipBlocked {
		return errors.New("user is not blocked")
	}

	// Remove the friendship record when unblocking
	return s.friendshipRepo.Delete(ctx, friendship.ID)
}

func (s *friendshipService) UnfriendUser(ctx context.Context, userID, friendID uuid.UUID) error {
	friendship, err := s.friendshipRepo.FindByUsers(ctx, userID, friendID)
	if err != nil {
		return err
	}

	if friendship == nil {
		return errors.New("friendship not found")
	}

	if friendship.Status != models.FriendshipAccepted {
		return errors.New("not friends")
	}

	return s.friendshipRepo.Delete(ctx, friendship.ID)
}

func (s *friendshipService) GetFriendsList(ctx context.Context, userID uuid.UUID, status string, limit, offset int) ([]FriendResponse, error) {
	if limit <= 0 {
		limit = 20 // Default limit
	}
	if offset < 0 {
		offset = 0
	}

	friendships, err := s.friendshipRepo.GetFriendsList(ctx, userID, status, limit, offset)
	if err != nil {
		return nil, err
	}

	friends := make([]FriendResponse, 0)
	for _, f := range friendships {
		var friend models.User
		var displayStatus string

		if f.UserID == userID {
			friend = f.Friend
			// User is the sender
			if f.Status == models.FriendshipPending {
				displayStatus = "sent" // You sent the request
			} else {
				displayStatus = string(f.Status)
			}
		} else {
			friend = f.User
			// User is the receiver
			displayStatus = string(f.Status) // Shows "pending" for received requests
		}

		friends = append(friends, FriendResponse{
			FriendshipID: f.ID,
			UserID:       friend.ID,
			Username:     friend.Username,
			FullName:     friend.FullName,
			Avatar:       friend.Avatar,
			IsOnline:     friend.IsOnline,
			Email:        friend.Email,
			Status:       displayStatus,
			CreatedAt:    f.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}

	return friends, nil
}
