package app

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/r-pine/demo_aggregation/app/internal/blockchain"
	"github.com/r-pine/demo_aggregation/app/internal/controller"
	"github.com/r-pine/demo_aggregation/app/internal/db/redis"
	sc "github.com/r-pine/demo_aggregation/app/internal/service"
	st "github.com/r-pine/demo_aggregation/app/internal/storage"
	"github.com/r-pine/demo_aggregation/app/pkg/config"
	"github.com/r-pine/demo_aggregation/app/pkg/logging"
	"github.com/r-pine/demo_aggregation/app/pkg/server"
)

func RunApplication() {
	// Init Logger
	logging.Init()
	log := logging.GetLogger()
	log.Infoln("Connect logger successfully!")

	// Init Config
	cfg := config.GetConfig()
	log.Infoln("Connect config successfully!")

	// Init Context
	ctx := context.Background()
	rcClient := redis.NewRedisClient(ctx, cfg, log)
	rc, err := rcClient.ConnectToRedis()
	if err != nil {
		log.Fatalf("redis connect to redis failed: %v", err)
		return
	}

	storage := st.NewStorage(ctx, rc)
	service := sc.NewService(ctx, storage)
	aggregation := blockchain.NewAggregation(ctx, *cfg, log)

	gin.SetMode(cfg.AppConfig.GinMode)
	ginRouter := gin.New()
	httpController := controller.NewController(log, *service, *cfg, aggregation)
	handlers := httpController.InitRoutes(ginRouter)

	server.RunServer(log, handlers, cfg.AppConfig.HttpAddr)
}
