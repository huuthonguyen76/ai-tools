package contextualizelink

import (
	"testing"

	mocks "example.com/m/v2/internal/mocks"
	"github.com/stretchr/testify/assert"
)

func TestGetContextualLink(t *testing.T) {
	url := "https://example.com/test"
	expectedLink := "example.com/test-contextual"

	mockDifyService := mocks.NewMockDifyServiceInterface(t)
	contextualRepo := mocks.NewMockContextualLinkRepositoryInterface(t)

	contextualRepo.EXPECT().Upsert(url, expectedLink).Return(nil)
	mockDifyService.EXPECT().GetContextualLink(url).Return(expectedLink, nil)

	handler := NewContextualizeLinkHandler(mockDifyService, contextualRepo)

	result, err := handler.Handler(url)

	assert.NoError(t, err)
	assert.Equal(t, result, expectedLink)

	assert.Equal(t, true, true)

}
