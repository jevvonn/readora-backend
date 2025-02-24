package postgresql

import (
	"github.com/jevvonn/readora-backend/internal/domain/entity"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB, command string) {
	migrator := db.Migrator()
	tables := []any{
		&entity.User{},
	}

	createUserRoleEnum := `
	DO $$ 
	BEGIN 
	    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE LOWER(typname) = 'userrole') THEN 
	        CREATE TYPE "userRole" AS ENUM('ADMIN', 'USER'); 
	    END IF; 
	END $$;
	`
	db.Exec(createUserRoleEnum)

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
