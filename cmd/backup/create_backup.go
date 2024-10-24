package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"main.go/internal/database"
)

func createBackup() {
	DB := database.DB
	var tables []string
	err := DB.Raw("SHOW TABLES").Find(&tables).Error
	if err != nil {
		log.Fatal(err)
		return
	}

	now := time.Now()
	backupFileName := fmt.Sprintf("data_backup_%v.json", now.Format("20060102_150405"))
	filePath := filepath.Join(backupDir, backupFileName)
	err = os.MkdirAll(backupDir, os.ModePerm)
	if err != nil {
		log.Fatal(err)
		return
	}

	file, err := os.Create(filePath)
	if err != nil {
		log.Fatal("something went wrong during creating backup file: ", err)
		return
	}
	defer file.Close()

	var AllData []TableData
	for _, table := range tables {
		var rows []map[string]any
		err = DB.Table(table).Find(&rows).Error
		if err != nil {
			log.Fatal(err)
			return
		}

		AllData = append(AllData, TableData{
			TableName: table,
			Rows:      rows,
		})
	}

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(AllData); err != nil {
		log.Fatal("something went wrong during decoding file: ", err)
		return
	}

	log.Println("database backup was created successfully!")
}