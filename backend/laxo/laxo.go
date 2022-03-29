package laxo

import (
	"log"

	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/hcl/v2/hclsimple"
	"github.com/urfave/negroni"
)

type Config struct {
  Port             string `hcl:"port"`
  LogLevel         string `hcl:"log_level"`
  AuthCookieName   string `hcl:"auth_cookie_name"`
  AuthCookieExpire int    `hcl:"auth_cookie_days_expire"`
}

var (
  Logger    hclog.Logger
  AppConfig Config
)

func InitConfig(testing bool) (hclog.Logger, Config) {
  err := hclsimple.DecodeFile("config.hcl", nil, &AppConfig)
  if err != nil {
    log.Fatalf("Failed to load configuration: %r", err)
  }

  logLevel := hclog.Off

  if(!testing) {
    logLevel = hclog.LevelFromString(AppConfig.LogLevel)
  }

  appLogger := hclog.New(&hclog.LoggerOptions{
    Name:  "laxo-backend",
    Level: logLevel,
  })

  Logger = appLogger

  appLogger.Info("Configuration loaded", "LogLevel", AppConfig.LogLevel)

  return appLogger, AppConfig
}

func SetupRouter(testing bool) *mux.Router {
  router := mux.NewRouter()

  // Common middlewares
  var commonMiddlewares []negroni.Handler

  if(!testing) {
    commonMiddlewares = append(commonMiddlewares, negroni.NewLogger())
  }

  common := negroni.New(
    commonMiddlewares...
  )

  subRouter := router.PathPrefix("/api").Subrouter().StrictSlash(true)

	subRouter.Handle("/shop", common.With(
		negroni.HandlerFunc(assureJSON),
		negroni.WrapFunc(assureAuth(HandleCreateShop)),
	)).Methods("POST")

	subRouter.Handle("/user", common.With(
		negroni.HandlerFunc(assureJSON),
		negroni.WrapFunc(HandleCreateUser),
	)).Methods("POST")

	subRouter.Handle("/login", common.With(
		negroni.HandlerFunc(assureJSON),
		negroni.WrapFunc(HandleLogin),
	)).Methods("POST")

	subRouter.Handle("/logout", common.With(
		negroni.WrapFunc(HandleLogout),
	)).Methods("POST")

	subRouter.Handle("/user", common.With(
		negroni.WrapFunc(assureAuth(HandleGetUser)),
	)).Methods("GET")


  return subRouter
}

