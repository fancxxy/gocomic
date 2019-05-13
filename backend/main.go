package main

import (
	"log"

	"github.com/fancxxy/gocomic/backend/config"
	"github.com/fancxxy/gocomic/backend/cron"
	"github.com/fancxxy/gocomic/backend/models"
	"github.com/fancxxy/gocomic/backend/routers"
	"github.com/fancxxy/gocomic/backend/rpc"
)

func main() {
	if err := config.Init(); err != nil {
		log.Fatalf("conf.Init() err: %v", err)
	}

	if err := rpc.Init(); err != nil {
		log.Fatalf("rpc.Init() err: %v", err)
	}

	if err := models.Init(); err != nil {
		log.Fatalf("models.Init() err: %v", err)
	}

	if err := cron.Init(); err != nil {
		log.Fatalf("cron.Init() err: %v", err)
	}

	if err := routers.Init(); err != nil {
		log.Fatalf("routers.Init() err: %v", err)
	}
}
