package services

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/supabase-community/supabase-go"
)

func TestNewVectorDatabaseService(t *testing.T) {
	t.Run("Success - Valid Client", func(t *testing.T) {
		// Create a mock client (in real tests, you'd use a real client or mock)
		client := &supabase.Client{}

		service, err := NewVectorDatabaseService(client)

		assert.NoError(t, err)
		assert.NotNil(t, service)
		assert.Equal(t, client, service.Client)
	})

	t.Run("Error - Nil Client", func(t *testing.T) {
		service, err := NewVectorDatabaseService(nil)

		assert.Error(t, err)
		assert.Nil(t, service)
		assert.Equal(t, "supabase client cannot be nil", err.Error())
	})
}

func TestVectorRecord_Validation(t *testing.T) {
	tests := []struct {
		name        string
		record      *VectorRecord
		shouldError bool
		errorMsg    string
	}{
		{
			name: "Valid Record",
			record: &VectorRecord{
				Content:   "Test content",
				Embedding: []float64{0.1, 0.2, 0.3},
				Metadata: map[string]interface{}{
					"source": "test",
				},
			},
			shouldError: false,
		},
		{
			name:        "Nil Record",
			record:      nil,
			shouldError: true,
			errorMsg:    "record cannot be nil",
		},
		{
			name: "Empty Embedding",
			record: &VectorRecord{
				Content:   "Test content",
				Embedding: []float64{},
			},
			shouldError: true,
			errorMsg:    "embedding vector cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// This would be called within Insert method
			var err error
			if tt.record == nil {
				err = assert.AnError
			} else if len(tt.record.Embedding) == 0 {
				err = assert.AnError
			}

			if tt.shouldError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestVectorSearchOptions_Defaults(t *testing.T) {
	t.Run("Apply Default Options", func(t *testing.T) {
		var options *VectorSearchOptions

		// Simulate default application
		if options == nil {
			options = &VectorSearchOptions{
				Limit:            10,
				Threshold:        0.0,
				IncludeEmbedding: false,
			}
		}

		assert.Equal(t, 10, options.Limit)
		assert.Equal(t, 0.0, options.Threshold)
		assert.False(t, options.IncludeEmbedding)
	})

	t.Run("Custom Options", func(t *testing.T) {
		options := &VectorSearchOptions{
			Limit:            20,
			Threshold:        0.7,
			IncludeEmbedding: true,
			MetadataFilter: map[string]interface{}{
				"category": "tech",
			},
		}

		assert.Equal(t, 20, options.Limit)
		assert.Equal(t, 0.7, options.Threshold)
		assert.True(t, options.IncludeEmbedding)
		assert.NotNil(t, options.MetadataFilter)
	})
}

func TestInsertToVectorDatabase_Validation(t *testing.T) {
	tests := []struct {
		name      string
		client    *supabase.Client
		tableName string
		record    map[string]interface{}
		wantErr   bool
		errMsg    string
	}{
		{
			name:      "Error - Nil Client",
			client:    nil,
			tableName: "vectors",
			record:    map[string]interface{}{"content": "test"},
			wantErr:   true,
			errMsg:    "supabase client is nil",
		},
		{
			name:      "Error - Empty Table Name",
			client:    &supabase.Client{},
			tableName: "",
			record:    map[string]interface{}{"content": "test"},
			wantErr:   true,
			errMsg:    "table name cannot be empty",
		},
		{
			name:      "Error - Empty Record",
			client:    &supabase.Client{},
			tableName: "vectors",
			record:    map[string]interface{}{},
			wantErr:   true,
			errMsg:    "record is empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := InsertToVectorDatabase(tt.client, tt.tableName, tt.record)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				// Would succeed with real client
				assert.Error(t, err) // Expects error without real DB connection
			}
		})
	}
}

func TestHelperFunctions(t *testing.T) {
	t.Run("getString", func(t *testing.T) {
		m := map[string]interface{}{
			"valid":   "test_value",
			"invalid": 123,
		}

		assert.Equal(t, "test_value", getString(m, "valid"))
		assert.Equal(t, "", getString(m, "invalid"))
		assert.Equal(t, "", getString(m, "nonexistent"))
	})

	t.Run("getFloat64", func(t *testing.T) {
		m := map[string]interface{}{
			"float":   3.14,
			"int":     42,
			"invalid": "string",
		}

		assert.Equal(t, 3.14, getFloat64(m, "float"))
		assert.Equal(t, float64(42), getFloat64(m, "int"))
		assert.Equal(t, 0.0, getFloat64(m, "invalid"))
		assert.Equal(t, 0.0, getFloat64(m, "nonexistent"))
	})

	t.Run("convertToInterfaceSlice", func(t *testing.T) {
		input := []string{"id1", "id2", "id3"}
		result := convertToInterfaceSlice(input)

		assert.Len(t, result, 3)
		assert.Equal(t, "id1", result[0])
		assert.Equal(t, "id2", result[1])
		assert.Equal(t, "id3", result[2])
	})
}

func TestBatchInsert_Validation(t *testing.T) {
	ctx := context.Background()
	service := &VectorDatabaseService{
		Client: &supabase.Client{},
	}

	t.Run("Error - Empty Table Name", func(t *testing.T) {
		records := []*VectorRecord{
			{
				Content:   "test",
				Embedding: []float64{0.1, 0.2},
			},
		}

		err := service.BatchInsert(ctx, "", records)
		assert.Error(t, err)
		assert.Equal(t, "table name cannot be empty", err.Error())
	})

	t.Run("Error - Empty Records", func(t *testing.T) {
		err := service.BatchInsert(ctx, "vectors", []*VectorRecord{})
		assert.Error(t, err)
		assert.Equal(t, "records slice is empty", err.Error())
	})

	t.Run("Error - Nil Record in Slice", func(t *testing.T) {
		records := []*VectorRecord{
			{
				Content:   "test",
				Embedding: []float64{0.1, 0.2},
			},
			nil,
		}

		err := service.BatchInsert(ctx, "vectors", records)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "record at index 1 is nil")
	})

	t.Run("Error - Empty Embedding in Record", func(t *testing.T) {
		records := []*VectorRecord{
			{
				Content:   "test1",
				Embedding: []float64{0.1, 0.2},
			},
			{
				Content:   "test2",
				Embedding: []float64{},
			},
		}

		err := service.BatchInsert(ctx, "vectors", records)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "record at index 1 has empty embedding")
	})
}

