package store

import (
	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type DbHandler struct {
	db     *gorm.DB
	logger *log.Entry
}

func InitialiseDB(path string, logger *log.Entry, reset bool) (*DbHandler, error) {
	// Connect to database
	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// If Resets flag is activated then purge the db
	if reset {
		db.Unscoped().Where("1 = 1").Delete(&ProcessedBlock{})
		db.Unscoped().Where("1 = 1").Delete(&HanldeMsgTxq{})
	}

	// Migrate the schema
	logger.Info("Migrate the DB schema")
	db.AutoMigrate(&ProcessedBlock{})
	db.AutoMigrate(&HanldeMsgTxq{})

	// create Db handler
	dbHandler := &DbHandler{
		db:     db,
		logger: logger,
	}

	return dbHandler, nil
}
