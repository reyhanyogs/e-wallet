package service

import (
	"github.com/hibiken/asynq"
	"github.com/reyhanyogs/e-wallet/domain"
	"github.com/reyhanyogs/e-wallet/internal/component"
	"github.com/reyhanyogs/e-wallet/internal/config"
)

type queueService struct {
	queueClient *asynq.Client
}

func NewQueue(config *config.Config) domain.QueueService {
	redisConnection := asynq.RedisClientOpt{
		Addr:     config.Queue.Addr,
		Password: config.Queue.Pass,
	}
	return &queueService{
		queueClient: asynq.NewClient(redisConnection),
	}
}

func (s *queueService) Enqueue(name string, data []byte, retry int) error {
	task := asynq.NewTask(name, data, asynq.MaxRetry(retry))

	info, err := s.queueClient.Enqueue(task)
	if err != nil {
		component.Log.Errorf("Enqueue; %s", err.Error())
	}
	component.Log.Info("Enqueue; ", info.Payload)
	return nil
}
