package database

import (
	"fmt"
	"gofiver/internal/config"
	"gofiver/internal/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Database struct {
	DB *gorm.DB
}

func NewDatabase(cfg config.DatabaseConfig) (*Database, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local&interpolateParams=true",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
	)

	logLevel := logger.Warn
	if cfg.Host == "localhost" || cfg.Host == "mysql" || cfg.Host == "common-mysql" {
		logLevel = logger.Silent
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger:                 logger.Default.LogMode(logLevel),
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
	})
	if err != nil {
		return nil, err
	}


	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}


	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	sqlDB.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)

	return &Database{DB: db}, nil
}

func (d *Database) AutoMigrate() error {
	return d.DB.AutoMigrate(
		&models.User{},
		&models.Blog{},
	)
}

func (d *Database) CreateIndexes() error {
	indexes := []string{
		"CREATE INDEX IF NOT EXISTS idx_blogs_id_desc ON blogs (id DESC)",
		"CREATE INDEX IF NOT EXISTS idx_blogs_deleted_at ON blogs (deleted_at)",
		"CREATE INDEX IF NOT EXISTS idx_blogs_user_id ON blogs (user_id)",
		"CREATE INDEX IF NOT EXISTS idx_blogs_title ON blogs (title(100))",
		"CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON users (deleted_at)",
	}

	for _, idx := range indexes {
		d.DB.Exec(idx)
	}
	return nil
}
