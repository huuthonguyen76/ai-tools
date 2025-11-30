package contextualizelink

import (
	"errors"
	"strings"

	"example.com/m/v2/internal/repositories"
	"example.com/m/v2/internal/services"
)

type ContextualizeLinkHandler struct {
	DifyService              services.DifyServiceInterface
	ContextualLinkRepository repositories.ContextualLinkRepositoryInterface
}

func NewContextualizeLinkHandler(difyService services.DifyServiceInterface, contextualLinkRepository repositories.ContextualLinkRepositoryInterface) *ContextualizeLinkHandler {
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
	if err != nil || contextualizedLink == "" {
		if contextualizedLink == "" {
			return "", errors.New("contextualized link is empty")
		}
		return "", err
	}

	contextualizedLink = strings.TrimPrefix(contextualizedLink, "https://")
	contextualizedLink = strings.TrimPrefix(contextualizedLink, "http://")
	contextualizedLink = strings.TrimPrefix(contextualizedLink, "/")
	contextualizedLink = strings.TrimSuffix(contextualizedLink, "/")
	contextualizedLink = strings.TrimSpace(contextualizedLink)
	contextualizedLink = strings.ToLower(contextualizedLink)

	err = h.ContextualLinkRepository.Upsert(link, contextualizedLink)
	if err != nil {
		return "", err
	}

	return contextualizedLink, nil
}
