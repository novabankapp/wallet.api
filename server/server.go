package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/novabankapp/common.infrastructure/logger"
	"github.com/novabankapp/wallet.api/controllers"
	"github.com/novabankapp/wallet.api/middlewares"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
	"time"
)

// NewServer @title Novabank Wallet API
// @version 1.0
// @description This is a Wallet API.
// @BasePath /
func NewServer(address string, port string, controllers controllers.Controllers, middlewares middlewares.Middlewares, logger logger.Logger) *http.Server {
	//ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	//defer cancel()
	// Init gin handler
	router := gin.Default()
	router.Use(
		gin.Recovery(),
		gin.Logger(),
	)
	router.Use(CORSMiddleware(), RequestID(), RequestLogger())

	//baseRoutes := router.Group("/api/v1")
	//registrationHttp.RegisterEndpoints(registrationRoutes, controllers.Registration, logger)
	//authMiddleware := middlewares.Auth.Handle
	//usersRoutes := router.Group("/api/v1")
	//usersHttp.UsersEndpoints(usersRoutes, controllers.Users, logger)
	//router.GET("try", Helloworld)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return &http.Server{
		Addr:           fmt.Sprintf("%s:%s", address, port),
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
}
