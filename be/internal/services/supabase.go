package services

import (
	"errors"

	"github.com/supabase-community/supabase-go"
)

func NewSupabaseClient(API_URL, API_KEY string) (*supabase.Client, error) {
	client, err := supabase.NewClient(API_URL, API_KEY, &supabase.ClientOptions{})
	if err != nil {
		return nil, errors.New("failed to initialize the client")
	}
	return client, nil
}
