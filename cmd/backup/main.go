package main

import (
	"flag"
	"fmt"
)

type TableData struct {
	TableName string           `json:"table_name"`
	Rows      []map[string]any `json:"rows"`
}

type DatabaseFlags struct {
	Host       *string
	User       *string
	Name       *string
	Password   *string
	Port       *string
	BackupFile *string
}

var (
	ScannedDatabaseFlags DatabaseFlags
	backupDir = "./backups"
) 

// Example use of restore
// make backup ARGS="--restore --dbname=golang_test --password=123456 --backupFile=data_backup_20241025_015150.json"
// or if you Makefiles dont work on your pc, you can use it without make files
// go run ./cmd/backup/ --restore --dbname=golang_test --password=123456 --backupFile=data_backup_20241025_015150.json
// ** the file must be inside "backups" folder which is the backupDir

// Example use of create
// make backup ARGS="--create"
// go run ./cmd/backup/ --create

// ** database must be created, if not will throw an error that database with specified name does nto exist
func main() {
	createBackupFlag := flag.Bool("create", false, "Create a new backup file")
	restoreDataFlag := flag.Bool("restore", false, "Restore data to a database connection")

	dbFlags := DatabaseFlags{
		Host:       flag.String("host", "127.0.0.1", "DB Host"),
		Password:   flag.String("password", "", "Database password"),
		User:       flag.String("user", "root", "Database user"),
		Port:       flag.String("port", "3306", "Database port (default 3306)"),
		Name:       flag.String("dbname", "", "Database name"),
		BackupFile: flag.String("backupFile", "", "backup file"),
	}
	ScannedDatabaseFlags = dbFlags

	flag.Parse()
	if !*createBackupFlag && !*restoreDataFlag {
		fmt.Println("You must specify the required operation whether its restore backup (--restore) or create backup (--create)")
		return
	}

	if *createBackupFlag {
		createBackup()
	}

	if *restoreDataFlag {
		restoreData()
	}
}
