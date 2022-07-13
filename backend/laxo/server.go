package laxo

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/hashicorp/hcl/v2/hclsimple"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/mediocregopher/radix/v4"
	"github.com/urfave/negroni"
	"laxo.vn/laxo/laxo/sqlc"
)

type Config struct {
	Port             string `hcl:"port"`
	GRPCPort         string `hcl:"grpc_port"`
	LogLevel         string `hcl:"log_level"`
	AuthCookieName   string `hcl:"auth_cookie_name"`
	AuthCookieExpire int    `hcl:"auth_cookie_days_expire"`
	CallbackBasePath string `hcl:"callback_base_path"`
	MaxAssetSize     int    `hcl:"max_asset_size"`
}

type Server struct {
	Router      *mux.Router
	Negroni     *negroni.Negroni
	Logger      *Logger
	Config      *Config
	RedisClient radix.Client
	PglClient   *pgxpool.Pool
	Queries     *sqlc.Queries
	Middleware  *Middleware
}

func NewServer(l *Logger, c *Config) (*Server, error) {
	base := mux.NewRouter()
	r := base.PathPrefix("/api").Subrouter().StrictSlash(false)

	// Common middlewares
	var commonMiddlewares []negroni.Handler

	commonMiddlewares = append(commonMiddlewares, NewNegroniZapLogger(l))

	n := negroni.New(
		commonMiddlewares...,
	)

	s := &Server{
		Router:  r,
		Negroni: n,
		Logger:  l,
		Config:  c,
	}

	return s, nil
}

func (s *Server) InitMiddleware() {
	m := NewMiddleware(s)
	s.Middleware = &m
}

func (s *Server) ServeStaticFiles(path string) {
	s.Logger.Infow("Serving static files...", "path", path)

	//@TODO: DISABLE DIR LISTING!
	fileServer := http.FileServer(http.Dir(path))
	s.Router.PathPrefix("/assets/").Handler(http.StripPrefix("/api/assets/", fileServer))
}

func InitConfig() (*Config, error) {
	var appConfig Config

	err := hclsimple.DecodeFile("config.hcl", nil, &appConfig)
	if err != nil {
		return nil, err
	}

	return &appConfig, err
}
