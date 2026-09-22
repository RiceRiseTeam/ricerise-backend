package database

import (
	"fmt"
	"ricerise/internal/config"
	"ricerise/internal/model"
	"time"

	"github.com/samber/do/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgres(injector do.Injector) (*gorm.DB, error) {
	appConfig := do.MustInvoke[*config.AppConfig](injector)

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		appConfig.SQLHost,
		appConfig.SQLUser,
		appConfig.SQLPassword,
		"ricerise",
		appConfig.SQLPort,
	)

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: dsn,
	}))
	if err != nil {
		panic("failed to open mysql connection: " + err.Error())
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic("failed to init mysql database: " + err.Error())
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)
	sqlDB.SetConnMaxIdleTime(time.Minute * 30)

	err = db.AutoMigrate(
		&model.UserModel{},
		&model.LocationModel{},
		&model.CommentModel{},
		&model.DinnerModel{},
		&model.ParticipantModel{},
	)
	if err != nil {
		panic("failed to auto migrate mysql model: " + err.Error())
	}

	return db, nil
}
