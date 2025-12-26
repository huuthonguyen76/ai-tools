-- =====================================================
-- Facebook Posts Table Setup Script
-- =====================================================
-- This script creates a table to store Facebook posts
-- scraped from Apify
-- =====================================================

-- Create the facebook_posts table
CREATE TABLE IF NOT EXISTS facebook_posts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    dataset_id TEXT NOT NULL,
    post_id TEXT,
    content TEXT,
    url TEXT,
    author_name TEXT,
    author_id TEXT,
    likes INTEGER DEFAULT 0,
    comments_count INTEGER DEFAULT 0,
    shares INTEGER DEFAULT 0,
    post_timestamp TEXT,
    group_name TEXT,
    group_id TEXT,
    raw_data JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- =====================================================
-- Create indexes for better query performance
-- =====================================================

-- Index on dataset_id for filtering by source dataset
CREATE INDEX IF NOT EXISTS facebook_posts_dataset_id_idx 
ON facebook_posts (dataset_id);

-- Index on post_id for unique post lookups
CREATE INDEX IF NOT EXISTS facebook_posts_post_id_idx 
ON facebook_posts (post_id);

-- Index on author_id for author-based queries
CREATE INDEX IF NOT EXISTS facebook_posts_author_id_idx 
ON facebook_posts (author_id);

-- Index on group_id for group-based queries
CREATE INDEX IF NOT EXISTS facebook_posts_group_id_idx 
ON facebook_posts (group_id);

-- Index on created_at for time-based queries
CREATE INDEX IF NOT EXISTS facebook_posts_created_at_idx 
ON facebook_posts (created_at DESC);

-- Full-text search index on content
CREATE INDEX IF NOT EXISTS facebook_posts_content_gin_idx 
ON facebook_posts 
USING gin (to_tsvector('english', COALESCE(content, '')));

-- GIN index on raw_data for JSONB queries
CREATE INDEX IF NOT EXISTS facebook_posts_raw_data_idx 
ON facebook_posts 
USING gin (raw_data);

-- =====================================================
-- Create updated_at trigger
-- =====================================================

CREATE OR REPLACE FUNCTION update_facebook_posts_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_facebook_posts_updated_at_trigger
BEFORE UPDATE ON facebook_posts
FOR EACH ROW
EXECUTE FUNCTION update_facebook_posts_updated_at();

-- =====================================================
-- Create unique constraint to prevent duplicate posts
-- =====================================================
-- This prevents inserting the same post twice from the same dataset

ALTER TABLE facebook_posts 
ADD CONSTRAINT facebook_posts_unique_post 
UNIQUE NULLS NOT DISTINCT (dataset_id, post_id);

-- =====================================================
-- Helpful queries
-- =====================================================

-- Get posts by dataset
-- SELECT * FROM facebook_posts WHERE dataset_id = 'your_dataset_id' ORDER BY created_at DESC;

-- Get posts by group
-- SELECT * FROM facebook_posts WHERE group_id = 'your_group_id' ORDER BY created_at DESC;

-- Search posts by content
-- SELECT * FROM facebook_posts 
-- WHERE to_tsvector('english', content) @@ to_tsquery('english', 'search_term');

-- Get engagement stats by group
-- SELECT 
--     group_name,
--     COUNT(*) as total_posts,
--     SUM(likes) as total_likes,
--     SUM(comments_count) as total_comments,
--     SUM(shares) as total_shares
-- FROM facebook_posts
-- GROUP BY group_name;

-- =====================================================
-- Clean Up (Use with caution!)
-- =====================================================
-- Uncomment to drop everything and start fresh
-- =====================================================

/*
DROP TRIGGER IF EXISTS update_facebook_posts_updated_at_trigger ON facebook_posts;
DROP FUNCTION IF EXISTS update_facebook_posts_updated_at();
DROP TABLE IF EXISTS facebook_posts;
*/









