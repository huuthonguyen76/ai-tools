-- =====================================================
-- Supabase Vector Database Setup Script
-- =====================================================
-- This script sets up the necessary database schema for
-- storing and querying vector embeddings using pgvector
-- =====================================================

-- Enable the pgvector extension
CREATE EXTENSION IF NOT EXISTS vector;

-- =====================================================
-- Create the vectors table
-- =====================================================
-- Adjust the embedding dimension (1536) based on your
-- embedding model:
-- - 384: sentence-transformers/all-MiniLM-L6-v2
-- - 768: BERT base models
-- - 1536: OpenAI text-embedding-ada-002
-- - 3072: OpenAI text-embedding-3-large
-- =====================================================

CREATE TABLE IF NOT EXISTS vectors (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    content TEXT NOT NULL,
    embedding VECTOR(1536),  -- Adjust dimension as needed
    metadata JSONB DEFAULT '{}'::jsonb,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- =====================================================
-- Create indexes for better performance
-- =====================================================

-- Index for faster similarity searches using cosine distance
-- Use ivfflat for datasets < 1M rows
CREATE INDEX IF NOT EXISTS vectors_embedding_idx 
ON vectors 
USING ivfflat (embedding vector_cosine_ops)
WITH (lists = 100);

-- Alternative: Use HNSW for larger datasets (better performance but more memory)
-- Uncomment the following and comment out the ivfflat index above if needed:
/*
CREATE INDEX IF NOT EXISTS vectors_embedding_idx 
ON vectors 
USING hnsw (embedding vector_cosine_ops)
WITH (m = 16, ef_construction = 64);
*/

-- Index for faster metadata queries
CREATE INDEX IF NOT EXISTS vectors_metadata_idx 
ON vectors 
USING gin (metadata);

-- Index for faster timestamp queries
CREATE INDEX IF NOT EXISTS vectors_created_at_idx 
ON vectors (created_at DESC);

-- =====================================================
-- Create updated_at trigger
-- =====================================================
-- Automatically update the updated_at timestamp

CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_vectors_updated_at 
BEFORE UPDATE ON vectors
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

-- =====================================================
-- Create RPC function for similarity search
-- =====================================================
-- This function performs efficient similarity searches
-- using cosine distance (1 - cosine similarity)
-- =====================================================

CREATE OR REPLACE FUNCTION match_vectors(
    query_embedding VECTOR,
    match_threshold FLOAT DEFAULT 0.0,
    match_count INT DEFAULT 10
)
RETURNS TABLE (
    id UUID,
    content TEXT,
    embedding VECTOR,
    metadata JSONB,
    similarity FLOAT,
    distance FLOAT
)
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY
    SELECT 
        v.id,
        v.content,
        v.embedding,
        v.metadata,
        1 - (v.embedding <=> query_embedding) AS similarity,
        v.embedding <=> query_embedding AS distance
    FROM vectors v
    WHERE 1 - (v.embedding <=> query_embedding) > match_threshold
    ORDER BY v.embedding <=> query_embedding
    LIMIT match_count;
END;
$$;

-- =====================================================
-- Create RPC function for similarity search with metadata filtering
-- =====================================================

CREATE OR REPLACE FUNCTION match_vectors_with_metadata(
    query_embedding VECTOR,
    metadata_filter JSONB DEFAULT '{}'::jsonb,
    match_threshold FLOAT DEFAULT 0.0,
    match_count INT DEFAULT 10
)
RETURNS TABLE (
    id UUID,
    content TEXT,
    embedding VECTOR,
    metadata JSONB,
    similarity FLOAT,
    distance FLOAT
)
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY
    SELECT 
        v.id,
        v.content,
        v.embedding,
        v.metadata,
        1 - (v.embedding <=> query_embedding) AS similarity,
        v.embedding <=> query_embedding AS distance
    FROM vectors v
    WHERE 
        (metadata_filter = '{}'::jsonb OR v.metadata @> metadata_filter)
        AND (1 - (v.embedding <=> query_embedding)) > match_threshold
    ORDER BY v.embedding <=> query_embedding
    LIMIT match_count;
END;
$$;

-- =====================================================
-- Create helper function to get vector stats
-- =====================================================

CREATE OR REPLACE FUNCTION get_vector_stats()
RETURNS TABLE (
    total_vectors BIGINT,
    avg_content_length NUMERIC,
    earliest_created TIMESTAMP WITH TIME ZONE,
    latest_created TIMESTAMP WITH TIME ZONE,
    unique_metadata_keys TEXT[]
)
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY
    SELECT 
        COUNT(*)::BIGINT AS total_vectors,
        AVG(LENGTH(content))::NUMERIC AS avg_content_length,
        MIN(created_at) AS earliest_created,
        MAX(created_at) AS latest_created,
        ARRAY_AGG(DISTINCT jsonb_object_keys(metadata)) AS unique_metadata_keys
    FROM vectors;
END;
$$;

-- =====================================================
-- Row Level Security (RLS) Setup - OPTIONAL
-- =====================================================
-- Uncomment the following if you want to enable RLS
-- =====================================================

/*
-- Enable RLS on the vectors table
ALTER TABLE vectors ENABLE ROW LEVEL SECURITY;

-- Policy: Allow authenticated users to read all vectors
CREATE POLICY "Allow authenticated read access" 
ON vectors 
FOR SELECT 
TO authenticated 
USING (true);

-- Policy: Allow authenticated users to insert vectors
CREATE POLICY "Allow authenticated insert access" 
ON vectors 
FOR INSERT 
TO authenticated 
WITH CHECK (true);

-- Policy: Allow authenticated users to update their own vectors
-- Adjust this policy based on your needs
CREATE POLICY "Allow authenticated update access" 
ON vectors 
FOR UPDATE 
TO authenticated 
USING (true)
WITH CHECK (true);

-- Policy: Allow authenticated users to delete vectors
CREATE POLICY "Allow authenticated delete access" 
ON vectors 
FOR DELETE 
TO authenticated 
USING (true);
*/

-- =====================================================
-- Sample Data - OPTIONAL
-- =====================================================
-- Uncomment to insert sample data for testing
-- =====================================================

/*
INSERT INTO vectors (content, embedding, metadata) VALUES
(
    'Artificial Intelligence is transforming the world',
    array_fill(0.1, ARRAY[1536])::vector,  -- Replace with real embedding
    '{"category": "AI", "source": "sample", "tags": ["technology", "ai"]}'::jsonb
),
(
    'Machine Learning enables computers to learn from data',
    array_fill(0.2, ARRAY[1536])::vector,  -- Replace with real embedding
    '{"category": "ML", "source": "sample", "tags": ["technology", "ml"]}'::jsonb
),
(
    'Deep Learning uses neural networks with multiple layers',
    array_fill(0.3, ARRAY[1536])::vector,  -- Replace with real embedding
    '{"category": "DL", "source": "sample", "tags": ["technology", "dl", "neural-networks"]}'::jsonb
);
*/

-- =====================================================
-- Verification Queries
-- =====================================================
-- Run these to verify your setup
-- =====================================================

-- Check if pgvector extension is enabled
SELECT * FROM pg_extension WHERE extname = 'vector';

-- Check table structure
SELECT column_name, data_type 
FROM information_schema.columns 
WHERE table_name = 'vectors';

-- Check indexes
SELECT indexname, indexdef 
FROM pg_indexes 
WHERE tablename = 'vectors';

-- Check RPC functions
SELECT proname, prosrc 
FROM pg_proc 
WHERE proname LIKE 'match_vectors%';

-- Get vector statistics (after inserting data)
-- SELECT * FROM get_vector_stats();

-- =====================================================
-- Maintenance Queries
-- =====================================================

-- Analyze table for query optimization (run after bulk inserts)
-- ANALYZE vectors;

-- Vacuum table to reclaim storage and update statistics
-- VACUUM ANALYZE vectors;

-- Reindex if search performance degrades
-- REINDEX INDEX vectors_embedding_idx;

-- =====================================================
-- Clean Up (Use with caution!)
-- =====================================================
-- Uncomment to drop everything and start fresh
-- =====================================================

/*
DROP TRIGGER IF EXISTS update_vectors_updated_at ON vectors;
DROP FUNCTION IF EXISTS update_updated_at_column();
DROP FUNCTION IF EXISTS match_vectors(VECTOR, FLOAT, INT);
DROP FUNCTION IF EXISTS match_vectors_with_metadata(VECTOR, JSONB, FLOAT, INT);
DROP FUNCTION IF EXISTS get_vector_stats();
DROP TABLE IF EXISTS vectors;
DROP EXTENSION IF EXISTS vector CASCADE;
*/


