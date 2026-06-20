package kafkaController

import (
	"context"
	"encoding/json"
	"log"

	"github.com/homka122/Gomka122/collector/internal/domain"
	apperror "github.com/homka122/Gomka122/internal/errors"
	kafkaClient "github.com/homka122/Gomka122/internal/kafka"
)

type Tasker interface {
	ProcessTask(domain.Task) error
}

type KafkaController struct {
	tasker   Tasker
	consumer kafkaClient.KafkaConsumer
}

func NewKafkaController(tasker Tasker, consumer kafkaClient.KafkaConsumer) KafkaController {
	return KafkaController{
		tasker:   tasker,
		consumer: consumer,
	}
}

func (kc KafkaController) Run() error {
	return kc.consumer.Run(context.Background(), func(_ context.Context, data []byte) error {
		var task kafkaClient.RepoTaskRequest

		err := json.Unmarshal(data, &task)
		if err != nil {
			return apperror.Wrap("unmarshal kafka task", err)
		}

		log.Printf("get task for proceed for %v/%v", task.Owner, task.Repo)

		return kc.tasker.ProcessTask(domain.Task{
			Owner: task.Owner,
			Repo:  task.Repo,
		})
	})
}
