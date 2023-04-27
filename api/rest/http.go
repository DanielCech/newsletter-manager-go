package http

import (
	_ "embed"
	"errors"
	"fmt"
	"net/http"

	"newsletter-manager-go/api/rest/middleware"
	httputil "newsletter-manager-go/api/rest/util"
	httpv1 "newsletter-manager-go/api/rest/v1"
	"newsletter-manager-go/util"

	httpx "go.strv.io/net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

//go:generate docker compose -f docker-compose-swagger-gen.yaml run -T --rm swagger-gen
//go:embed openapi.yaml
var OpenAPI []byte

// Controller handles all /api endpoints.
// It is responsible for routing requests to appropriate handlers.
// Versioned endpoints are handled by subcontrollers.
type Controller struct {
	*chi.Mux

	authorService     httpv1.AuthorService
	newsletterService httpv1.NewsletterService
	tokenParser       middleware.TokenParser
	corsConfig        middleware.CORSConfig
	logger            *zap.Logger
}

// NewController returns new instance of a HTTP REST controller.
func NewController(
	authorService httpv1.AuthorService,
	newsletterService httpv1.NewsletterService,
	tokenParser middleware.TokenParser,
	logger *zap.Logger,
) (*Controller, error) {
	if err := newControllerValidate(authorService, newsletterService, tokenParser, logger); err != nil {
		return nil, err
	}
	controller := &Controller{
		authorService: authorService,
		tokenParser:   tokenParser,
		logger:        logger,
	}
	controller.initRouter()
	return controller, nil
}

// initRouter initializes chi router for the controller.
func (c *Controller) initRouter() {
	r := chi.NewRouter()

	r.Use(middleware.NewCORSHandler(c.corsConfig))
	r.Use(httpx.RequestIDMiddleware(func(h http.Header) string {
		return h.Get(httpx.Header.XRequestID)
	}))
	r.Use(httpx.LoggingMiddleware(util.NewServerLogger("httpx.LoggingMiddleware")))
	r.Use(httpx.RecoverMiddleware(util.NewServerLogger("httpx.RecoverMiddleware")))
	r.Use(middleware.LimitBodySize(c.logger, middleware.DefaultByteCountLimit))

	authenticate := middleware.Authenticate(c.logger, c.tokenParser)

	v1Handler := httpv1.NewHandler(c.authorService, c.newsletterService, c.tokenParser, c.logger)

	r.Route("/api", func(r chi.Router) {
		r.With(authenticate).Get("/openapi.yaml", c.OpenAPI)
		r.Mount("/v1", v1Handler)
	})

	r.Get("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	c.Mux = r
}

// OpenAPI serves rendered OpenAPI file.
func (c *Controller) OpenAPI(w http.ResponseWriter, r *http.Request) {
	encodeFunc := func(w http.ResponseWriter, data any) error {
		d, ok := data.([]byte)
		if !ok {
			return fmt.Errorf("expected byte slice: got %T", data)
		}
		if _, err := w.Write(d); err != nil {
			return fmt.Errorf("writing openapi content: %w", err)
		}
		return nil
	}
	httputil.WriteResponse(
		util.WithCtx(r.Context(), c.logger),
		w,
		OpenAPI,
		http.StatusOK,
		httpx.WithEncodeFunc(encodeFunc),
		httpx.WithContentType(httpx.TextYAML),
	)
}

func newControllerValidate(
	authorService httpv1.AuthorService,
	newsletterService httpv1.NewsletterService,
	tokenParser middleware.TokenParser,
	logger *zap.Logger,
) error {
	if authorService == nil {
		return errors.New("invalid user service")
	}
	if newsletterService == nil {
		return errors.New("invalid newsletter service")
	}
	if tokenParser == nil {
		return errors.New("invalid token parser")
	}
	if logger == nil {
		return errors.New("invalid logger")
	}
	return nil
}
