package laxo

import (
  "log"

  "github.com/gorilla/mux"
  "github.com/urfave/negroni"
  "github.com/hashicorp/hcl/v2/hclsimple"
  "github.com/hashicorp/go-hclog"
)

type Config struct {
  Port     string `hcl:"port"`
  LogLevel string `hcl:"log_level"`
}

var Logger hclog.Logger

func InitConfig(config *Config, testing bool) hclog.Logger {
  err := hclsimple.DecodeFile("config.hcl", nil, config)
  if err != nil {
    log.Fatalf("Failed to load configuration: %r", err)
  }

  logLevel := hclog.Off

  if(!testing) {
    logLevel = hclog.LevelFromString(config.LogLevel)
  }

  appLogger := hclog.New(&hclog.LoggerOptions{
    Name:  "laxo-backend",
    Level: logLevel,
  })

  Logger = appLogger

  appLogger.Info("Configuration loaded", "LogLevel", config.LogLevel)

  return appLogger
}

func SetupRouter() *mux.Router {
  router := mux.NewRouter()

  // Common middlewares
  common := negroni.New()

  subRouter := router.PathPrefix("/api").Subrouter().StrictSlash(true)

	subRouter.Handle("/user", common.With(
		negroni.HandlerFunc(assureJSON),
		negroni.WrapFunc(HandleCreateUser),
	)).Methods("POST")

	subRouter.Handle("/test", common.With(
		negroni.WrapFunc(HandleTest),
	)).Methods("GET")


  return subRouter
}

