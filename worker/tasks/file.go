package tasks

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"github.com/jevvonn/readora-backend/config"
	"github.com/jevvonn/readora-backend/internal/constant"
	"github.com/jevvonn/readora-backend/internal/domain/entity"
	"github.com/jevvonn/readora-backend/internal/infra/logger"
	"github.com/jevvonn/readora-backend/internal/infra/storage"
	"gorm.io/gorm"
)

const (
	BooksFileUploadTaskName  = "file-upload:books"
	BooksFileDeleteTaskName  = "file-delete:books"
	BooksFileProcessTaskName = "file-process:books"
)

type BooksFileUploadPayload struct {
	TmpFile  string
	Filename string
	BooksId  string
	FileType string
}

type BooksFileParsePayload struct {
	TmpFile  string
	Filename string
	BooksId  string
}

type BooksFileProcessingResponse struct {
	Message        string `json:"message"`
	FilePublicURL  string `json:"filePublicUrl"`
	CoverPublicURL string `json:"coverPublicUrl"`
}

type BooksFileDeletePayload struct {
	Filename string
}

// Create Task
func HandleBooksFileUploadTask(db *gorm.DB) asynq.HandlerFunc {
	return func(ctx context.Context, t *asynq.Task) error {
		var payload BooksFileUploadPayload
		if err := json.Unmarshal(t.Payload(), &payload); err != nil {
			return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
		}

		log := logger.New()
		err := db.Model(&entity.Book{
			ID: uuid.MustParse(payload.BooksId),
		}).Updates(entity.Book{
			FileUploadStatus: constant.BookFileUploadStatusUploading,
		}).Error
		if err != nil {
			log.Error("[Task][BooksFileUpload]", err)
			return err
		}

		// Check path
		_, err = os.Stat(payload.TmpFile)
		if os.IsNotExist(err) {
			log.Error("[Task][BooksFileUpload]", err)
			return err
		}

		// Upload to Supabase
		file, err := os.Open(payload.TmpFile)
		if err != nil {
			return err
		}

		storage := storage.New()
		publicUrl, err := storage.UploadFile(file, "books", payload.Filename, payload.FileType)
		if err != nil {
			log.Error("[Task][BooksFileUpload]", err)
			return err
		}

		// Update file_url in books
		err = db.Model(&entity.Book{
			ID: uuid.MustParse(payload.BooksId),
		}).Updates(entity.Book{
			FileURL:          publicUrl,
			FileUploadStatus: constant.BookFileUploadStatusUploaded,
		}).Error
		if err != nil {
			log.Error("[Task][BooksFileUpload]", err)
			return err
		}

		file.Close()
		defer os.Remove(payload.TmpFile)

		fmt.Println("[Task][BooksFileUpload] Run Succesfully")
		return nil
	}
}

func HandleBooksFileParseTask(db *gorm.DB) asynq.HandlerFunc {
	return func(ctx context.Context, t *asynq.Task) error {
		var payload BooksFileParsePayload
		if err := json.Unmarshal(t.Payload(), &payload); err != nil {
			return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
		}

		log := logger.New()
		conf := config.Load()
		err := db.Model(&entity.Book{
			ID: uuid.MustParse(payload.BooksId),
		}).Updates(entity.Book{
			FileAIStatus: constant.BookFileAIStatusProcessing,
		}).Error
		if err != nil {
			log.Error("[Task][BooksFileParsePayload]", err)
			return err
		}

		// Check path
		_, err = os.Stat(payload.TmpFile)
		if os.IsNotExist(err) {
			log.Error("[Task][BooksFileParsePayload]", err)
			return err
		}

		// Upload to Supabase
		file, err := os.Open(payload.TmpFile)
		if err != nil {
			return err
		}

		client := resty.New()
		url := fmt.Sprintf("%s/service/books/parse", conf.NodeApiBaseURL)
		resp, err := client.R().
			SetFile("file", payload.TmpFile).
			SetFormData(map[string]string{
				"bookId": payload.BooksId,
			}).
			Post(url)

		if err != nil {
			log.Error("Error:", err)
			return err
		}

		if resp.StatusCode() != 200 {
			return errors.New("failed to parse file")
		}

		file.Close()
		defer os.Remove(payload.TmpFile)

		var response BooksFileProcessingResponse
		err = json.Unmarshal(resp.Body(), &response)
		if err != nil {
			log.Error("[Task][BooksFileParsePayload]", err)
			return err
		}

		err = db.Model(&entity.Book{
			ID: uuid.MustParse(payload.BooksId),
		}).Updates(entity.Book{
			CoverImageURL: response.CoverPublicURL,
			FileAIStatus:  constant.BookFileAIStatusReady,
		}).Error
		if err != nil {
			log.Error("[Task][BooksFileParsePayload]", err)
			return err
		}

		fmt.Println("[Task][BooksFileParsePayload] Run Succesfully")
		return nil
	}
}

func HandleBooksFileDeleteTask(ctx context.Context, t *asynq.Task) error {
	var payload BooksFileDeletePayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	log := logger.New()

	storage := storage.New()
	err := storage.DeleteFile("books", []string{payload.Filename})
	if err != nil {
		log.Error("[Task][BooksFileDelete]", err)
		return err
	}

	fmt.Println("[Task][BooksFileDelete] Run Succesfully")
	return nil
}
