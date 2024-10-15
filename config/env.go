package config

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"
)

type Config struct {
	PublicHost             string
	Port                   string
	DBUser                 string
	DBPassword             string
	DBName                 string
	DBAddress              string
	JWT_EXPIRATION_IN_SECONDS string
	JWT_SECRET              string
	Env                    string
	CLOUDINARY_APIKEY      string
	CLOUDINARY_SECRET      string
	CLOUDINARY_NAME 		string
	DSN string
	DSN_NO_DB string
	AUTH_STORE_KEY string
}

var Envs = initConfig()
//dsn := "user:pass@tcp(127.0.0.1:3306)/go-shop?charset=utf8mb4&parseTime=True&loc=Local"
func initConfig() Config {
	loadEnvFile()

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
		CLOUDINARY_NAME: getEnv("CLOUDINARY_NAME", ""),
		DSN: fmt.Sprintf("%s:%s@tcp(%s%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		getEnv("DB_USER",""),getEnv("DB_PASSWORD",""),getEnv("DB_HOST",""),
		getEnv("DB_PORT",""),getEnv("DB_NAME","")),
		DSN_NO_DB: fmt.Sprintf("%s:%s@tcp(%s%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		getEnv("DB_USER",""),getEnv("DB_PASSWORD",""),getEnv("DB_HOST",""),
		getEnv("DB_PORT",""),""),
		AUTH_STORE_KEY: getEnv("AUTH_STORE_KEY",""),
		JWT_SECRET: getEnv("JWT_SECRET",""),
		JWT_EXPIRATION_IN_SECONDS: getEnv("JWT_EXPIRATION_IN_SECONDS",""),
	}
}

var envs map[string]string

func getEnv(key, fallback string) string {
	if envs[key] != "" {
		return envs[key]
	}

	return fallback
}

func loadEnvFile() {
	godotenv.Load()
	envVariables, err := godotenv.Read()
	if err != nil {
		log.Fatal(err)
	}

	envs = envVariables
}
