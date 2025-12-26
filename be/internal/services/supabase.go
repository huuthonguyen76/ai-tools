package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	logger "example.com/m/v2/internal/pkg"
	"github.com/supabase-community/supabase-go"
	"go.uber.org/zap"
)

// SupabaseClientInterface defines the interface for Supabase client operations
type SupabaseClientInterface interface {
	NewClient(API_URL, API_KEY string) (*supabase.Client, error)
}

// VectorRecord represents a record in the vector database
type VectorRecord struct {
	ID        string                 `json:"id,omitempty"`
	Content   string                 `json:"content"`
	Embedding []float64              `json:"embedding"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt string                 `json:"created_at,omitempty"`
	UpdatedAt string                 `json:"updated_at,omitempty"`
}

// VectorSearchResult represents a result from vector similarity search
type VectorSearchResult struct {
	ID         string                 `json:"id"`
	Content    string                 `json:"content"`
	Embedding  []float64              `json:"embedding,omitempty"`
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
	Similarity float64                `json:"similarity"` // Cosine similarity score
	Distance   float64                `json:"distance"`   // Vector distance
}

// VectorSearchOptions defines options for vector search operations
type VectorSearchOptions struct {
	Limit            int                    `json:"limit"`             // Maximum number of results to return
	Threshold        float64                `json:"threshold"`         // Minimum similarity threshold (0-1)
	IncludeEmbedding bool                   `json:"include_embedding"` // Whether to include embedding in results
	MetadataFilter   map[string]interface{} `json:"metadata_filter"`   // Filter by metadata
	SelectFields     []string               `json:"select_fields"`     // Specific fields to return
}

// VectorDatabaseService provides vector database operations with Supabase
type VectorDatabaseService struct {
	Client *supabase.Client
}

// NewSupabaseClient creates a new Supabase client with the provided credentials
func NewSupabaseClient(API_URL, API_KEY string) (*supabase.Client, error) {
	client, err := supabase.NewClient(API_URL, API_KEY, &supabase.ClientOptions{})
	if err != nil {
		return nil, errors.New("failed to initialize the client")
	}
	return client, nil
}

// NewVectorDatabaseService creates a new vector database service instance
func NewVectorDatabaseService(client *supabase.Client) (*VectorDatabaseService, error) {
	if client == nil {
		return nil, errors.New("supabase client cannot be nil")
	}
	return &VectorDatabaseService{Client: client}, nil
}

// InsertToVectorDatabase inserts a single record into the vector database
// tableName: the name of the table to insert into
// record: the data/embedding to insert; should be a map of column names to values
// Returns an error if insertion fails
func InsertToVectorDatabase(client *supabase.Client, tableName string, record map[string]interface{}) error {
	if client == nil {
		return errors.New("supabase client is nil")
	}
	if tableName == "" {
		return errors.New("table name cannot be empty")
	}
	if len(record) == 0 {
		return errors.New("record is empty")
	}

	_, _, err := client.From(tableName).Insert(record, false, "", "", "").Execute()
	if err != nil {
		return fmt.Errorf("failed to insert record into vector database: %w", err)
	}

	return nil
}

// Insert inserts a single vector record into the specified table
func (s *VectorDatabaseService) Insert(ctx context.Context, tableName string, record *VectorRecord) error {
	if tableName == "" {
		return errors.New("table name cannot be empty")
	}
	if record == nil {
		return errors.New("record cannot be nil")
	}
	if len(record.Embedding) == 0 {
		return errors.New("embedding vector cannot be empty")
	}

	recordMap := map[string]interface{}{
		"content":   record.Content,
		"embedding": record.Embedding,
	}

	if len(record.Metadata) > 0 {
		recordMap["metadata"] = record.Metadata
	}

	_, _, err := s.Client.From(tableName).Insert(recordMap, false, "", "", "").Execute()
	if err != nil {
		logger.Error(ctx, "Failed to insert vector record", zap.Error(err), zap.String("table", tableName))
		return fmt.Errorf("failed to insert vector record: %w", err)
	}

	logger.Info(ctx, "Vector record inserted successfully", zap.String("table", tableName))
	return nil
}

// BatchInsert inserts multiple vector records in a single operation
func (s *VectorDatabaseService) BatchInsert(ctx context.Context, tableName string, records []*VectorRecord) error {
	if tableName == "" {
		return errors.New("table name cannot be empty")
	}
	if len(records) == 0 {
		return errors.New("records slice is empty")
	}

	recordMaps := make([]map[string]interface{}, 0, len(records))
	for i, record := range records {
		if record == nil {
			return fmt.Errorf("record at index %d is nil", i)
		}
		if len(record.Embedding) == 0 {
			return fmt.Errorf("record at index %d has empty embedding", i)
		}

		recordMap := map[string]interface{}{
			"content":   record.Content,
			"embedding": record.Embedding,
		}

		if len(record.Metadata) > 0 {
			recordMap["metadata"] = record.Metadata
		}

		recordMaps = append(recordMaps, recordMap)
	}

	_, _, err := s.Client.From(tableName).Insert(recordMaps, false, "", "", "").Execute()
	if err != nil {
		logger.Error(ctx, "Failed to batch insert vector records",
			zap.Error(err),
			zap.String("table", tableName),
			zap.Int("count", len(records)))
		return fmt.Errorf("failed to batch insert vector records: %w", err)
	}

	logger.Info(ctx, "Batch insert completed successfully",
		zap.String("table", tableName),
		zap.Int("count", len(records)))
	return nil
}

// SimilaritySearch performs a vector similarity search using cosine similarity
// tableName: the table containing vector embeddings
// queryEmbedding: the query vector to search for
// options: search configuration options
// Returns matching records sorted by similarity (highest first)
func (s *VectorDatabaseService) SimilaritySearch(ctx context.Context, tableName string, queryEmbedding []float64, options *VectorSearchOptions) ([]*VectorSearchResult, error) {
	if tableName == "" {
		return nil, errors.New("table name cannot be empty")
	}
	if len(queryEmbedding) == 0 {
		return nil, errors.New("query embedding cannot be empty")
	}
	if options == nil {
		options = &VectorSearchOptions{
			Limit:            10,
			Threshold:        0.0,
			IncludeEmbedding: false,
		}
	}

	// Default limit
	if options.Limit <= 0 {
		options.Limit = 10
	}

	// Build the RPC call for vector similarity search
	// Note: This assumes you have a PostgreSQL function that performs similarity search
	// Example function: match_documents(query_embedding, match_threshold, match_count)
	// params := map[string]interface{}{
	// 	"query_embedding": queryEmbedding,
	// 	"match_threshold": options.Threshold,
	// 	"match_count":     options.Limit,
	// }

	// If using pgvector, you need to call a custom RPC function in Supabase
	//
	// IMPORTANT: You must create the RPC function in Supabase first!
	// Example SQL to create the function:
	//
	// CREATE OR REPLACE FUNCTION match_vectors(
	//   query_embedding vector,
	//   match_threshold float,
	//   match_count int
	// )
	// RETURNS TABLE (
	//   id uuid,
	//   content text,
	//   embedding vector,
	//   metadata jsonb,
	//   similarity float
	// )
	// LANGUAGE plpgsql
	// AS $$
	// BEGIN
	//   RETURN QUERY
	//   SELECT
	//     v.id,
	//     v.content,
	//     v.embedding,
	//     v.metadata,
	//     1 - (v.embedding <=> query_embedding) as similarity
	//   FROM your_table_name v
	//   WHERE 1 - (v.embedding <=> query_embedding) > match_threshold
	//   ORDER BY v.embedding <=> query_embedding
	//   LIMIT match_count;
	// END;
	// $$;

	rpcName := fmt.Sprintf("match_%s", tableName)

	// The Rpc method in supabase-go client works differently than other methods
	// It typically returns the result directly as a string
	// We need to call it and then parse the JSON response
	logger.Info(ctx, "Calling RPC function for vector similarity search",
		zap.String("rpc_name", rpcName),
		zap.String("table", tableName))

	// Note: The current version of supabase-go may have limitations with RPC calls
	// For now, we'll return a helpful error message
	// In a production environment, you would need to:
	// 1. Create the RPC function in Supabase SQL Editor
	// 2. Use the client.Rpc() method appropriately based on your library version
	// 3. Or use raw SQL queries via PostgREST

	logger.Error(ctx, "RPC function not yet implemented",
		zap.String("rpc_name", rpcName),
		zap.String("table", tableName))

	return nil, fmt.Errorf("similarity search via RPC not yet fully implemented - please create RPC function '%s' in Supabase and update this method to call it correctly based on your supabase-go library version", rpcName)

	// TODO: Once you've confirmed the correct Rpc API for your supabase-go version:
	// 1. Call the RPC function
	// 2. Unmarshal the results into []map[string]interface{}
	// 3. Convert to VectorSearchResult structs as shown below:
	//
	// searchResults := make([]*VectorSearchResult, 0, len(results))
	// for _, result := range results {
	//     searchResult := &VectorSearchResult{
	//         ID:         getString(result, "id"),
	//         Content:    getString(result, "content"),
	//         Similarity: getFloat64(result, "similarity"),
	//         Distance:   getFloat64(result, "distance"),
	//     }
	//     if options.IncludeEmbedding {
	//         if embData, ok := result["embedding"].([]interface{}); ok {
	//             embedding := make([]float64, len(embData))
	//             for i, v := range embData {
	//                 if val, ok := v.(float64); ok {
	//                     embedding[i] = val
	//                 }
	//             }
	//             searchResult.Embedding = embedding
	//         }
	//     }
	//     if metadata, ok := result["metadata"].(map[string]interface{}); ok {
	//         searchResult.Metadata = metadata
	//     }
	//     searchResults = append(searchResults, searchResult)
	// }
	// return searchResults, nil
}

// GetByID retrieves a vector record by its ID
func (s *VectorDatabaseService) GetByID(ctx context.Context, tableName string, id string) (*VectorRecord, error) {
	if tableName == "" {
		return nil, errors.New("table name cannot be empty")
	}
	if id == "" {
		return nil, errors.New("id cannot be empty")
	}

	var results []VectorRecord
	count, err := s.Client.From(tableName).
		Select("*", "exact", false).
		Eq("id", id).
		ExecuteTo(&results)

	if err != nil {
		logger.Error(ctx, "Failed to get vector record by ID",
			zap.Error(err),
			zap.String("table", tableName),
			zap.String("id", id))
		return nil, fmt.Errorf("failed to get vector record: %w", err)
	}

	if count == 0 {
		return nil, nil // Not found
	}

	return &results[0], nil
}

// Update updates a vector record by its ID
func (s *VectorDatabaseService) Update(ctx context.Context, tableName string, id string, updates map[string]interface{}) error {
	if tableName == "" {
		return errors.New("table name cannot be empty")
	}
	if id == "" {
		return errors.New("id cannot be empty")
	}
	if len(updates) == 0 {
		return errors.New("updates cannot be empty")
	}

	_, _, err := s.Client.From(tableName).
		Update(updates, "", "").
		Eq("id", id).
		Execute()

	if err != nil {
		logger.Error(ctx, "Failed to update vector record",
			zap.Error(err),
			zap.String("table", tableName),
			zap.String("id", id))
		return fmt.Errorf("failed to update vector record: %w", err)
	}

	logger.Info(ctx, "Vector record updated successfully",
		zap.String("table", tableName),
		zap.String("id", id))
	return nil
}

// Delete deletes a vector record by its ID
func (s *VectorDatabaseService) Delete(ctx context.Context, tableName string, id string) error {
	if tableName == "" {
		return errors.New("table name cannot be empty")
	}
	if id == "" {
		return errors.New("id cannot be empty")
	}

	_, _, err := s.Client.From(tableName).
		Delete("", "").
		Eq("id", id).
		Execute()

	if err != nil {
		logger.Error(ctx, "Failed to delete vector record",
			zap.Error(err),
			zap.String("table", tableName),
			zap.String("id", id))
		return fmt.Errorf("failed to delete vector record: %w", err)
	}

	logger.Info(ctx, "Vector record deleted successfully",
		zap.String("table", tableName),
		zap.String("id", id))
	return nil
}

// BatchDelete deletes multiple vector records by their IDs
// Note: This performs individual delete operations for each ID
// For better performance with large batches, consider using a custom RPC function
func (s *VectorDatabaseService) BatchDelete(ctx context.Context, tableName string, ids []string) error {
	if tableName == "" {
		return errors.New("table name cannot be empty")
	}
	if len(ids) == 0 {
		return errors.New("ids slice is empty")
	}

	// Delete each record individually
	// The supabase-go library may not support In() for arrays properly in all versions
	// So we delete one by one to ensure compatibility
	var deletedCount int
	var lastErr error

	for _, id := range ids {
		_, _, err := s.Client.From(tableName).
			Delete("", "").
			Eq("id", id).
			Execute()

		if err != nil {
			logger.Error(ctx, "Failed to delete vector record in batch",
				zap.Error(err),
				zap.String("table", tableName),
				zap.String("id", id))
			lastErr = err
			// Continue with other deletions
			continue
		}
		deletedCount++
	}

	if lastErr != nil && deletedCount == 0 {
		return fmt.Errorf("failed to delete any vector records: %w", lastErr)
	}

	if lastErr != nil {
		logger.Warn(ctx, "Batch delete completed with some errors",
			zap.String("table", tableName),
			zap.Int("successful", deletedCount),
			zap.Int("total", len(ids)))
		return fmt.Errorf("batch delete completed with errors: deleted %d/%d records", deletedCount, len(ids))
	}

	logger.Info(ctx, "Batch delete completed successfully",
		zap.String("table", tableName),
		zap.Int("count", deletedCount))
	return nil
}

// QueryByMetadata retrieves vector records filtered by metadata
func (s *VectorDatabaseService) QueryByMetadata(ctx context.Context, tableName string, metadataFilters map[string]interface{}, limit int) ([]*VectorRecord, error) {
	if tableName == "" {
		return nil, errors.New("table name cannot be empty")
	}
	if len(metadataFilters) == 0 {
		return nil, errors.New("metadata filters cannot be empty")
	}
	if limit <= 0 {
		limit = 100
	}

	query := s.Client.From(tableName).Select("*", "exact", false)

	// Apply metadata filters
	// Note: This assumes metadata is stored as JSONB in PostgreSQL
	for key, value := range metadataFilters {
		// Using JSONB containment operator
		filterJSON, err := json.Marshal(map[string]interface{}{key: value})
		if err != nil {
			return nil, fmt.Errorf("failed to marshal metadata filter: %w", err)
		}
		query = query.Filter("metadata", "cs", string(filterJSON))
	}

	var results []VectorRecord
	_, err := query.Limit(limit, "").ExecuteTo(&results)

	if err != nil {
		logger.Error(ctx, "Failed to query by metadata",
			zap.Error(err),
			zap.String("table", tableName))
		return nil, fmt.Errorf("failed to query by metadata: %w", err)
	}

	// Convert to pointer slice
	resultPtrs := make([]*VectorRecord, len(results))
	for i := range results {
		resultPtrs[i] = &results[i]
	}

	logger.Info(ctx, "Metadata query completed",
		zap.String("table", tableName),
		zap.Int("results_count", len(results)))

	return resultPtrs, nil
}

// Upsert inserts a new vector record or updates an existing one based on a unique constraint
func (s *VectorDatabaseService) Upsert(ctx context.Context, tableName string, record *VectorRecord, conflictColumn string) error {
	if tableName == "" {
		return errors.New("table name cannot be empty")
	}
	if record == nil {
		return errors.New("record cannot be nil")
	}
	if conflictColumn == "" {
		return errors.New("conflict column cannot be empty")
	}

	recordMap := map[string]interface{}{
		"content":   record.Content,
		"embedding": record.Embedding,
	}

	if len(record.Metadata) > 0 {
		recordMap["metadata"] = record.Metadata
	}

	_, _, err := s.Client.From(tableName).
		Upsert(recordMap, conflictColumn, "*", "exact").
		Execute()

	if err != nil {
		logger.Error(ctx, "Failed to upsert vector record",
			zap.Error(err),
			zap.String("table", tableName))
		return fmt.Errorf("failed to upsert vector record: %w", err)
	}

	logger.Info(ctx, "Vector record upserted successfully", zap.String("table", tableName))
	return nil
}

// Helper functions

func getString(m map[string]interface{}, key string) string {
	if v, ok := m[key].(string); ok {
		return v
	}
	return ""
}

func getFloat64(m map[string]interface{}, key string) float64 {
	if v, ok := m[key].(float64); ok {
		return v
	}
	if v, ok := m[key].(int); ok {
		return float64(v)
	}
	return 0.0
}

func convertToInterfaceSlice(slice []string) []interface{} {
	result := make([]interface{}, len(slice))
	for i, v := range slice {
		result[i] = v
	}
	return result
}
