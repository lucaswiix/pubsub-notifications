package api

import (
	"meli/notifications/pkg/cors"
	"meli/notifications/repository"
	"meli/notifications/service"
	"meli/notifications/usecase"
	"meli/notifications/utils"
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

func InitWebServer() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	gin.DisableConsoleColor()
	ginServer := gin.New()
	time.Local, _ = time.LoadLocation("America/Sao_Paulo")
	// Initial packages
	cors := cors.NewCors()
	utils.InitValidation()
	// Middleware
	ginServer.Use(gin.Recovery())
	ginServer.Use(gzip.Gzip(gzip.BestSpeed))
	ginServer.Use(cors.CORS())
	// ginServer.Use(middleware.JSONLogMiddleware())

	skipPaths := strings.Split(os.Getenv("LOG_SKIP_PATH"), ",")
	ginServer.Use(gin.LoggerWithConfig(gin.LoggerConfig{SkipPaths: skipPaths}))

	RegisterHealthCheckHandlers(ginServer)

	notifyRepository := repository.NewNotifyRepository(repository.DB)
	queueRepository := repository.NewQueueRepository(repository.CH)
	optOutRepository := repository.NewOptOutRepository(repository.DB)

	notifyService := service.NewNotifyService(notifyRepository)
	queueService := service.NewQueueService(queueRepository)
	optOutService := service.NewOptOutService(optOutRepository)

	sendNotificationUsecase := usecase.NewNotificationUsecase(notifyService, queueService, optOutService)

	RegisterNotifyHandlers(ginServer, sendNotificationUsecase)
	RegisterOptOutHandlers(ginServer, optOutService)
	return ginServer
}
