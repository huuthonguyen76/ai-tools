package contextualizelink

import (
	"errors"

	"example.com/m/v2/internal/services"
)

type ContextualizeLinkHandler struct {
	DifyService *services.DifyService
}

func NewContextualizeLinkHandler(difyService *services.DifyService) *ContextualizeLinkHandler {
	return &ContextualizeLinkHandler{DifyService: difyService}
}

func (h *ContextualizeLinkHandler) Handler(link string) (string, error) {
	if link == "" {
		return "", errors.New("link query parameter is required")
	}

	contextualizedLink, err := h.DifyService.GetContextualLink(link)
	if err != nil {
		return "", err
	}

	return contextualizedLink, nil
}
