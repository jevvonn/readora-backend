package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"github.com/jevvonn/readora-backend/internal/constant"
	"github.com/jevvonn/readora-backend/internal/domain/entity"
	"github.com/jevvonn/readora-backend/internal/infra/logger"
	"github.com/jevvonn/readora-backend/internal/infra/storage"
	"gorm.io/gorm"
)

// A list of task types.
const (
	BooksFileUploadTaskName = "file-upload:books"
)

type BooksFileUploadPayload struct {
	TmpFile  string
	Filename string
	BooksId  string
}

// Create Task
func HandleBooksFileUploadTask(db *gorm.DB) asynq.HandlerFunc {
	return func(ctx context.Context, t *asynq.Task) error {
		var payload BooksFileUploadPayload
		if err := json.Unmarshal(t.Payload(), &payload); err != nil {
			return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
		}

		log := logger.New()

		// Check path
		_, err := os.Stat(payload.TmpFile)
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
		publicUrl, err := storage.UploadFile(file, "books", payload.Filename, "application/pdf")
		if err != nil {
			log.Error("[Task][BooksFileUpload]", err)
			return err
		}

		// Update file_url in books
		err = db.Model(&entity.Book{
			ID: uuid.MustParse(payload.BooksId),
		}).Updates(entity.Book{
			FileURL:        publicUrl,
			BookFileStatus: constant.BookFileStatusReady,
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
