package repositories

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/supabase-community/supabase-go"
)

const TABLE_NAME = "social_post_raw"

type SocialPostRawRepository struct {
	Client *supabase.Client
}

type SocialPostRawRepositoryInterface interface {
	FindByID(id string) (*SocialPostRaw, error)
	FindByURL(url string) (*SocialPostRaw, error)
	FindByChannel(channel string) ([]*SocialPostRaw, error)
	FindBySocialMonitorID(socialMonitorID string) ([]*SocialPostRaw, error)
	Create(post *SocialPostRaw) error
	Upsert(post *SocialPostRaw) error
}

type SocialPostRaw struct {
	ID          string                 `json:"id"`
	Channel     string                 `json:"channel"`
	RawContent  map[string]interface{} `json:"raw_content"`
	HashContent string                 `json:"hash_content"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}

type ApifyFacebookGroup struct {
	FacebookGroupURL string `json:"inputUrl"`
	Text             string `json:"text"`
	PermaURL         string `json:"url"`
	Attachments      []struct {
		Image struct {
			Height int    `json:"height"`
			Width  int    `json:"width"`
			URI    string `json:"uri"`
		} `json:"image"`
	} `json:"attachments"`
}

func NewSocialPostRawRepository(client *supabase.Client) *SocialPostRawRepository {
	return &SocialPostRawRepository{Client: client}
}

func (r *SocialPostRawRepository) FindByID(id string) (*SocialPostRaw, error) {
	var data []SocialPostRaw

	count, err := r.Client.From(TABLE_NAME).
		Select("*", "exact", false).
		Eq("id", id).
		ExecuteTo(&data)

	if err != nil {
		return nil, fmt.Errorf("failed to query social post raw by id: %w", err)
	}

	if count == 0 {
		return nil, nil
	}

	return &data[0], nil
}

func (r *SocialPostRawRepository) Create(post *SocialPostRaw) error {
	if post == nil {
		return errors.New("social post raw cannot be nil")
	}

	if post.Channel == "" {
		return errors.New("channel cannot be empty")
	}

	// Don't set created_at and updated_at manually - let Supabase handle them with DEFAULT values
	_, count, err := r.Client.From(TABLE_NAME).
		Insert(map[string]interface{}{
			"id":           uuid.NewString(),
			"channel":      post.Channel,
			"raw_content":  post.RawContent,
			"hash_content": post.HashContent,
		}, false, "", "*", "exact").
		Execute()

	if err != nil {
		return fmt.Errorf("failed to insert social post: %w", err)
	}

	if count == 0 {
		return errors.New("failed to create social post: no data returned")
	}

	return nil
}

func (r *SocialPostRawRepository) FindAll(ctx context.Context) ([]SocialPostRaw, error) {
	var posts []SocialPostRaw

	count, err := r.Client.From(TABLE_NAME).
		Select("*", "exact", false).
		ExecuteTo(&posts)

	if err != nil {
		return nil, fmt.Errorf("failed to get all social post raws: %w", err)
	}

	if count == 0 {
		return []SocialPostRaw{}, nil
	}

	return posts, nil
}
