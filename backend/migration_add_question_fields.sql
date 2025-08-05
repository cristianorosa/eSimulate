-- Script de migração para adicionar campos na tabela questions
-- Execute este script no banco de dados para atualizar a estrutura

-- 1. Adicionar coluna problem (se não existir)
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns 
                   WHERE table_name = 'questions' AND column_name = 'problem') THEN
        ALTER TABLE questions ADD COLUMN problem TEXT;
    END IF;
END $$;

-- 2. Adicionar coluna content_type (se não existir)
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns 
                   WHERE table_name = 'questions' AND column_name = 'content_type') THEN
        ALTER TABLE questions ADD COLUMN content_type VARCHAR(20) DEFAULT 'text';
    END IF;
END $$;

-- 3. Adicionar coluna question_type (se não existir)
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns 
                   WHERE table_name = 'questions' AND column_name = 'question_type') THEN
        ALTER TABLE questions ADD COLUMN question_type VARCHAR(20) DEFAULT 'objective';
    END IF;
END $$;

-- 4. Adicionar constraints para content_type
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM information_schema.check_constraints 
                   WHERE constraint_name = 'questions_content_type_check') THEN
        ALTER TABLE questions ADD CONSTRAINT questions_content_type_check 
        CHECK (content_type IN ('text', 'code'));
    END IF;
END $$;

-- 5. Adicionar constraints para question_type
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM information_schema.check_constraints 
                   WHERE constraint_name = 'questions_question_type_check') THEN
        ALTER TABLE questions ADD CONSTRAINT questions_question_type_check 
        CHECK (question_type IN ('objective', 'multiple_choice'));
    END IF;
END $$;

-- 6. Renomear coluna domain_id para topic_id (se existir)
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.columns 
               WHERE table_name = 'questions' AND column_name = 'domain_id') THEN
        ALTER TABLE questions RENAME COLUMN domain_id TO topic_id;
    END IF;
END $$;

-- 7. Verificar se a tabela topics existe, se não, criar
CREATE TABLE IF NOT EXISTS topics (
    id SERIAL PRIMARY KEY,
    exam_id INTEGER NOT NULL,
    name VARCHAR(150) NOT NULL,
    weight_percentage DECIMAL(5,2) NOT NULL,
    order_index INTEGER DEFAULT 0,
    questions_count INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 8. Adicionar foreign key para topics se não existir
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM information_schema.table_constraints 
                   WHERE constraint_name = 'fk_questions_topic') THEN
        ALTER TABLE questions ADD CONSTRAINT fk_questions_topic 
        FOREIGN KEY (topic_id) REFERENCES topics(id);
    END IF;
END $$;

-- 9. Verificar se a tabela options existe, se não, criar
CREATE TABLE IF NOT EXISTS options (
    id SERIAL PRIMARY KEY,
    question_id INTEGER NOT NULL,
    text TEXT NOT NULL,
    is_correct BOOLEAN NOT NULL,
    explanation TEXT,
    order_index INTEGER DEFAULT 0
);

-- 10. Adicionar foreign key para options se não existir
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM information_schema.table_constraints 
                   WHERE constraint_name = 'fk_options_question') THEN
        ALTER TABLE options ADD CONSTRAINT fk_options_question 
        FOREIGN KEY (question_id) REFERENCES questions(id);
    END IF;
END $$;

-- 11. Atualizar dados existentes (se necessário)
UPDATE questions SET problem = statement WHERE problem IS NULL;
UPDATE questions SET content_type = 'text' WHERE content_type IS NULL;
UPDATE questions SET question_type = 'objective' WHERE question_type IS NULL;

-- 12. Verificar estrutura final
SELECT column_name, data_type, is_nullable, column_default 
FROM information_schema.columns 
WHERE table_name = 'questions' 
ORDER BY ordinal_position; 