package laxo

import (
  "log"

  "github.com/hashicorp/hcl/v2/hclsimple"
  "github.com/hashicorp/go-hclog"
  "github.com/gorilla/mux"
)

type Config struct {
  Port string `hcl:"port"`
  LogLevel string `hcl:"log_level"`
}


func InitConfig(config *Config) hclog.Logger {
  err := hclsimple.DecodeFile("config.hcl", nil, config)
  if err != nil {
    log.Fatalf("Failed to load configuration: %r", err)
  }

  appLogger := hclog.New(&hclog.LoggerOptions{
    Name:  "laxo-backend",
    Level: hclog.LevelFromString(config.LogLevel),
  })

  appLogger.Info("Configuration loaded", "LogLevel", config.LogLevel)

  return appLogger
}

func SetupRouter() *mux.Router {
  r := mux.NewRouter()
  s := r.PathPrefix("/api").Subrouter()
  s.HandleFunc("/login", handleLogin).Methods("POST")

  return s
}

