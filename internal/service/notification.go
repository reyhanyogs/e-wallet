package service

import (
	"context"

	"github.com/reyhanyogs/e-wallet/domain"
	"github.com/reyhanyogs/e-wallet/dto"
)

type notificationService struct {
	notificationRepository domain.NotificationRepository
}

func NewNotification(notificationRepository domain.NotificationRepository) domain.NotificationService {
	return &notificationService{
		notificationRepository: notificationRepository,
	}
}

func (s *notificationService) FindByUser(ctx context.Context, user int64) ([]dto.NotificationData, error) {
	notifications, err := s.notificationRepository.FindByUser(ctx, user)
	if err != nil {
		return nil, err
	}

	var result []dto.NotificationData
	for _, v := range notifications {
		result = append(result, dto.NotificationData{
			ID:        v.ID,
			Title:     v.Title,
			Body:      v.Body,
			Status:    v.Status,
			IsRead:    v.IsRead,
			CreatedAt: v.CreatedAt,
		})
	}
	if result == nil {
		result = make([]dto.NotificationData, 0)
	}

	return result, err
}
