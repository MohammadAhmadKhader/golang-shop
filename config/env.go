package config

import (
	"fmt"
	"path/filepath"
	"runtime"

	"github.com/joho/godotenv"
)

type Config struct {
	PublicHost                string
	Port                      string
	DBUser                    string
	DBPassword                string
	DBName                    string
	DBAddress                 string
	ACCESS_JWT_EXPIRATION_IN_SECONDS string
	REFRESH_JWT_EXPIRATION_IN_SECONDS string
	JWT_SECRET                string
	Env                       string
	CLOUDINARY_APIKEY         string
	CLOUDINARY_SECRET         string
	CLOUDINARY_NAME           string
	DSN                       string
	DSN_NO_DB                 string
	AUTH_STORE_KEY            string
}

var Envs = initConfig()

// dsn := "user:pass@tcp(127.0.0.1:3306)/go-shop?charset=utf8mb4&parseTime=True&loc=Local"
func initConfig() Config {
	err := loadEnvFile()
	if err != nil {
		panic(err)
	}

	return Config{
		PublicHost:        getEnv("PUBLIC_HOST", "http://localhost"),
		Port:              getEnv("PORT", ":8080"),
		DBUser:            getEnv("DB_USER", "root"),
		DBPassword:        getEnv("DB_PASSWORD", "mypassword"),
		DBAddress:         fmt.Sprintf("%s:%s", getEnv("DB_HOST", "127.0.0.1"), getEnv("DB_PORT", "3306")),
		DBName:            getEnv("DB_NAME", "dbname"),
		Env:               getEnv("env", "production"),
		CLOUDINARY_APIKEY: getEnv("CLOUDINARY_APIKEY", ""),
		CLOUDINARY_SECRET: getEnv("CLOUDINARY_SECRET", ""),
		CLOUDINARY_NAME:   getEnv("CLOUDINARY_NAME", ""),
		DSN: fmt.Sprintf("%s:%s@tcp(%s%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			getEnv("DB_USER", ""), getEnv("DB_PASSWORD", ""), getEnv("DB_HOST", ""),
			getEnv("DB_PORT", ""), getEnv("DB_NAME", "")),
		DSN_NO_DB: fmt.Sprintf("%s:%s@tcp(%s%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			getEnv("DB_USER", ""), getEnv("DB_PASSWORD", ""), getEnv("DB_HOST", ""),
			getEnv("DB_PORT", ""), ""),
		AUTH_STORE_KEY:            getEnv("AUTH_STORE_KEY", ""),
		JWT_SECRET:                getEnv("JWT_SECRET", ""),
		ACCESS_JWT_EXPIRATION_IN_SECONDS: getEnv("ACCESS_JWT_EXPIRATION_IN_SECONDS", ""),
		REFRESH_JWT_EXPIRATION_IN_SECONDS:  getEnv("REFRESH_JWT_EXPIRATION_IN_SECONDS", ""),
	}
}

var envs map[string]string

func getEnv(key, fallback string) string {
	if envs[key] != "" {
		return envs[key]
	}

	return fallback
}

func loadEnvFile() error {
	err := godotenv.Load()
	if err != nil {
		err = handleTestEnv()
		if err != nil { 
			return err
		}
		return nil
	}
	envVariables, err := godotenv.Read()
	if err != nil {
		return err
	}

	envs = envVariables
	return nil
}

// this function used to handle the tests environment, the test files in sub directory can not see the main ".env" file
//
// TODO: must be reworked to a better approach
func handleTestEnv() error {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return fmt.Errorf("failed to load .env.test file")
	}
	basepath := filepath.Dir(file)
	err := godotenv.Load(filepath.Join(basepath, "../.env.test"))
	if err != nil {
		return err
	}

	envVariables, err := godotenv.Read(filepath.Join(basepath, "../.env.test"))
	if err != nil {
		return err
	}

	envs = envVariables
	return err
}
