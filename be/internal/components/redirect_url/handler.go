package redirecturl

import (
	"errors"

	"example.com/m/v2/internal/repositories"
)

type RedirectURLHandler struct {
	ContextualLinkRepository *repositories.ContextualLinkRepository
}

func NewRedirectURLHandler(contextualLinkRepository *repositories.ContextualLinkRepository) *RedirectURLHandler {
	return &RedirectURLHandler{
		ContextualLinkRepository: contextualLinkRepository,
	}
}

func (h *RedirectURLHandler) Handler(contextualizedLink string) (string, error) {
	if contextualizedLink == "" {
		return "", errors.New("contextualizedLink parameter is required")
	}

	record, err := h.ContextualLinkRepository.FindByContextualizedLink(contextualizedLink)
	if err != nil {
		return "", err
	}

	if record == nil {
		return "", errors.New("contextualized link not found")
	}

	return record.Link, nil
}
