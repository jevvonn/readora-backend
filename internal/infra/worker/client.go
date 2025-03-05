package worker

import (
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
	"github.com/jevvonn/readora-backend/config"
	"github.com/jevvonn/readora-backend/worker/tasks"
)

type WorkerItf interface {
	NewSendOTPRegisterTask(email string, otp string) error
}

type Worker struct {
	client *asynq.Client
}

func NewWorkerClient() WorkerItf {
	conf := config.Load()
	redisAddr := fmt.Sprintf("%s:%s", conf.RedisHost, conf.RedisPort)
	workerClient := asynq.NewClient(asynq.RedisClientOpt{
		Addr: redisAddr,
	})

	err := workerClient.Ping()
	if err != nil {
		panic("could not ping worker client")
	}

	return &Worker{
		client: workerClient,
	}
}

func (w *Worker) NewSendOTPRegisterTask(email string, otp string) error {
	payload, err := json.Marshal(tasks.SendOTPRegisterPayload{
		Email: email,
		OTP:   otp,
	})

	if err != nil {
		return err
	}

	task := asynq.NewTask(tasks.SendOTPRegisterTaskName, payload)
	_, err = w.client.Enqueue(task)

	return err
}
