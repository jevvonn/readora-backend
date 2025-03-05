package worker

import (
	"fmt"

	"github.com/hibiken/asynq"
	"github.com/jevvonn/readora-backend/config"
)

func NewWorkerClient() *asynq.Client {
	conf := config.Load()
	redisAddr := fmt.Sprintf("%s:%s", conf.RedisHost, conf.RedisPort)
	workerClient := asynq.NewClient(asynq.RedisClientOpt{
		Addr: redisAddr,
	})

	err := workerClient.Ping()
	if err != nil {
		panic("could not ping worker client")
	}

	return workerClient
}
