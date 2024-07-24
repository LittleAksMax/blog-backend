package config

import (
	"github.com/lpernett/godotenv"
	"log"
	"os"
	"strconv"
)

// Keys used in environment file
const (
	apiPortKey    = "PORT"
	dbHostnameKey = "MONGO_HOSTNAME"
	dbPortKey     = "MONGO_PORT"
	dbUserKey     = "MONGO_INITDB_ROOT_USERNAME"
	dbPasswdKey   = "MONGO_INITDB_ROOT_PASSWORD"
	dbNameKey     = "MONGO_DATABASE_NAME"
)

type Config struct {
	ApiPort int

	DbHost   string
	DbPort   int
	DbUser   string
	DbPasswd string
	DbName   string
}

func readPort(key string) int {
	port, err := strconv.Atoi(os.Getenv(key))

	if err != nil || port < 1024 || port > 65353 {
		log.Fatalf("Error converting environment variable <%s> to int between 1024 and 65353", key)
	}
	return port
}

func InitDotenv(filenames ...string) {
	err := godotenv.Load(filenames...)
	if err != nil {
		log.Fatal("Error loading environment file")
	}
}

func InitConfig() *Config {
	apiPort := readPort(apiPortKey)
	dbPort := readPort(dbPortKey)

	return &Config{
		ApiPort: apiPort,

		DbHost:   os.Getenv(dbHostnameKey),
		DbPort:   dbPort,
		DbUser:   os.Getenv(dbUserKey),
		DbPasswd: os.Getenv(dbPasswdKey),
		DbName:   os.Getenv(dbNameKey),
	}
}
