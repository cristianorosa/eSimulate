-- Migration: Exam-Questions N:N Relationship with Tags System
-- Version: 2.3
-- Description: Creates exam_questions table, question_tags system, and validation functions
-- Based on requirements: Questions can belong to multiple exams, with hierarchical validation

-- ============================================================================
-- EXAM-QUESTIONS N:N RELATIONSHIP
-- ============================================================================

-- Table: exam_questions (N:N relationship between exams and questions)
CREATE TABLE IF NOT EXISTS exam_questions (
    exam_id INTEGER NOT NULL REFERENCES exams(id) ON DELETE CASCADE,
    question_id INTEGER NOT NULL REFERENCES questions(id) ON DELETE CASCADE,
    order_index INTEGER NOT NULL DEFAULT 1,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    PRIMARY KEY (exam_id, question_id),
    CONSTRAINT exam_questions_order_index_check CHECK (order_index >= 1)
);

-- ============================================================================
-- QUESTION TAGS SYSTEM
-- ============================================================================

-- Table: question_tags (tag definitions)
CREATE TABLE IF NOT EXISTS question_tags (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT question_tags_name_check CHECK (LENGTH(TRIM(name)) >= 2)
);

-- Table: question_tag_associations (N:N between questions and tags)
CREATE TABLE IF NOT EXISTS question_tag_associations (
    question_id INTEGER NOT NULL REFERENCES questions(id) ON DELETE CASCADE,
    tag_id INTEGER NOT NULL REFERENCES question_tags(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    PRIMARY KEY (question_id, tag_id)
);

-- ============================================================================
-- VALIDATION FUNCTION
-- ============================================================================

-- Function: validate_exam_question_topic
-- Purpose: Ensures a question associated with an exam belongs to one of the exam's topics
-- Exception: If exam has only one topic with 100% weight, any question can be associated
CREATE OR REPLACE FUNCTION validate_exam_question_topic()
RETURNS TRIGGER AS $$
DECLARE
    question_topic_id INTEGER;
    exam_topic_count INTEGER;
    exam_single_topic_weight DECIMAL;
    exam_has_topic BOOLEAN;
BEGIN
    -- Get the topic_id of the question
    SELECT topic_id INTO question_topic_id
    FROM questions
    WHERE id = NEW.question_id;
    
    -- Check if question has a topic
    IF question_topic_id IS NULL THEN
        RAISE EXCEPTION 'Question % does not have a topic assigned', NEW.question_id;
    END IF;
    
    -- Count how many topics the exam has
    SELECT COUNT(*) INTO exam_topic_count
    FROM exam_topics
    WHERE exam_id = NEW.exam_id;
    
    -- If exam has no topics, reject the association
    IF exam_topic_count = 0 THEN
        RAISE EXCEPTION 'Exam % has no topics configured', NEW.exam_id;
    END IF;
    
    -- Check if exam has only one topic with 100% weight (exception case)
    IF exam_topic_count = 1 THEN
        SELECT weight_percentage INTO exam_single_topic_weight
        FROM exam_topics
        WHERE exam_id = NEW.exam_id;
        
        -- If single topic has 100% weight, allow any question
        IF exam_single_topic_weight = 100.0 THEN
            RETURN NEW;
        END IF;
    END IF;
    
    -- Check if the question's topic is associated with the exam
    SELECT EXISTS(
        SELECT 1
        FROM exam_topics et
        WHERE et.exam_id = NEW.exam_id
        AND et.topic_id = question_topic_id
    ) INTO exam_has_topic;
    
    -- If topic is not associated with exam, reject
    IF NOT exam_has_topic THEN
        RAISE EXCEPTION 'Question % belongs to topic % which is not associated with exam %', 
            NEW.question_id, question_topic_id, NEW.exam_id;
    END IF;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- ============================================================================
-- TRIGGERS
-- ============================================================================

-- Trigger: validate exam-question-topic relationship
DROP TRIGGER IF EXISTS trigger_validate_exam_question_topic ON exam_questions;
CREATE TRIGGER trigger_validate_exam_question_topic
    BEFORE INSERT OR UPDATE ON exam_questions
    FOR EACH ROW
    EXECUTE FUNCTION validate_exam_question_topic();

-- ============================================================================
-- VIEWS
-- ============================================================================

-- View: v_exam_questions_with_topics
-- Purpose: Provides easy access to exam-question relationships with topic information
CREATE OR REPLACE VIEW v_exam_questions_with_topics AS
SELECT 
    eq.exam_id,
    eq.question_id,
    eq.order_index,
    eq.created_at as association_created_at,
    e.title as exam_title,
    e.description as exam_description,
    q.question_text,
    q.difficulty_level,
    q.topic_id,
    t.title as topic_title,
    t.description as topic_description,
    et.weight_percentage as topic_weight_in_exam,
    et.questions_count as topic_questions_count_in_exam
FROM exam_questions eq
JOIN exams e ON eq.exam_id = e.id
JOIN questions q ON eq.question_id = q.id
JOIN topics t ON q.topic_id = t.id
LEFT JOIN exam_topics et ON (et.exam_id = eq.exam_id AND et.topic_id = q.topic_id);

-- ============================================================================
-- INDEXES
-- ============================================================================

-- Indexes for exam_questions
CREATE INDEX IF NOT EXISTS idx_exam_questions_exam_id ON exam_questions(exam_id);
CREATE INDEX IF NOT EXISTS idx_exam_questions_question_id ON exam_questions(question_id);
CREATE INDEX IF NOT EXISTS idx_exam_questions_order_index ON exam_questions(exam_id, order_index);

-- Indexes for question_tags
CREATE INDEX IF NOT EXISTS idx_question_tags_name ON question_tags(name);

-- Indexes for question_tag_associations
CREATE INDEX IF NOT EXISTS idx_question_tag_associations_question_id ON question_tag_associations(question_id);
CREATE INDEX IF NOT EXISTS idx_question_tag_associations_tag_id ON question_tag_associations(tag_id);

-- ============================================================================
-- SAMPLE DATA
-- ============================================================================

-- Insert sample question tags
INSERT INTO question_tags (name, description) VALUES 
    ('Básico', 'Questões de nível básico/fundamental'),
    ('Intermediário', 'Questões de nível intermediário'),
    ('Avançado', 'Questões de nível avançado'),
    ('Prático', 'Questões com enfoque prático'),
    ('Teórico', 'Questões com enfoque teórico'),
    ('Conceitual', 'Questões conceituais'),
    ('Aplicação', 'Questões de aplicação prática'),
    ('Análise', 'Questões que requerem análise'),
    ('Síntese', 'Questões que requerem síntese'),
    ('Avaliação', 'Questões de avaliação crítica')
ON CONFLICT (name) DO NOTHING;

-- ============================================================================
-- COMMENTS
-- ============================================================================

-- Table comments
COMMENT ON TABLE exam_questions IS 'N:N relationship between exams and questions with order control';
COMMENT ON TABLE question_tags IS 'Tag definitions for categorizing questions';
COMMENT ON TABLE question_tag_associations IS 'N:N relationship between questions and tags';

-- Column comments
COMMENT ON COLUMN exam_questions.order_index IS 'Order of question within the exam (1-based)';
COMMENT ON COLUMN question_tags.name IS 'Unique tag name for categorization';
COMMENT ON COLUMN question_tags.description IS 'Optional description of the tag purpose';

-- Function comments
COMMENT ON FUNCTION validate_exam_question_topic() IS 'Validates that questions belong to exam topics, with exception for single 100% weight topic exams';

-- View comments
COMMENT ON VIEW v_exam_questions_with_topics IS 'Consolidated view of exam-question relationships with topic details';
