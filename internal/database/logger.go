package database

import (
	"log"
	"os"
	"time"

	"gorm.io/gorm/logger"
)

var SQLLogger = logger.New(
	log.New(os.Stdout, "\r\n", log.LstdFlags),
	logger.Config{
		SlowThreshold:             time.Second,
		LogLevel:                  logger.Info,
		IgnoreRecordNotFoundError: true,          
		ParameterizedQueries:      true,          
		Colorful:                  true,         
	},
)