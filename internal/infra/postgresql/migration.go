package postgresql

import (
	"github.com/jevvonn/readora-backend/internal/domain/entity"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB, command string) {
	migrator := db.Migrator()
	createUserRoleEnum := `
	DO $$ 
	BEGIN 
		IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'userrole') THEN 
			CREATE TYPE userrole AS ENUM ('ADMIN', 'USER'); 
		END IF; 
	END $$;
	`
	createBookFileUploadStatusEnum := `
	DO $$ 
	BEGIN 
		IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'book_file_upload_status') THEN 
			CREATE TYPE book_file_upload_status AS ENUM ('QUEUE','UPLOADING', 'UPLOADED'); 
		END IF; 
	END $$;
	`

	createBookFileAIStatusEnum := `
	DO $$ 
	BEGIN 
		IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'book_file_ai_status') THEN 
			CREATE TYPE book_file_ai_status AS ENUM ('QUEUE','PROCESSING', 'READY'); 
		END IF; 
	END $$;
	`
	db.Exec(createUserRoleEnum)
	db.Exec(createBookFileUploadStatusEnum)
	db.Exec(createBookFileAIStatusEnum)

	tables := []any{
		&entity.User{},
		&entity.Book{},
		&entity.Comment{},
		&entity.Reply{},
		&entity.Genre{},
	}

	var err error
	if command == "up" {
		err = migrator.AutoMigrate(tables...)
	}

	if command == "down" {
		err = migrator.DropTable(tables...)
		db.Exec(`
			DROP SCHEMA public CASCADE;
			CREATE SCHEMA public;

			GRANT ALL ON SCHEMA public TO public;
		`)
	}

	if err != nil {
		panic(err)
	}
}
