package db

import (
	"log"

	"github.com/jonathanCaamano/portfolio-gateway/internal/config"
	"github.com/jonathanCaamano/portfolio-gateway/internal/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Connect(cfg *config.Config) *gorm.DB {
	db, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		log.Fatalf("db connect: %v", err)
	}
	return db
}

func Migrate(db *gorm.DB) {
	if err := db.AutoMigrate(
		&domain.User{},
		&domain.Category{},
		&domain.Product{},
		&domain.RefreshToken{},
	); err != nil {
		log.Fatalf("migrate: %v", err)
	}
}
