-- Migration: Create question tags tables
-- Description: Creates question_tags and question_tag_associations tables for the tag system
-- Date: 2025-01-22

-- Create question_tags table
CREATE TABLE IF NOT EXISTS question_tags (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create question_tag_associations table for N:N relationship
CREATE TABLE IF NOT EXISTS question_tag_associations (
    question_id INTEGER NOT NULL REFERENCES questions(id) ON DELETE CASCADE,
    tag_id INTEGER NOT NULL REFERENCES question_tags(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (question_id, tag_id)
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_question_tag_associations_question_id ON question_tag_associations(question_id);
CREATE INDEX IF NOT EXISTS idx_question_tag_associations_tag_id ON question_tag_associations(tag_id);
CREATE INDEX IF NOT EXISTS idx_question_tags_name ON question_tags(name);

-- Insert some initial tags
INSERT INTO question_tags (name) VALUES 
    ('Básico'),
    ('Intermediário'),
    ('Avançado'),
    ('Conceitual'),
    ('Prático'),
    ('Teórico'),
    ('Aplicação'),
    ('Análise'),
    ('Síntese'),
    ('Avaliação')
ON CONFLICT (name) DO NOTHING;

COMMIT;
