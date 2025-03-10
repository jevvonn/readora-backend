package main

import (
	"fmt"

	"github.com/hibiken/asynq"
	"github.com/jevvonn/readora-backend/config"
	"github.com/jevvonn/readora-backend/internal/infra/logger"
	"github.com/jevvonn/readora-backend/internal/infra/postgresql"
	"github.com/jevvonn/readora-backend/worker/tasks"
)

func main() {
	conf := config.New()
	log := logger.New()

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Jakarta",
		conf.DbHost,
		conf.DbPort,
		conf.DbUser,
		conf.DbPassword,
		conf.DbName,
	)
	db, err := postgresql.New(dsn)
	if err != nil {
		log.Error("[Worker][DB]", err)
		panic("could not connect to database")
	}

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

	// Mount Tasks
	mux.HandleFunc(tasks.SendOTPRegisterTaskName, tasks.HandleSendOTPRegisterTask)
	mux.HandleFunc(tasks.BooksFileUploadTaskName, tasks.HandleBooksFileUploadTask(db))
	mux.HandleFunc(tasks.BooksFileDeleteTaskName, tasks.HandleBooksFileDeleteTask)
	mux.HandleFunc(tasks.BooksFileProcessTaskName, tasks.HandleBooksFileParseTask(db))

	if err := srv.Run(mux); err != nil {
		log.Error("[Worker][Server]", err)
		panic("could not run worker server")
	}
}
