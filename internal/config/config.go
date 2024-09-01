package config

import (
	"github.com/lpernett/godotenv"
	"log"
	"os"
	"strconv"
	"strings"
)

// Keys used in environment file
const (
	apiPortKey                = "PORT"
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
	corsAllowedOriginsKey     = "CORS_ALLOWED_ORIGINS"
)

type Config struct {
	ApiPort int

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

	CorsAllowedOrigins []string
}

func readPort(key string) int {
	port, err := strconv.Atoi(os.Getenv(key))

	if err != nil || port < 1024 || port > 65353 {
		log.Fatalf("Error converting environment variable <%s> to int between 1024 and 65353", key)
	}
	return port
}

func readList(value string) []string {
	return strings.Split(value, ",")
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

		CorsAllowedOrigins: readList(os.Getenv(corsAllowedOriginsKey)),
	}
}
