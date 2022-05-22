package laxo

import (
	"github.com/gorilla/mux"
	"github.com/hashicorp/hcl/v2/hclsimple"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/mediocregopher/radix/v4"
	"github.com/urfave/negroni"
	"laxo.vn/laxo/laxo/sqlc"
)

type Config struct {
  Port             string `hcl:"port"`
  LogLevel         string `hcl:"log_level"`
  AuthCookieName   string `hcl:"auth_cookie_name"`
  AuthCookieExpire int    `hcl:"auth_cookie_days_expire"`
  CallbackBasePath string `hcl:"callback_base_path"`
}

type Server struct {
  Router       *mux.Router
  Negroni      *negroni.Negroni
  Logger       *Logger
  Config       *Config
  RedisClient  radix.Client
  PglClient    *pgxpool.Pool
  Queries      *sqlc.Queries
}

func NewServer(l *Logger, c *Config) (*Server, error) {
  r := mux.NewRouter().PathPrefix("/api").Subrouter().StrictSlash(true)

  // Common middlewares
  var commonMiddlewares []negroni.Handler

  commonMiddlewares = append(commonMiddlewares, negroni.NewLogger())

  n := negroni.New(
    commonMiddlewares...
  )

  s := &Server{
    Router: r,
    Negroni: n,
    Logger: l,
    Config: c,
  }

  return s, nil
}

func InitConfig() (*Config, error) {
  var appConfig Config

  err := hclsimple.DecodeFile("config.hcl", nil, &appConfig)
  if err != nil {
    return nil, err
  }

  return &appConfig, err
}

