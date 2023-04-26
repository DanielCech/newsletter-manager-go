package v1

import (
	"net/http"

	"newsletter-manager-go/api/rest/middleware"
	httputil "newsletter-manager-go/api/rest/util"
	domuser "newsletter-manager-go/domain/user"

	"go.strv.io/net/http/signature"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

// Handler for v1 endpoints.
type Handler struct {
	*chi.Mux

	userService    UserService
	sessionService SessionService
	tokenParser    middleware.TokenParser
	logger         *zap.Logger
}

// NewHandler returns new instance of handler handling /v1 endpoints.
func NewHandler(
	userService UserService,
	sessionService SessionService,
	tokenParser middleware.TokenParser,
	logger *zap.Logger,
) *Handler {
	h := &Handler{
		userService:    userService,
		sessionService: sessionService,
		tokenParser:    tokenParser,
		logger:         logger,
	}
	h.initRouter()
	return h
}

// initRouter initializes chi router for the handler.
func (h *Handler) initRouter() {
	r := chi.NewRouter()

	authenticate := middleware.Authenticate(h.logger, h.tokenParser)
	authorizeAdmin := middleware.Authorize(h.logger, domuser.RoleAdmin)
	authorizeUser := middleware.Authorize(h.logger, domuser.RoleUser)

	w := signature.DefaultWrapper().
		WithInputGetter(httputil.ParseRequestBody).
		WithErrorHandler(func(w http.ResponseWriter, r *http.Request, err error) {
			httputil.WriteErrorResponse(r.Context(), h.logger, w, err)
		})
	wCreated := w.WithResponseMarshaler(signature.FixedResponseCodeMarshal(http.StatusCreated))

	r.Route("/users", func(r chi.Router) {
		r.Route("/register", func(r chi.Router) {
			r.Post("/", signature.WrapHandler(wCreated, h.CreateUser))
		})
		r.Group(func(r chi.Router) {
			r.Use(authenticate)
			r.Group(func(r chi.Router) {
				r.Use(authorizeUser)
				r.Route("/me", func(r chi.Router) {
					r.Get("/", signature.WrapHandlerResponse(w, h.ReadLoggedUser))
				})
				r.Route("/change-password", func(r chi.Router) {
					r.Patch("/", signature.WrapHandlerInput(w, h.ChangeUserPassword))
				})
			})
			r.Group(func(r chi.Router) {
				r.Use(authorizeAdmin)
				r.Get("/", signature.WrapHandlerResponse(w, h.ListUsers))
			})
		})
	})

	r.Route("/sessions", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Post("/native", signature.WrapHandler(wCreated, h.CreateSession))
			r.Post("/refresh", signature.WrapHandler(wCreated, h.RefreshSession))
		})
		r.Post("/destroy", signature.WrapHandlerInput(w, h.DestroySession))
	})

	h.Mux = r
}
