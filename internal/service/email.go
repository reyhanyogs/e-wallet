package service

import (
	"encoding/json"

	"github.com/reyhanyogs/e-wallet/domain"
	"github.com/reyhanyogs/e-wallet/dto"
	"github.com/reyhanyogs/e-wallet/internal/component"
)

type emailService struct {
	queueService domain.QueueService
}

func NewEmail(queueService domain.QueueService) domain.EmailService {
	return &emailService{
		queueService: queueService,
	}
}

func (s *emailService) Send(to, subject, body string) error {
	payload := dto.SendEmailReq{
		To:      to,
		Subject: subject,
		Body:    body,
	}
	data, err := json.Marshal(payload)
	if err != nil {
		component.Log.Errorf("Send(Marshal): to = %s :err = %s", payload.To, err.Error())
		return err
	}

	return s.queueService.Enqueue("send:email", data, 3)
}
