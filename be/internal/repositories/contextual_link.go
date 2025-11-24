package repositories

import (
	"errors"
	"fmt"
	"time"

	"github.com/supabase-community/supabase-go"
)

type ContextualLinkRepository struct {
	Client *supabase.Client
}

type ContextualLink struct {
	ID                 string    `json:"id"`
	Link               string    `json:"link"`
	ContextualizedLink string    `json:"contextualized_link"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

func NewContextualLinkRepository(client *supabase.Client) *ContextualLinkRepository {
	return &ContextualLinkRepository{Client: client}
}

func (r *ContextualLinkRepository) FindByLink(link string) (*ContextualLink, error) {
	var data []ContextualLink

	count, err := r.Client.From("contextual-links").
		Select("*", "exact", false).
		Eq("link", link).
		ExecuteTo(&data)

	if err != nil {
		return nil, fmt.Errorf("failed to query link: %w", err)
	}

	if count == 0 {
		return nil, nil
	}

	return &data[0], nil
}

func (r *ContextualLinkRepository) FindByContextualizedLink(contextualizedLink string) (*ContextualLink, error) {
	var data []ContextualLink

	count, err := r.Client.From("contextual-links").
		Select("*", "exact", false).
		Eq("contextualized_link", contextualizedLink).
		ExecuteTo(&data)

	if err != nil {
		return nil, fmt.Errorf("failed to query contextualized link: %w", err)
	}

	if count == 0 {
		return nil, nil
	}

	return &data[0], nil
}

func (r *ContextualLinkRepository) Create(link, contextualizedLink string) error {
	if link == "" {
		return errors.New("link cannot be empty")
	}

	if contextualizedLink == "" {
		return errors.New("contextualizedLink cannot be empty")
	}

	// Don't set created_at and updated_at manually - let Supabase handle them with DEFAULT values
	_, count, err := r.Client.From("contextual-links").
		Insert(map[string]interface{}{
			"link":                link,
			"contextualized_link": contextualizedLink,
		}, false, "", "*", "exact").
		Execute()

	if err != nil {
		return fmt.Errorf("failed to insert link: %w", err)
	}

	if count == 0 {
		return errors.New("failed to create record: no data returned")
	}

	return nil
}

// Upsert inserts a new record or updates an existing one based on the link field
// If a record with the same link exists, it updates the contextualized_link
// Otherwise, it creates a new record
func (r *ContextualLinkRepository) Upsert(link, contextualizedLink string) error {
	if link == "" {
		return errors.New("link cannot be empty")
	}

	if contextualizedLink == "" {
		return errors.New("contextualizedLink cannot be empty")
	}

	// Upsert operation: insert if not exists, update if exists
	// The "link" parameter specifies which column to check for conflicts
	_, count, err := r.Client.From("contextual-links").
		Upsert(map[string]interface{}{
			"link":                link,
			"contextualized_link": contextualizedLink,
		}, "link", "*", "exact").
		Execute()

	if err != nil {
		return fmt.Errorf("failed to upsert link: %w", err)
	}

	if count == 0 {
		return errors.New("failed to upsert record: no data returned")
	}

	return nil
}
