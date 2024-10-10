package main

import (
	"context"
	"fmt"
	"github.com/LittleAksMax/blog-backend/internal/api"
	"github.com/LittleAksMax/blog-backend/internal/auth"
	"github.com/LittleAksMax/blog-backend/internal/cache"
	"github.com/LittleAksMax/blog-backend/internal/config"
	"github.com/LittleAksMax/blog-backend/internal/db"
	"github.com/LittleAksMax/blog-backend/internal/logging"
	"github.com/gin-gonic/gin"
	"log"
	"runtime"
)

func main() {
	fmt.Println("Running application!")

	ctx := context.Background()

	// make sure at least 3 usable threads
	runtime.GOMAXPROCS(max(runtime.GOMAXPROCS(-1), 3))

	if gin.Mode() == gin.DebugMode {
		config.InitDotenv(".env.Dev")
	} else if gin.Mode() == gin.ReleaseMode {
		gin.DisableConsoleColor()
		// config.InitDotenv(".env")
	} else {
		log.Fatalf("Unsupported Gin Mode: %s", gin.Mode())
	}

	cfg := config.InitConfig()

	logWriter := logging.InitLogging(cfg.LogFile)

	dbCfg := db.InitDb(ctx, cfg.DbHost, cfg.DbPort, cfg.DbUser, cfg.DbPasswd, cfg.DbName)
	defer dbCfg.CloseDb()

	cacheCfg := cache.InitCache(ctx, cfg.CacheHost, cfg.CachePort, cfg.CachePasswd)
	defer cacheCfg.CloseCache()

	authCfg := auth.InitAuth(ctx, cfg.FirebaseProjectID, cfg.FirebaseCredentialFile)

	api.RunApi(cfg.ApiPort, cfg.CorsAllowedOrigins, logWriter, dbCfg, cacheCfg, authCfg)
}
