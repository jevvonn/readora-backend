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
	NewBooksFileUpload(tmpFile, fileName, booksId, fileType string) error
	NewBooksFileDelete(fileName string) error
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

	task := asynq.NewTask(tasks.SendOTPRegisterTaskName, payload, asynq.MaxRetry(3))
	_, err = w.client.Enqueue(task)

	return err
}

func (w *Worker) NewBooksFileUpload(tmpFile, fileName, booksId, fileType string) error {
	payload, err := json.Marshal(tasks.BooksFileUploadPayload{
		TmpFile:  tmpFile,
		BooksId:  booksId,
		Filename: fileName,
		FileType: fileType,
	})

	if err != nil {
		return err
	}

	task := asynq.NewTask(tasks.BooksFileUploadTaskName, payload, asynq.MaxRetry(3))
	w.client.Enqueue(task)

	w.NewBooksFileParsing(tmpFile, fileName, booksId)
	return err
}

func (w *Worker) NewBooksFileParsing(tmpFile, fileName, booksId string) error {
	payloadParse, err := json.Marshal(tasks.BooksFileParsePayload{
		TmpFile:  tmpFile,
		BooksId:  booksId,
		Filename: fileName,
	})

	if err != nil {
		return err
	}

	task := asynq.NewTask(tasks.BooksFileProcessTaskName, payloadParse, asynq.MaxRetry(3))
	w.client.Enqueue(task)

	return err
}

func (w *Worker) NewBooksFileDelete(fileName string) error {
	payload, err := json.Marshal(tasks.BooksFileUploadPayload{
		Filename: fileName,
	})

	if err != nil {
		return err
	}

	task := asynq.NewTask(tasks.BooksFileDeleteTaskName, payload, asynq.MaxRetry(3))
	_, err = w.client.Enqueue(task)

	return err
}
