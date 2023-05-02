package v1

import (
	"net/http"

	"go.strv.io/net/http/signature"
	"newsletter-manager-go/api/rest/middleware"
	httputil "newsletter-manager-go/api/rest/util"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

// Handler for v1 endpoints.
type Handler struct {
	*chi.Mux

	authorService     AuthorService
	newsletterService NewsletterService
	tokenParser       middleware.TokenParser
	logger            *zap.Logger
}

// NewHandler returns new instance of handler handling /v1 endpoints.
func NewHandler(
	authorService AuthorService,
	newsletterService NewsletterService,
	logger *zap.Logger,
) *Handler {
	h := &Handler{
		authorService:     authorService,
		newsletterService: newsletterService,
		logger:            logger,
	}
	h.initRouter()
	return h
}

// initRouter initializes chi router for the handler.
func (h *Handler) initRouter() {
	r := chi.NewRouter()

	//authenticate := middleware.Authenticate(h.logger, h.tokenParser)

	w := signature.DefaultWrapper().
		WithInputGetter(httputil.ParseRequestBody).
		WithErrorHandler(func(w http.ResponseWriter, r *http.Request, err error) {
			httputil.WriteErrorResponse(r.Context(), h.logger, w, err)
		})
	wCreated := w.WithResponseMarshaler(signature.FixedResponseCodeMarshal(http.StatusCreated))

	r.Route("/authors", func(r chi.Router) {
		r.Get("/", signature.WrapHandlerResponse(w, h.ListAuthors))

		r.Route("/{authorId}", func(r chi.Router) {
			//r.Use(authenticate)
			r.Get("/", signature.WrapHandler(wCreated, h.GetAuthor))
			r.Patch("/", signature.WrapHandler(wCreated, h.UpdateAuthor))
			r.Delete("/", signature.WrapHandlerInput(wCreated, h.DeleteAuthor))
			r.Get("/newsletters", signature.WrapHandler(wCreated, h.GetAuthorNewsletters))
			r.Post("/newsletters", signature.WrapHandler(wCreated, h.CreateNewsletter))
		})

		r.Post("/sign-in", signature.WrapHandler(wCreated, h.AuthorSignIn))
		r.Post("/sign-up", signature.WrapHandler(wCreated, h.AuthorSignUp))
		r.Get("/subscriptions", signature.WrapHandler(wCreated, h.AuthorSubscriptions))

		// TODO: maybe handled by Update
		// r.Post("/change-password", signature.WrapHandler(wCreated, h.ChangePassword))

		r.Post("/refresh-token", signature.WrapHandler(wCreated, h.RefreshToken))
	})

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
