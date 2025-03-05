package main

import (
	"fmt"

	"github.com/hibiken/asynq"
	"github.com/jevvonn/readora-backend/config"
	"github.com/jevvonn/readora-backend/internal/infra/logger"
	"github.com/jevvonn/readora-backend/worker/tasks"
)

func main() {
	conf := config.New()
	log := logger.New()

	redisAddr := fmt.Sprintf("%s:%s", conf.RedisHost, conf.RedisPort)
	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: redisAddr},
		asynq.Config{
			Concurrency: 4,
			Queues: map[string]int{
				"critical": 6,
				"default":  3,
				"low":      1,
			},
		},
	)

	if err := srv.Ping(); err != nil {
		log.Error("[Worker][Server]", err)
		panic("could not ping worker server")
	}

	mux := asynq.NewServeMux()

	mountTasks(mux)

	if err := srv.Run(mux); err != nil {
		log.Error("[Worker][Server]", err)
		panic("could not run worker server")
	}
}

func mountTasks(mux *asynq.ServeMux) {
	mux.HandleFunc(tasks.SendOTPRegisterTaskName, tasks.HandleSendOTPRegisterTask)
}
