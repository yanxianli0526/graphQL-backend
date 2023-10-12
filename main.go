package main

import (
	log "log"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	orm "graphql-go-template/internal/database"
	"graphql-go-template/internal/handlers"
	"graphql-go-template/internal/restful/router"
	objectStorage "graphql-go-template/pkg/gcp"

	"gitlab.smart-aging.tech/devops/ms-go-kit/observability"

	"github.com/gin-contrib/cors"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	config "graphql-go-template/envconfig"
)

var (
	env    config.EnvConfig
	logger *zap.Logger
)

var gqlPath, gqlPgPath string
var isPgEnabled bool

func init() {
	gqlPath = "/api/graphql"
	gqlPgPath = "/"
	isPgEnabled = true
}

func run(msg string, err error, fields ...zapcore.Field) {
	logger.Info(msg, fields...)
	if err != nil {
		log.Fatalf("failed to run %s: %s", msg, err)
	}
}

func loadEnv() {
	if err := config.Process(&env); err != nil {
		log.Fatal("load config error > ", err)
	}
}

// Run spins up the server
func main() {
	loadEnv()
	var err error
	//
	// MUST: setup an instrumented logger
	//
	logger, err = observability.SetupLogger(env.Env.Debug)
	if err != nil {
		log.Fatalf("failed to setup logger: %s", err)
	}

	// Create a new ORM instance to send it to our server
	orm, err := orm.Factory(env.Database)
	if err != nil {
		logger.Error("orm.Factory error", zap.Error(err))
	}
	// DB Migration
	err = orm.UpdateMigration()
	if err != nil {
		logger.Error("UpdateMigration error", zap.Error(err))
	}

	err = orm.DBScript()
	if err != nil {
		logger.Error("DBScript", zap.Error(err))
	}

	storage, err := objectStorage.NewGCS("test-inventory-toll-files", "test-jubo-file")
	if err != nil {
		logger.Error("failed to connect object store :", zap.Error(err))
	}

	engine := observability.NewGinEngine(logger, env.Env.Debug)

	routerUrl := engine.Group("/api")
	engine.Use(cors.New(getCorsConfig()))
	routerUrl.Use(cors.New(getCorsConfig()))

	// restful api implement
	engine.GET("api/healthCheck", handlers.HealthCheck())
	router.RegisterNIS(orm, routerUrl, env.AuthEndpoint)

	//
	// GraphQL handlers & Playground handler
	//

	if isPgEnabled {
		engine.GET(gqlPgPath, handlers.PlaygroundHandler(gqlPath))
		logger.Info("GraphQL Playground " + strconv.Itoa(env.Env.HTTPPort) + gqlPgPath)
	}
	// Pass in the ORM instance to the GraphqlHandler
	engine.POST(gqlPath, handlers.GraphqlHandler(orm, storage, env.AuthEndpoint, env.PhotoService, logger))
	logger.Info("GraphQL @ " + strconv.Itoa(env.Env.HTTPPort) + gqlPath)

	//
	// run server
	//
	run("start gin server", observability.StartGinServer(engine, env.Env.HTTPPort), zap.Int("port", env.Env.HTTPPort))

	//
	// RECOMMENDED: implement your graceful shutdown logic
	//
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("shutdown servers...")

	//
	// RECOMMENDED: use graceful shutdown methods to close the server(s)
	//
	observability.ShutdownGinServer(observability.WithGracefulPeriod(1))
}

// getCorsConfig generates a config to use in gin cors middleware based on server configuration
func getCorsConfig() cors.Config {
	corsConf := cors.DefaultConfig()

	corsConf.AllowOrigins = []string{"http://localhost:8010"}
	corsConf.AllowMethods = []string{"GET", "POST", "DELETE", "PATCH", "OPTIONS", "PUT"}
	corsConf.AllowHeaders = []string{"Authorization", "Content-Type", "Upgrade", "Origin",
		"Connection", "Accept-Encoding", "Accept-Language", "Host"}
	corsConf.AllowCredentials = true
	return corsConf
}