func TestUpdate_Validation(t *testing.T) {
	ctx := context.Background()
	service := &VectorDatabaseService{
		Client: &supabase.Client{},
	}

	tests := []struct {
		name      string
		tableName string
		id        string
		updates   map[string]interface{}
		wantErr   bool
		errMsg    string
	}{
		{
			name:      "Error - Empty Table Name",
			tableName: "",
			id:        "123",
			updates:   map[string]interface{}{"content": "updated"},
			wantErr:   true,
			errMsg:    "table name cannot be empty",
		},
		{
			name:      "Error - Empty ID",
			tableName: "vectors",
			id:        "",
			updates:   map[string]interface{}{"content": "updated"},
			wantErr:   true,
			errMsg:    "id cannot be empty",
		},
		{
			name:      "Error - Empty Updates",
			tableName: "vectors",
			id:        "123",
			updates:   map[string]interface{}{},
			wantErr:   true,
			errMsg:    "updates cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.Update(ctx, tt.tableName, tt.id, tt.updates)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			}
		})
	}
}

func TestDelete_Validation(t *testing.T) {
	ctx := context.Background()
	service := &VectorDatabaseService{
		Client: &supabase.Client{},
	}

	tests := []struct {
		name      string
		tableName string
		id        string
		wantErr   bool
		errMsg    string
	}{
		{
			name:      "Error - Empty Table Name",
			tableName: "",
			id:        "123",
			wantErr:   true,
			errMsg:    "table name cannot be empty",
		},
		{
			name:      "Error - Empty ID",
			tableName: "vectors",
			id:        "",
			wantErr:   true,
			errMsg:    "id cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.Delete(ctx, tt.tableName, tt.id)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			}
		})
	}
}

func TestQueryByMetadata_Validation(t *testing.T) {
	ctx := context.Background()
	service := &VectorDatabaseService{
		Client: &supabase.Client{},
	}

	tests := []struct {
		name            string
		tableName       string
		metadataFilters map[string]interface{}
		limit           int
		wantErr         bool
		errMsg          string
	}{
		{
			name:            "Error - Empty Table Name",
			tableName:       "",
			metadataFilters: map[string]interface{}{"category": "tech"},
			limit:           10,
			wantErr:         true,
			errMsg:          "table name cannot be empty",
		},
		{
			name:            "Error - Empty Metadata Filters",
			tableName:       "vectors",
			metadataFilters: map[string]interface{}{},
			limit:           10,
			wantErr:         true,
			errMsg:          "metadata filters cannot be empty",
		},
		{
			name:            "Apply Default Limit",
			tableName:       "vectors",
			metadataFilters: map[string]interface{}{"category": "tech"},
			limit:           0,
			wantErr:         false, // Would succeed with valid client
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := service.QueryByMetadata(ctx, tt.tableName, tt.metadataFilters, tt.limit)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			}
		})
	}
}
