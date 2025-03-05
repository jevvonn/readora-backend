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
	createBookFileStatusEnum := `
	DO $$ 
	BEGIN 
		IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'bookfilestatus') THEN 
			CREATE TYPE bookfilestatus AS ENUM ('PROCESSING', 'READY'); 
		END IF; 
	END $$;
	`
	db.Exec(createUserRoleEnum)
	db.Exec(createBookFileStatusEnum)

	tables := []any{
		&entity.User{},
		&entity.Book{},
		&entity.Genre{},
	}

	var err error
	if command == "up" {
		err = migrator.AutoMigrate(tables...)
	}

	if command == "down" {
		err = migrator.DropTable(tables...)
	}

	if err != nil {
		panic(err)
	}
}
