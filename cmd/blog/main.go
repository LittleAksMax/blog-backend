package main

import (
	"context"
	"fmt"
	"github.com/LittleAksMax/blog-backend/internal/api"
	"github.com/LittleAksMax/blog-backend/internal/config"
	"github.com/LittleAksMax/blog-backend/internal/db"
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
		config.InitDotenv(".env")
	} else {
		log.Fatalf("Unsupported Gin Mode: %s", gin.Mode())
	}

	cfg := config.InitConfig()

	dbCfg := db.InitDb(ctx, cfg.DbHost, cfg.DbPort, cfg.DbUser, cfg.DbPasswd, cfg.DbName)
	defer dbCfg.CloseDb()

	api.RunApi(cfg.ApiPort, dbCfg)
}
