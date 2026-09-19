package database

import (
	"ricerise/internal/config"
	"ricerise/internal/model"
	"time"

	gomysql "github.com/go-sql-driver/mysql"
	"github.com/samber/do/v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewMySQL(injector do.Injector) (*gorm.DB, error) {
	appConfig := do.MustInvoke[*config.AppConfig](injector)

	mysqlConfig := gomysql.Config{
		User:      appConfig.MySQLUser,
		Passwd:    appConfig.MySQLPassword,
		Net:       "tcp",
		Addr:      appConfig.MySQLUrl,
		DBName:    "ricerise",
		ParseTime: true,
		Loc:       time.Local,
		Params: map[string]string{
			"charset": "utf8mb4",
		},
	}

	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN: mysqlConfig.FormatDSN(),
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
