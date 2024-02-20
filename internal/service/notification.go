package service

import (
	"bytes"
	"context"
	"html/template"
	"time"

	"github.com/reyhanyogs/e-wallet/domain"
	"github.com/reyhanyogs/e-wallet/dto"
	"github.com/reyhanyogs/e-wallet/internal/component"
)

type notificationService struct {
	notificationRepository domain.NotificationRepository
	templateRepository     domain.TemplateRepository
	hub                    *dto.Hub
}

func NewNotification(notificationRepository domain.NotificationRepository, templateRepository domain.TemplateRepository, hub *dto.Hub) domain.NotificationService {
	return &notificationService{
		notificationRepository: notificationRepository,
		templateRepository:     templateRepository,
		hub:                    hub,
	}
}

func (s *notificationService) FindByUser(ctx context.Context, user int64) ([]dto.NotificationData, error) {
	notifications, err := s.notificationRepository.FindByUser(ctx, user)
	if err != nil {
		component.Log.Errorf("FindByUser(FindByUser): user_id = %d: err = %s", user, err.Error())
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

	return result, nil
}

func (s *notificationService) Insert(ctx context.Context, userId int64, code string, data map[string]string) error {
	tmpl, err := s.templateRepository.FindByCode(ctx, code)
	if err != nil {
		component.Log.Errorf("Insert(FindByCode): user_id = %d: err = %s", userId, err.Error())
		return err
	}
	if tmpl == (domain.Template{}) {
		return domain.ErrTemplateNotFound
	}

	body := new(bytes.Buffer)
	t := template.Must(template.New("notif").Parse(tmpl.Body))
	err = t.Execute(body, data)
	if err != nil {
		component.Log.Errorf("Insert(Execute): user_id = %d: err = %s", userId, err.Error())
		return err
	}

	notification := domain.Notification{
		UserID:    userId,
		Title:     tmpl.Title,
		Body:      body.String(),
		Status:    1,
		IsRead:    0,
		CreatedAt: time.Now(),
	}
	err = s.notificationRepository.Insert(ctx, &notification)
	if err != nil {
		component.Log.Errorf("Insert(Insert): user_id = %d: err = %s", userId, err.Error())
		return err
	}

	if channel, ok := s.hub.NotificationChannel[userId]; ok {
		channel <- dto.NotificationData{
			ID:        notification.ID,
			Title:     notification.Title,
			Body:      notification.Body,
			Status:    notification.Status,
			IsRead:    notification.IsRead,
			CreatedAt: notification.CreatedAt,
		}
	}

	return nil
}
