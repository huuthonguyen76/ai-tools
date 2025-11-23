package contextualizelink

import (
	"errors"

	"example.com/m/v2/internal/repositories"
	"example.com/m/v2/internal/services"
)

type ContextualizeLinkHandler struct {
	DifyService              *services.DifyService
	ContextualLinkRepository *repositories.ContextualLinkRepository
}

func NewContextualizeLinkHandler(difyService *services.DifyService, contextualLinkRepository *repositories.ContextualLinkRepository) *ContextualizeLinkHandler {
	return &ContextualizeLinkHandler{
		DifyService:              difyService,
		ContextualLinkRepository: contextualLinkRepository,
	}
}

func (h *ContextualizeLinkHandler) Handler(link string) (string, error) {
	if link == "" {
		return "", errors.New("link query parameter is required")
	}

	contextualizedLink, err := h.DifyService.GetContextualLink(link)
	if err != nil {
		return "", err
	}

	err = h.ContextualLinkRepository.Upsert(link, contextualizedLink)
	if err != nil {
		return "", err
	}

	return contextualizedLink, nil
}
