package graph

import (
	_ "embed"
	"errors"
	"net/http"

	"strv-template-backend-go-api/api/graphql/middleware"
	"strv-template-backend-go-api/database/sql"
	"strv-template-backend-go-api/util"
	"strv-template-backend-go-api/util/timesource"

	httpx "go.strv.io/net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

const fixedComplexityLimit = 15

// Controller handles all /api endpoints.
// It is responsible for routing requests to appropriate handlers.
type Controller struct {
	*chi.Mux

	userService    UserService
	sessionService SessionService
	tokenParser    middleware.TokenParser
	corsConfig     middleware.CORSConfig
	logger         *zap.Logger
}

// NewController returns new instance of HTTP Graphql controller.
func NewController(
	userService UserService,
	sessionService SessionService,
	tokenParser middleware.TokenParser,
	dataSource sql.DataSource,
	corsConfig middleware.CORSConfig,
	logger *zap.Logger,
) (*Controller, error) {
	if err := newControllerValidate(userService, sessionService, tokenParser, dataSource, logger); err != nil {
		return nil, err
	}
	controller := &Controller{
		userService:    userService,
		sessionService: sessionService,
		tokenParser:    tokenParser,
		corsConfig:     corsConfig,
		logger:         logger,
	}
	controller.initRouter(dataSource)
	return controller, nil
}

// Initialize router for controller.
func (c *Controller) initRouter(dataSource sql.DataSource) {
	r := chi.NewRouter()

	r.Use(middleware.NewCORSHandler(c.corsConfig))
	r.Use(middleware.RequestStartTime(timesource.DefaultTimeSource{}))
	r.Use(httpx.RequestIDMiddleware(func(h http.Header) string {
		return h.Get(httpx.Header.XRequestID)
	}))
	r.Use(httpx.RecoverMiddleware(util.NewServerLogger("httpx.RecoverMiddleware")))
	r.Use(middleware.LimitBodySize(c.logger, middleware.DefaultByteCountLimit))
	r.Use(middleware.DataLoader(dataSource))

	authenticate := middleware.Authenticate(c.logger, c.tokenParser)

	gqlServer := handler.NewDefaultServer(NewExecutableSchema(Config{
		Resolvers:  NewResolver(c.userService, c.sessionService),
		Directives: NewDirectiveHandler(true),
	}))
	gqlServer.SetErrorPresenter(ErrorPresenter(c.logger))
	gqlServer.Use(extension.FixedComplexityLimit(fixedComplexityLimit))

	r.With(authenticate).Handle("/api/graphql", gqlServer)

	r.Get("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	c.Mux = r
}

func newControllerValidate(
	userService UserService,
	sessionService SessionService,
	tokenParser middleware.TokenParser,
	dataSource sql.DataSource,
	logger *zap.Logger,
) error {
	if userService == nil {
		return errors.New("invalid user service")
	}
	if sessionService == nil {
		return errors.New("invalid session service")
	}
	if tokenParser == nil {
		return errors.New("invalid token parser")
	}
	if dataSource == nil {
		return errors.New("invalid data source")
	}
	if logger == nil {
		return errors.New("invalid logger")
	}
	return nil
}
