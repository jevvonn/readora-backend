package bootstrap

import (
	"flag"
	"os"

	"github.com/jevvonn/readora-backend/internal/infra/postgresql"
	"gorm.io/gorm"
)

func CommandHandler(db *gorm.DB) {
	var migrationCmd string
	var seederCmd bool

	flag.StringVar(&migrationCmd, "m", "", "Migrate database 'up' or 'down'")
	flag.BoolVar(&seederCmd, "s", false, "Seed database")
	flag.Parse()

	if migrationCmd != "" {
		postgresql.Migrate(db, migrationCmd)
		os.Exit(0)
	}

	if seederCmd {
		postgresql.Seed(db)
		os.Exit(0)
	}
}
