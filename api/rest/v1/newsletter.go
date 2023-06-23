package v1

import (
	"fmt"
	"net/http"
	"newsletter-manager-go/api/rest/v1/model"
	domnewsletter "newsletter-manager-go/domain/newsletter"
	util "newsletter-manager-go/util"
	utilctx "newsletter-manager-go/util/context"
)

func (h *Handler) CreateNewsletter(_ http.ResponseWriter, r *http.Request, input model.CreateNewsletterReq) (*model.Newsletter, error) {
	authorID, _ := utilctx.AuthorIDFromCtx(r.Context())

	createNewsletterInput := domnewsletter.CreateNewsletterInput{
		AuthorID:    authorID,
		Name:        input.Name,
		Description: input.Description,
	}

	newsletter, err := h.newsletterService.Create(r.Context(), createNewsletterInput)
	if err != nil {
		return nil, fmt.Errorf("creating newsletter: %w", err)
	}

	modelNewsletter := model.FromDomainNewsletter(newsletter)

	return &modelNewsletter, nil
}

func (h *Handler) ListNewsletters(_ http.ResponseWriter, r *http.Request) ([]model.Newsletter, error) {
	_, ok := utilctx.AuthorIDFromCtx(r.Context())
	if !ok {
		return nil, fmt.Errorf("author id not found")
	}

	newsletters, err := h.newsletterService.List(r.Context())
	if err != nil {
		return nil, err
	}

	modelNewsletters := util.MapFuncRef(newsletters, model.FromDomainNewsletter)

	return modelNewsletters, nil
}

func (h *Handler) GetCurrentAuthorNewsletters(_ http.ResponseWriter, r *http.Request) ([]model.Newsletter, error) {
	authorID, _ := utilctx.AuthorIDFromCtx(r.Context())

	newsletters, err := h.newsletterService.ListCurrentAuthorNewsletters(r.Context(), authorID)
	if err != nil {
		return nil, err
	}

	modelNewsletters := util.MapFuncRef(newsletters, model.FromDomainNewsletter)

	return modelNewsletters, nil
}

func (h *Handler) GetNewsletter(_ http.ResponseWriter, r *http.Request, input model.PathNewsletterInput) (*model.Newsletter, error) {
	newsletter, err := h.newsletterService.Read(r.Context(), input.NewsletterID)
	if err != nil {
		return nil, err
	}

	modelNewsletter := model.FromDomainNewsletter(newsletter)

	return &modelNewsletter, nil
}

func (h *Handler) UpdateNewsletter(_ http.ResponseWriter, r *http.Request, input model.PathNewsletterInput) (*model.Newsletter, error) {
	// TODO
	return nil, nil
}

func (h *Handler) DeleteNewsletter(_ http.ResponseWriter, r *http.Request, input model.PathNewsletterInput) error {
	// TODO
	return nil
}
