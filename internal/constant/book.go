package constant

import (
	"fmt"

	"github.com/jevvonn/readora-backend/config"
)

const (
	BookFileUploadStatusQueue     = "QUEUE"
	BookFileUploadStatusUploading = "UPLOADING"
	BookFileUploadStatusUploaded  = "UPLOADED"

	BookFileAIStatusQueue      = "QUEUE"
	BookFileAIStatusProcessing = "PROCESSING"
	BookFileAIStatusReady      = "READY"
)

func GetBookDefultCoverImage() string {
	conf := config.Load()
	return fmt.Sprintf("%s/storage/v1/object/public/images//default-book-cover.png", conf.SupabaseProjectURL)
}

func GetBookTxtFile(bookId string) string {
	conf := config.Load()
	return fmt.Sprintf("%s/storage/v1/object/public/texts/%s.txt", conf.SupabaseProjectURL, bookId)
}
