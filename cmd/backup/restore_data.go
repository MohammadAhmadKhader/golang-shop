package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"main.go/pkg/models"
)

func restoreData() {
	err := validateInputs()
	if err != nil {
		log.Fatal(err)
		return
	}

	dsn := createDSN()
	backupFile := ScannedDatabaseFlags.BackupFile
	DB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
		return
	}

	err = migrateDB(DB)
	if err != nil {
		log.Fatal(err)
		return
	}

	filePath := filepath.Join(backupDir, *backupFile)
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("an error has occurred during opening backup file with path: '%v'", *backupFile)
		return
	}
	defer file.Close()

	fileData, err := io.ReadAll(file)
	if err != nil {
		log.Fatal("an error has occurred during reading from backup file: ", err)
		return
	}

	var tables []TableData
	err = json.Unmarshal(fileData, &tables)
	if err != nil {
		log.Fatal("an error has occurred during parsing backup file: ", err)
		return
	}


	disableForeignKeyChecks(DB)
	defer enableForeignKeyChecks(DB)
	for _, table := range tables {
		log.Printf("Restoring data to table %v\n has begun", table.TableName)

		for _, row := range table.Rows {
			err := DB.Table(table.TableName).Create(row).Error
			if err != nil {
				log.Printf("an error has occurred during inserting row to table: %v: %v\n", table.TableName, err)
			}
		}

		log.Printf("Restoring data to table %v\n has finished", table.TableName)
	}
}

func validateInputs() error {
	dbPassword := ScannedDatabaseFlags.Password
	dbPort := ScannedDatabaseFlags.Port
	dbName := ScannedDatabaseFlags.Name
	backupFile := ScannedDatabaseFlags.BackupFile

	if *dbName == "" {
		return fmt.Errorf("database name is required (use --dbname)")
	}
	if *dbPassword == "" {
		return fmt.Errorf("database password is required (use --password)")
	}
	if *backupFile == "" {
		return fmt.Errorf("backup file is required (use --backupFile)")
	}

	_, err := strconv.Atoi(*dbPort)
	if err != nil {
		return fmt.Errorf("invalid database port, received '%v' ", *dbPort)
	}

	if len(*dbPort) == 0 {
		return fmt.Errorf("database user is required (use --user)")
	}

	return nil
}

func createDSN() string {
	dbUser := ScannedDatabaseFlags.User
	dbPassword := ScannedDatabaseFlags.Password
	dbPort := ScannedDatabaseFlags.Port
	dbName := ScannedDatabaseFlags.Name
	dbHost := ScannedDatabaseFlags.Host

	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8mb4&parseTime=True&loc=Local",
		*dbUser, *dbPassword, *dbHost,*dbPort, *dbName)

	return dsn
}

func migrateDB(db *gorm.DB) error {
	err := db.AutoMigrate(
		&models.Category{}, &models.Product{}, &models.User{},
		&models.CartItem{}, &models.Review{}, &models.Order{},
		&models.OrderItem{}, &models.Image{}, &models.Role{},
		&models.UserRoles{}, &models.Address{}, &models.Message{},
	)
	if err != nil {
		return err
	}
	err = db.SetupJoinTable(&models.User{}, "Roles", &models.UserRoles{})
	if err != nil {
		return err
	}

	return nil
}

func disableForeignKeyChecks(db *gorm.DB) error {
	return db.Exec("SET FOREIGN_KEY_CHECKS = 0;").Error
}

func enableForeignKeyChecks(db *gorm.DB) error {
	return db.Exec("SET FOREIGN_KEY_CHECKS = 1;").Error
}