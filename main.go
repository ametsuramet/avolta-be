package main

import (
	"avolta/cmd"
	"avolta/config"
	"avolta/database"
	"avolta/router"
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

func init() {
	ctx := context.Background()
	config.LoadConfig()
	database.InitDB(ctx)
	initLog()
}

func initLog() {
	t := time.Now()
	filename := t.Format("2006-01-02")
	f, err := os.OpenFile("log/"+filename+".log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)
}

func main() {
	args := os.Args[1:]
	if config.App.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	if len(args) > 0 {
		fmt.Println("RUN", args)
		cmd.Init(args)
		os.Exit(0)
	}
	router := router.SetupRouter()

	router.Run(config.App.Server.Addr)
}
