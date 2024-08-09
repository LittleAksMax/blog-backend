package config

import (
	"github.com/lpernett/godotenv"
	"log"
	"os"
	"strconv"
)

// Keys used in environment file
const (
	apiPortKey                = "PORT"
	apiKeyKey                 = "API_KEY"
	firebaseProjectIdKey      = "FIREBASE_PROJECT_ID"
	firebaseCredentialFileKey = "FIREBASE_CREDENTIAL_FILE"
	dbHostnameKey             = "MONGO_HOSTNAME"
	dbPortKey                 = "MONGO_PORT"
	dbUserKey                 = "MONGO_INITDB_ROOT_USERNAME"
	dbPasswdKey               = "MONGO_INITDB_ROOT_PASSWORD"
	dbNameKey                 = "MONGO_DATABASE_NAME"
	cacheHostKey              = "REDIS_HOSTNAME"
	cachePortKey              = "REDIS_PORT"
	cachePasswdKey            = "REDIS_PASSWORD"
)

type Config struct {
	ApiPort int
	ApiKey  string

	FirebaseProjectID      string
	FirebaseCredentialFile string

	DbHost   string
	DbPort   int
	DbUser   string
	DbPasswd string
	DbName   string

	CacheHost   string
	CachePort   int
	CachePasswd string
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
	cachePort := readPort(cachePortKey)

	return &Config{
		ApiPort: apiPort,
		ApiKey:  os.Getenv(apiKeyKey),

		FirebaseProjectID:      os.Getenv(firebaseProjectIdKey),
		FirebaseCredentialFile: os.Getenv(firebaseCredentialFileKey),

		DbHost:   os.Getenv(dbHostnameKey),
		DbPort:   dbPort,
		DbUser:   os.Getenv(dbUserKey),
		DbPasswd: os.Getenv(dbPasswdKey),
		DbName:   os.Getenv(dbNameKey),

		CacheHost:   os.Getenv(cacheHostKey),
		CachePort:   cachePort,
		CachePasswd: os.Getenv(cachePasswdKey),
	}
}
