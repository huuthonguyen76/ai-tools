package syncdatabase

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"

	logger "example.com/m/v2/internal/pkg"
	"example.com/m/v2/internal/repositories"
	"example.com/m/v2/internal/services"
	"github.com/supabase-community/supabase-go"
	"go.uber.org/zap"
)

// SyncService handles syncing data from Apify to Supabase
type SyncService struct {
	ApifyService            *services.ApifyService
	SupabaseClient          *supabase.Client
	SocialPostRawRepository *repositories.SocialPostRawRepository
}

// NewSyncService creates a new sync service instance
func NewSyncService(apifyService *services.ApifyService, supabaseClient *supabase.Client, socialPostRawRepository *repositories.SocialPostRawRepository) (*SyncService, error) {
	if apifyService == nil {
		return nil, fmt.Errorf("apify service cannot be nil")
	}
	if supabaseClient == nil {
		return nil, fmt.Errorf("supabase client cannot be nil")
	}

	return &SyncService{
		ApifyService:            apifyService,
		SupabaseClient:          supabaseClient,
		SocialPostRawRepository: socialPostRawRepository,
	}, nil
}

// SyncAllDatasets fetches all dataset IDs, gets all items, and inserts them into Supabase
func (s *SyncService) SyncAllDatasets(ctx context.Context) {
	channel := "apify_dataset"

	// Step 1: Get all dataset IDs
	logger.Info(ctx, "Fetching all dataset IDs from Apify")
	datasetIDs, err := s.ApifyService.GetAllDatasetIDs(ctx)
	if err != nil {
		logger.Error(ctx, "Failed to get dataset IDs", zap.Error(err))
	}

	logger.Info(ctx, "Found datasets", zap.Int("count", len(datasetIDs)))

	// Step 2: For each dataset, get all items and insert into Supabase
	for _, datasetID := range datasetIDs {
		logger.Info(ctx, "Processing dataset", zap.String("dataset_id", datasetID))

		// Get all items from this dataset
		items, err := s.ApifyService.GetAllRawDatasetItems(ctx, datasetID)
		if err != nil {
			errMsg := fmt.Sprintf("failed to get items from dataset %s: %v", datasetID, err)
			logger.Error(ctx, errMsg, zap.Error(err))
			continue
		}
		for _, item := range items {
			hashContent := md5.Sum([]byte(fmt.Sprintf("%v", item)))

			record := &repositories.SocialPostRaw{
				Channel:     channel,
				RawContent:  item,
				HashContent: hex.EncodeToString(hashContent[:]),
			}
			err = s.SocialPostRawRepository.Create(record)
			if err != nil {
				logger.Error(ctx, "Failed to upsert social post raw", zap.Error(err))
			}
		}
	}

	logger.Info(ctx, "Sync completed for dataset")
}
