package storage

import (
	"os"

	"github.com/jevvonn/readora-backend/config"
	storage_go "github.com/supabase-community/storage-go"
	"github.com/supabase-community/supabase-go"
)

type StorageItf interface {
	UploadFile(file *os.File, bucket, fileName, mimeType string) (string, error)
	DeleteFile(bucket string, fileNames []string) error
}

type Storage struct {
	client *supabase.Client
}

func New() StorageItf {
	conf := config.Load()
	client, err := supabase.NewClient(conf.SupabaseProjectURL, conf.SupabaseProjectToken, nil)

	if err != nil {
		panic(err)
	}

	return &Storage{
		client,
	}
}

func (s *Storage) UploadFile(file *os.File, bucket, fileName, mimeType string) (string, error) {
	typeFile := new(string)
	*typeFile = mimeType
	_, err := s.client.Storage.UploadFile(bucket, fileName, file, storage_go.FileOptions{
		ContentType: typeFile,
	})
	if err != nil {
		return "", err
	}

	publicURL := s.client.Storage.GetPublicUrl(bucket, fileName).SignedURL

	return publicURL, nil
}

func (s *Storage) DeleteFile(bucket string, fileNames []string) error {
	_, err := s.client.Storage.RemoveFile(bucket, fileNames)
	return err
}
