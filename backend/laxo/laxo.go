package laxo

import (
  "log"

  "github.com/gorilla/mux"
  "github.com/urfave/negroni"
  "github.com/hashicorp/hcl/v2/hclsimple"
  "github.com/joho/godotenv"
  "github.com/hashicorp/go-hclog"
)

type Config struct {
  Port     string `hcl:"port"`
  LogLevel string `hcl:"log_level"`
}

var Logger hclog.Logger

func InitConfig(config *Config) hclog.Logger {
  if err := godotenv.Load(".env"); err != nil {
    log.Fatal("Failed to load .env file")
  }

  err := hclsimple.DecodeFile("config.hcl", nil, config)
  if err != nil {
    log.Fatalf("Failed to load configuration: %r", err)
  }

  appLogger := hclog.New(&hclog.LoggerOptions{
    Name:  "laxo-backend",
    Level: hclog.LevelFromString(config.LogLevel),
  })

  Logger = appLogger

  appLogger.Info("Configuration loaded", "LogLevel", config.LogLevel)

  return appLogger
}

func SetupRouter() *negroni.Negroni {
  n := negroni.New()

  // common middleware
  n.Use(negroni.HandlerFunc(assureJSON))

  r := mux.NewRouter()

  s := r.PathPrefix("/api").Subrouter()
  s.HandleFunc("/login", handleLogin).Methods("POST")
  s.HandleFunc("/user", handleCreateUser).Methods("POST")

  n.UseHandler(r)
  return n
}

