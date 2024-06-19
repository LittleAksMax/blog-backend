package main

import (
	"fmt"

	"github.com/LittleAksMax/blog-backend/internal/api"
	"github.com/LittleAksMax/blog-backend/internal/db"
	"github.com/LittleAksMax/blog-backend/internal/s3"
)

func main() {
	fmt.Println("Running application!")
	s3.InitS3()
	db.InitDb()
	api.RunApi()
}
