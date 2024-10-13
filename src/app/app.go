package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"shin/src/app/views"
	"shin/src/config"
	"shin/src/lib"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-openapi/runtime/middleware"
)

func Init() *gin.Engine {

	router := gin.New()

	//Set Logger
	logger := lib.CreateGinLogger(os.Stdout, lib.LOGGER_TEXT_FORMATTER)
	logger.AddHook(lib.CreateLogrusDiscordHook(config.Config.Logger.Discord["shin_channel"]))
	router.Use(views.GinLoggerMiddleware(logger))

	//Gin Recovery
	router.Use(gin.Recovery())

	router.Use(func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
		defer cancel()
		c.Set("ctx", ctx)
		c.Next()
	})

	router.Use(cors.New(cors.Config{
		AllowOrigins:     config.Config.Cors.Origins,
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"*"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	views.Init(router)

	//docs
	opts := middleware.SwaggerUIOpts{SpecURL: "/swagger.yaml"}
	router.GET("/docs", gin.WrapH(middleware.SwaggerUI(opts, nil)))
	router.GET("/swagger.yaml", gin.WrapH(http.FileServer(http.Dir("./docs"))))

	return router
}

func Serve() {
	router := Init()
	router.Run(fmt.Sprintf("0.0.0.0:%d", config.Config.Port))
}
