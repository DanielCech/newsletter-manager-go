package v1

import (
	"net/http"

	"newsletter-manager-go/api/rest/middleware"
	httputil "newsletter-manager-go/api/rest/util"

	"github.com/go-chi/chi/v5"
	"go.strv.io/net/http/signature"
	"go.uber.org/zap"
)

// Handler for v1 endpoints.
type Handler struct {
	*chi.Mux

	authorService     AuthorService
	sessionService    SessionService
	newsletterService NewsletterService
	tokenParser       middleware.TokenParser
	logger            *zap.Logger
}

// NewHandler returns new instance of handler handling /v1 endpoints.
func NewHandler(
	authorService AuthorService,
	sessionService SessionService,
	newsletterService NewsletterService,
	tokenParser middleware.TokenParser,
	logger *zap.Logger,
) *Handler {
	h := &Handler{
		authorService:     authorService,
		sessionService:    sessionService,
		newsletterService: newsletterService,
		tokenParser:       tokenParser,
		logger:            logger,
	}
	h.initRouter()
	return h
}

// initRouter initializes chi router for the handler.
func (h *Handler) initRouter() {
	r := chi.NewRouter()

	authenticate := middleware.Authenticate(h.logger, h.tokenParser)

	w := signature.DefaultWrapper().
		WithInputGetter(httputil.ParseRequestBody).
		WithErrorHandler(func(w http.ResponseWriter, r *http.Request, err error) {
			httputil.WriteErrorResponse(r.Context(), h.logger, w, err)
		})
	wCreated := w.WithResponseMarshaler(signature.FixedResponseCodeMarshal(http.StatusCreated))

	r.Route("/authors", func(r chi.Router) {

		r.Route("/sign-up", func(r chi.Router) {
			r.Post("/", signature.WrapHandler(wCreated, h.AuthorSignUp))
		})
		r.Route("/sign-in", func(r chi.Router) {
			r.Post("/", signature.WrapHandler(wCreated, h.CreateSession))
		})

		r.Route("/current", func(r chi.Router) {
			r.Use(authenticate)

			r.Get("/", signature.WrapHandlerResponse(w, h.ReadCurrentAuthor))
			r.Patch("/", signature.WrapHandler(wCreated, h.UpdateCurrentAuthor))
			r.Delete("/", signature.WrapHandlerError(w, h.DeleteCurrentAuthor))

			r.Post("/change-password", signature.WrapHandlerInput(wCreated, h.ChangeAuthorPassword))

			r.Route("/refresh-token", func(r chi.Router) {
				r.Post("/", signature.WrapHandler(wCreated, h.RefreshSession))
			})

			r.Route("/logout", func(r chi.Router) {
				r.Post("/", signature.WrapHandlerInput(w, h.DestroySession))
			})

			r.Route("/newsletters", func(r chi.Router) {
				r.Get("/", signature.WrapHandler(wCreated, h.GetCurrentAuthorNewsletters))
				r.Post("/", signature.WrapHandler(wCreated, h.CreateNewsletter))
			})
		})
	})

	r.Get("/subscriptions", signature.WrapHandler(wCreated, h.AuthorSubscriptions))

	r.Route("/newsletters", func(r chi.Router) {
		// TODO: add to OpenAPI
		r.Get("/", signature.WrapHandlerResponse(w, h.ListNewsletters))

		r.Route("/{newsletterId}", func(r chi.Router) {
			r.Get("/", signature.WrapHandler(w, h.GetNewsletter))
			r.Patch("/", signature.WrapHandler(w, h.UpdateNewsletter))
			r.Delete("/", signature.WrapHandlerInput(w, h.DeleteNewsletter))
			r.Route("/emails", func(r chi.Router) {
				r.Get("/", signature.WrapHandler(w, h.ListNewsletterEmails))
				r.Post("/", signature.WrapHandler(w, h.CreateNewsletterEmail))
			})

			r.Post("/subscribe", signature.WrapHandlerResponse(w, h.SubscribeToNewsletter))
			r.Post("/unsubscribe", signature.WrapHandlerResponse(w, h.UnsubscribeFromNewsletter))
		})
	})

	h.Mux = r
}
