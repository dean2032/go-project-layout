package repo

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/dean2032/go-project-layout/config"
	"github.com/dean2032/go-project-layout/utils/errors"
	"github.com/dean2032/go-project-layout/utils/logging"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"gorm.io/plugin/dbresolver"
)

// Database modal
type Database struct {
	*gorm.DB
}

// NewDatabase creates a new database instance
func NewDatabase(cfg *config.Config) *Database {
	db, err := connectTo(cfg.MainDB, NewGormConfig(cfg), cfg.DBConnectionPoolSize)

	if err != nil {
		logging.Infof("Url: ", cfg.MainDB)
		logging.Panic(err.Error())
	}

	logging.Info("Database connection established")

	return &Database{
		DB: db,
	}
}

func setPoolParam(db *gorm.DB, connPoolSize int) {
	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}
	sqlDB.SetMaxIdleConns(connPoolSize)
	sqlDB.SetMaxOpenConns(connPoolSize)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)
	sqlDB.SetConnMaxIdleTime(5 * time.Minute)
}

func connectTo(dbUrls string, gormConfig gorm.Config, connPoolSize int) (*gorm.DB, error) {
	if connPoolSize == 0 {
		connPoolSize = 256
	}
	urls := strings.Split(dbUrls, "\n")
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       urls[0],
		SkipInitializeWithVersion: true,
	}), &gormConfig)
	if err != nil {
		return nil, errors.Wrap(err, "mysql can't connect")
	}
	replicas := make([]gorm.Dialector, 0, len(urls))
	for _, url := range urls[1:] {
		if len(url) > 0 {
			replicas = append(replicas, mysql.New(mysql.Config{DSN: url, SkipInitializeWithVersion: true}))
		}
	}
	if len(replicas) > 0 {
		err = db.Use(dbresolver.Register(dbresolver.Config{Replicas: replicas}).
			SetConnMaxIdleTime(time.Hour).
			SetConnMaxLifetime(time.Hour).
			SetMaxIdleConns(connPoolSize).
			SetMaxOpenConns(connPoolSize))
		if err != nil {
			return nil, errors.WithStack(err)
		}
	}
	setPoolParam(db, connPoolSize)
	return db, nil
}

// NewGormConfig make gormConfig
func NewGormConfig(cfg *config.Config) gorm.Config {
	config := gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		SkipDefaultTransaction:                   true,
		DisableForeignKeyConstraintWhenMigrating: true,
	}
	if cfg.Debug {
		config.Logger = gormLogger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
			gormLogger.Config{
				SlowThreshold: time.Second,
				LogLevel:      gormLogger.Info,
				Colorful:      true,
			},
		)
	} else {
		config.Logger = logging.NewGormLogger(zap.WarnLevel, time.Second)
	}
	return config
}
