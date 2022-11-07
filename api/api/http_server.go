package api

import (
	"os"
	"strings"
	"time"

	"github.com/lucaswiix/meli/notifications/pkg/cors"
	"github.com/lucaswiix/meli/notifications/repository"
	"github.com/lucaswiix/meli/notifications/service"
	"github.com/lucaswiix/meli/notifications/usecase"
	"github.com/lucaswiix/meli/notifications/utils"
	"github.com/prometheus/client_golang/prometheus"
	"go.elastic.co/apm/module/apmgin"

	_ "time/tzdata"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

var notificationStatus = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_request_get_user_status_count", // metric name
		Help: "Count of status returned by user.",
	},
	[]string{"user", "status"}, // labels
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
	ginServer.Use(apmgin.Middleware(ginServer))
	// ginServer.Use(gin.Recovery())
	ginServer.Use(gzip.Gzip(gzip.BestSpeed))
	ginServer.Use(cors.CORS())

	skipPaths := strings.Split(os.Getenv("LOG_SKIP_PATH"), ",")
	ginServer.Use(gin.LoggerWithConfig(gin.LoggerConfig{SkipPaths: skipPaths}))

	RegisterHealthCheckHandlers(ginServer)
	prometheus.MustRegister(notificationStatus)

	RegisterMetricsHandlers(ginServer)
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
