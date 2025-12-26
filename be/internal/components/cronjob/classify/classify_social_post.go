package classify

import (
	"context"
	"encoding/json"
	"fmt"

	"example.com/m/v2/internal/repositories"
	"github.com/bytedance/gopkg/util/logger"
	"go.uber.org/zap"

	b "example.com/m/v2/baml_client"
)

type ClassifySocialPost struct {
	SocialPostRawRepository *repositories.SocialPostRawRepository
}

func NewClassifySocialPost(socialPostRawRepository *repositories.SocialPostRawRepository) *ClassifySocialPost {
	return &ClassifySocialPost{
		SocialPostRawRepository: socialPostRawRepository,
	}
}

func (c *ClassifySocialPost) ClassifySocialPost(ctx context.Context) error {
	socialPosts, err := c.SocialPostRawRepository.FindAll(ctx)
	if err != nil {
		return err
	}

	logger.Info(ctx, "Found social posts", zap.Int("count", len(socialPosts)))

	for _, rawSocialPost := range socialPosts {
		socialPost, err := ConvertToApifyFacebookGroup(rawSocialPost.RawContent)
		if err != nil {
			logger.Error(ctx, "Failed to convert raw content", zap.Error(err))
			continue
		}

		result, err := b.ClassifySocialPost(ctx, socialPost.Text)
		if err != nil {
			logger.Error(ctx, "Failed to classify social post", zap.Error(err))
			return err
		}
		logger.Info(ctx, "Classified social post", zap.Any("result", result))
	}

	return nil
}

// convertToApifyFacebookGroup converts a map[string]interface{} to ApifyFacebookGroup struct
func ConvertToApifyFacebookGroup(rawContent map[string]interface{}) (repositories.ApifyFacebookGroup, error) {
	var result repositories.ApifyFacebookGroup

	// Marshal map to JSON bytes, then unmarshal into struct
	jsonBytes, err := json.Marshal(rawContent)
	if err != nil {
		return result, fmt.Errorf("failed to marshal raw content: %w", err)
	}

	if err := json.Unmarshal(jsonBytes, &result); err != nil {
		return result, fmt.Errorf("failed to unmarshal to ApifyFacebookGroup: %w", err)
	}

	return result, nil
}
