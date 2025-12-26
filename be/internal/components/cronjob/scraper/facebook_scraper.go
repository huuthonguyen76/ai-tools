package scraper

import (
	"context"

	"example.com/m/v2/internal/services"
)

type FacebookScraper struct {
	ApifyService *services.ApifyService
}

func NewFacebookScraper(apifyService *services.ApifyService) *FacebookScraper {
	return &FacebookScraper{
		ApifyService: apifyService,
	}
}

func ScrapeGroup(ctx context.Context, groupURL string) error {
	return nil
}
