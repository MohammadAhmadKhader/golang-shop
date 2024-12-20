package database

import (
	"database/sql"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"main.go/config"
	"main.go/pkg/models"
)

var DB *gorm.DB = InitDB()

func InitDB() *gorm.DB {
	Envs := config.Envs
	config := &gorm.Config{}

	// adding queries logger
	if Envs.Env == "development" {
		config.Logger = SQLLogger
	}
	dsnNoDB := Envs.DSN_NO_DB
	dbName := Envs.DBName

	db, err := sql.Open("mysql", dsnNoDB)
	if err != nil {
		panic(err)
	}
	createDbIfNotExistQ := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", dbName)
	_, err = db.Exec(createDbIfNotExistQ)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	DB, err := gorm.Open(mysql.Open(Envs.DSN), config)
	if err != nil {
		panic(err)
	}

	err = DB.AutoMigrate(
		&models.Category{}, &models.Product{}, &models.User{},
		&models.CartItem{}, &models.Review{}, &models.Order{},
		&models.OrderItem{}, &models.Image{}, &models.Role{},
		&models.UserRoles{}, &models.Address{}, &models.Message{},
	)
	if err != nil {
		panic(err)
	}
	err = DB.SetupJoinTable(&models.User{}, "Roles", &models.UserRoles{})
	if err != nil {
		panic(err)
	}

	SeedData(DB)
	
	return DB
}