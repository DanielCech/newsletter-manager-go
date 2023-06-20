package v1

import (
	"fmt"
	"net/http"
	"newsletter-manager-go/api/rest/v1/model"
	domnewsletter "newsletter-manager-go/domain/newsletter"
	"newsletter-manager-go/types/id"
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
	// TODO
	return nil, nil
}

func (h *Handler) GetCurrentAuthorNewsletters(_ http.ResponseWriter, r *http.Request, input model.AuthorIDInput) ([]model.Newsletter, error) {
	authorID, _ := utilctx.AuthorIDFromCtx(r.Context())

	newsletters, err := h.newsletterService.ListCurrentAuthorNewsletters(r.Context(), authorID)
	if err != nil {
		return nil, err
	}

	modelNewsletters := util.MapFuncRef(newsletters, model.FromDomainNewsletter)

	return modelNewsletters, nil
}

func (h *Handler) GetNewsletter(_ http.ResponseWriter, r *http.Request, newsletterId id.Newsletter) (*model.Newsletter, error) {
	newsletter, err := h.newsletterService.Read(r.Context(), newsletterId)
	if err != nil {
		return nil, err
	}

	modelNewsletter := model.FromDomainNewsletter(newsletter)

	return &modelNewsletter, nil
}

func (h *Handler) UpdateNewsletter(_ http.ResponseWriter, r *http.Request, input model.NewsletterIDInput) (*model.Newsletter, error) {
	// TODO
	return nil, nil
}

func (h *Handler) DeleteNewsletter(_ http.ResponseWriter, r *http.Request, input model.NewsletterIDInput) error {
	// TODO
	return nil
}
