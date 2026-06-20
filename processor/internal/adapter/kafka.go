package adapter

import (
	"context"
	"encoding/json"
	"log"

	apperror "github.com/homka122/Gomka122/internal/errors"
	kafkaClient "github.com/homka122/Gomka122/internal/kafka"
	"github.com/homka122/Gomka122/processor/internal/domain"
)

type KafkaAdapter struct {
	producer kafkaClient.KafkaProducer
	consumer kafkaClient.KafkaConsumer
}

func NewKafkaAdapter(producer kafkaClient.KafkaProducer, consumer kafkaClient.KafkaConsumer) KafkaAdapter {
	return KafkaAdapter{producer: producer, consumer: consumer}
}

func (ka KafkaAdapter) SendTaskRequest(owner, repo string) error {
	req := kafkaClient.RepoTaskRequest{
		Owner: owner,
		Repo:  repo,
	}

	err := ka.producer.Send(context.Background(), "repo_request", req)
	if err != nil {
		return apperror.WrapCode(apperror.CodeInternal, "send kafka task requests", err)
	}
	log.Printf("sended task request %v", req)

	return nil
}

func (ka KafkaAdapter) RunGetTaskResponse(ctx context.Context, handle func(ctx context.Context, task domain.TaskResponse) error) error {
	return ka.consumer.Run(ctx, func(_ context.Context, data []byte) error {
		var task kafkaClient.RepoTaskResponse

		err := json.Unmarshal(data, &task)
		if err != nil {
			return apperror.Wrap("unmarshal kafka task", err)
		}

		log.Printf("get task for proceed for %v/%v", task.Owner, task.Repo)

		switch task.Status {
		case kafkaClient.STATUS_OK:
			return handle(ctx, domain.TaskResponse{
				Owner:       task.Owner,
				Repo:        task.Repo,
				Description: task.Description,
				Stars:       int32(task.Stars),
				Forks:       int32(task.Forks),
				CreateDate:  task.CreatedAt,
			})
		case kafkaClient.STATUS_NOT_FOUND:
			return apperror.New(kafkaClient.STATUS_NOT_FOUND, "repo not found")
		case kafkaClient.STATUS_INVALID_ARGUMENT:
			return apperror.New(apperror.CodeInvalidArgument, "invalid repo owner or name")
		case kafkaClient.STATUS_UNAVAILABLE:
			return apperror.New(apperror.CodeUnavailable, "service unavaliable")
		default:
			log.Printf("kafka adapter internal error %v", task)
			return apperror.New(apperror.CodeInternal, "internal error")
		}
	})
}
