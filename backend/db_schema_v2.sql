-- Script de criação do banco de dados para o eSimulate v2 (PostgreSQL/ANSI SQL)
-- Compatível com PostgreSQL, MySQL, MariaDB, SQL Server, Oracle, etc.
-- Versão: 2.3 - Inclui relacionamento N:N questões-exames e sistema de tags

-- 1. Tabela de Papéis (Roles)
CREATE TABLE roles (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE
);

-- 2. Usuários com níveis de acesso
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(150) NOT NULL UNIQUE,
    password_hash VARCHAR(255),
    role_id INTEGER NOT NULL,
    google_id VARCHAR(100),
    facebook_id VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_users_role FOREIGN KEY (role_id) REFERENCES roles(id)
);

-- 3. Áreas de Conhecimento
CREATE TABLE areas (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 4. Provas (Exams)
CREATE TABLE exams (
    id SERIAL PRIMARY KEY,
    title VARCHAR(200) NOT NULL,
    description TEXT,
    area_id INTEGER NOT NULL,
    max_time_minutes INTEGER NOT NULL,
    passing_score DECIMAL(5,2) DEFAULT 70.0,
    questions_count INTEGER DEFAULT 0,
    is_active BOOLEAN DEFAULT true,
    created_by INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_exams_area FOREIGN KEY (area_id) REFERENCES areas(id),
    CONSTRAINT fk_exams_user FOREIGN KEY (created_by) REFERENCES users(id)
);

-- 5. Tópicos (antigamente "domains") - Agora independentes de exames
CREATE TABLE topics (
    id SERIAL PRIMARY KEY,
    name VARCHAR(150) NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 5.1. Relação N:N entre Exames e Tópicos
CREATE TABLE exam_topics (
    exam_id INTEGER NOT NULL,
    topic_id INTEGER NOT NULL,
    questions_count INTEGER NOT NULL DEFAULT 0,
    weight_percentage DECIMAL(5,2) NOT NULL,
    order_index INTEGER NOT NULL DEFAULT 1,
    difficulty_easy_percentage DECIMAL(5,2) DEFAULT 33.33,
    difficulty_medium_percentage DECIMAL(5,2) DEFAULT 33.33,
    difficulty_hard_percentage DECIMAL(5,2) DEFAULT 33.34,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (exam_id, topic_id),
    CONSTRAINT fk_exam_topics_exam FOREIGN KEY (exam_id) REFERENCES exams(id) ON DELETE CASCADE,
    CONSTRAINT fk_exam_topics_topic FOREIGN KEY (topic_id) REFERENCES topics(id) ON DELETE CASCADE,
    CONSTRAINT exam_topics_weight_check CHECK (weight_percentage >= 0 AND weight_percentage <= 100),
    CONSTRAINT exam_topics_questions_check CHECK (questions_count >= 0),
    CONSTRAINT exam_topics_order_check CHECK (order_index >= 1),
    CONSTRAINT exam_topics_difficulty_check CHECK (
        difficulty_easy_percentage >= 0 AND difficulty_easy_percentage <= 100 AND
        difficulty_medium_percentage >= 0 AND difficulty_medium_percentage <= 100 AND
        difficulty_hard_percentage >= 0 AND difficulty_hard_percentage <= 100 AND
        (difficulty_easy_percentage + difficulty_medium_percentage + difficulty_hard_percentage) = 100
    )
);

-- 6. Questões com suporte a enunciado/problema separados e tipos de questão
CREATE TABLE questions (
    id SERIAL PRIMARY KEY,
    topic_id INTEGER NOT NULL,                  -- Relacionamento direto apenas com tópico
    statement TEXT NOT NULL,                    -- Enunciado da questão
    problem TEXT NOT NULL,                      -- Problema (texto ou código)
    content_type VARCHAR(20) NOT NULL DEFAULT 'text' CHECK (content_type IN ('text', 'code')),
    question_type VARCHAR(20) NOT NULL DEFAULT 'objective' CHECK (question_type IN ('objective', 'multiple_choice')),
    explanation TEXT,                           -- Explicação da questão
    difficulty_level INTEGER DEFAULT 1,         -- 1=Fácil, 2=Médio, 3=Difícil, 4=Muito Difícil, 5=Expert
    created_by INTEGER,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_questions_topic FOREIGN KEY (topic_id) REFERENCES topics(id),
    CONSTRAINT fk_questions_user FOREIGN KEY (created_by) REFERENCES users(id)
);

-- 7. Opções de resposta das questões
CREATE TABLE options (
    id SERIAL PRIMARY KEY,
    question_id INTEGER NOT NULL,
    text TEXT NOT NULL,                         -- Texto da opção
    is_correct BOOLEAN NOT NULL,                -- Se é a resposta correta
    explanation TEXT,                           -- Explicação da opção (opcional)
    order_index INTEGER DEFAULT 0,              -- Ordem de exibição
    CONSTRAINT fk_options_question FOREIGN KEY (question_id) REFERENCES questions(id)
);

-- 7.1. Relação N:N entre Exames e Questões
CREATE TABLE exam_questions (
    exam_id INTEGER NOT NULL,
    question_id INTEGER NOT NULL,
    order_index INTEGER NOT NULL DEFAULT 1,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (exam_id, question_id),
    CONSTRAINT fk_exam_questions_exam FOREIGN KEY (exam_id) REFERENCES exams(id) ON DELETE CASCADE,
    CONSTRAINT fk_exam_questions_question FOREIGN KEY (question_id) REFERENCES questions(id) ON DELETE CASCADE,
    CONSTRAINT exam_questions_order_index_check CHECK (order_index >= 1)
);

-- 7.2. Sistema de Tags para Questões
CREATE TABLE question_tags (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT question_tags_name_check CHECK (LENGTH(TRIM(name)) >= 2)
);

-- 7.3. Relação N:N entre Questões e Tags
CREATE TABLE question_tag_associations (
    question_id INTEGER NOT NULL,
    tag_id INTEGER NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (question_id, tag_id),
    CONSTRAINT fk_question_tag_assoc_question FOREIGN KEY (question_id) REFERENCES questions(id) ON DELETE CASCADE,
    CONSTRAINT fk_question_tag_assoc_tag FOREIGN KEY (tag_id) REFERENCES question_tags(id) ON DELETE CASCADE
);

-- 8. Aplicação de Provas por Usuário
CREATE TABLE user_exams (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    exam_id INTEGER NOT NULL,
    started_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    finished_at TIMESTAMP,
    total_score DECIMAL(5,2),
    passed BOOLEAN,
    time_spent_minutes INTEGER,
    CONSTRAINT fk_user_exams_user FOREIGN KEY (user_id) REFERENCES users(id),
    CONSTRAINT fk_user_exams_exam FOREIGN KEY (exam_id) REFERENCES exams(id)
);

-- 9. Respostas do usuário em cada prova
CREATE TABLE user_answers (
    id SERIAL PRIMARY KEY,
    user_exam_id INTEGER NOT NULL,
    question_id INTEGER NOT NULL,
    option_id INTEGER,
    is_correct BOOLEAN,
    is_marked_for_review BOOLEAN DEFAULT false,
    answered_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_user_answers_user_exam FOREIGN KEY (user_exam_id) REFERENCES user_exams(id),
    CONSTRAINT fk_user_answers_question FOREIGN KEY (question_id) REFERENCES questions(id),
    CONSTRAINT fk_user_answers_option FOREIGN KEY (option_id) REFERENCES options(id)
);

-- 10. Desempenho por Tópico (antigamente "domain_performance")
CREATE TABLE topic_performance (
    id SERIAL PRIMARY KEY,
    user_exam_id INTEGER NOT NULL,
    topic_id INTEGER NOT NULL,
    correct_answers INTEGER,
    total_questions INTEGER,
    score DECIMAL(5,2),
    CONSTRAINT fk_topic_perf_user_exam FOREIGN KEY (user_exam_id) REFERENCES user_exams(id),
    CONSTRAINT fk_topic_perf_topic FOREIGN KEY (topic_id) REFERENCES topics(id)
);

-- 11. Simulados (quizzes)
CREATE TABLE quizzes (
    id SERIAL PRIMARY KEY,
    title VARCHAR(150) NOT NULL,
    description TEXT,
    theme_id INTEGER,
    created_by INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_quizzes_theme FOREIGN KEY (theme_id) REFERENCES topics(id),
    CONSTRAINT fk_quizzes_user FOREIGN KEY (created_by) REFERENCES users(id)
);

-- 12. Relação entre simulados e questões
CREATE TABLE quiz_questions (
    quiz_id INTEGER NOT NULL,
    question_id INTEGER NOT NULL,
    PRIMARY KEY (quiz_id, question_id),
    CONSTRAINT fk_qq_quiz FOREIGN KEY (quiz_id) REFERENCES quizzes(id),
    CONSTRAINT fk_qq_question FOREIGN KEY (question_id) REFERENCES questions(id)
);

-- 13. Simulados realizados por usuário
CREATE TABLE user_quiz (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    quiz_id INTEGER NOT NULL,
    started_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    finished_at TIMESTAMP,
    CONSTRAINT fk_user_quiz_user FOREIGN KEY (user_id) REFERENCES users(id),
    CONSTRAINT fk_user_quiz_quiz FOREIGN KEY (quiz_id) REFERENCES quizzes(id)
);

-- 14. Índices para otimização de performance
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_role ON users(role_id);
CREATE INDEX idx_exams_area ON exams(area_id);
CREATE INDEX idx_exams_active ON exams(is_active);
CREATE INDEX idx_questions_topic ON questions(topic_id);
CREATE INDEX idx_questions_active ON questions(is_active);
CREATE INDEX idx_questions_type ON questions(question_type);
CREATE INDEX idx_questions_content_type ON questions(content_type);
CREATE INDEX idx_options_question ON options(question_id);
CREATE INDEX idx_options_correct ON options(is_correct);
CREATE INDEX idx_user_exams_user ON user_exams(user_id);
CREATE INDEX idx_user_exams_exam ON user_exams(exam_id);
CREATE INDEX idx_user_answers_exam ON user_answers(user_exam_id);
CREATE INDEX idx_topic_performance_exam ON topic_performance(user_exam_id);

-- Índices para as novas tabelas N:N
CREATE INDEX idx_exam_topics_exam ON exam_topics(exam_id);
CREATE INDEX idx_exam_topics_topic ON exam_topics(topic_id);
CREATE INDEX idx_exam_topics_weight ON exam_topics(exam_id, weight_percentage);
CREATE INDEX idx_exam_questions_exam ON exam_questions(exam_id);
CREATE INDEX idx_exam_questions_question ON exam_questions(question_id);
CREATE INDEX idx_exam_questions_order ON exam_questions(exam_id, order_index);
CREATE INDEX idx_question_tags_name ON question_tags(name);
CREATE INDEX idx_question_tag_assoc_question ON question_tag_associations(question_id);
CREATE INDEX idx_question_tag_assoc_tag ON question_tag_associations(tag_id);

-- 15. Dados iniciais

-- Inserir papéis (roles)
INSERT INTO roles (name) VALUES 
('user'),
('redator'),
('admin');

-- Inserir áreas de conhecimento
INSERT INTO areas (name, description) VALUES 
('TI - Tecnologia da Informação', 'Certificações e provas de tecnologia'),
('OAB - Ordem dos Advogados do Brasil', 'Exame da OAB'),
('CRC - Conselho Regional de Contabilidade', 'Exame de suficiência CRC'),
('Concurso Petrobras', 'Concursos públicos da Petrobras'),
('Vestibular', 'Vestibulares de universidades'),
('Concurso Público', 'Outros concursos públicos');

-- Inserir usuário administrador padrão (senha: password)
INSERT INTO users (name, email, password_hash, role_id) VALUES 
('Administrador', 'admin@esimulate.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 3);

-- Inserir redator padrão (senha: password)
INSERT INTO users (name, email, password_hash, role_id) VALUES 
('Redator', 'redator@esimulate.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 2);

-- Inserir usuário padrão (senha: password)
INSERT INTO users (name, email, password_hash, role_id) VALUES 
('Usuário', 'user@esimulate.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 1);

-- 16. Views

-- VALIDATION MOVED TO BUSINESS LOGIC
-- Validation rules are now handled in the application layer (use cases)
-- This provides better error handling and more flexible business rules

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
    q.statement as question_statement,
    q.problem as question_problem,
    q.difficulty_level,
    q.topic_id,
    t.name as topic_name,
    t.description as topic_description,
    et.weight_percentage as topic_weight_in_exam,
    et.questions_count as topic_questions_count_in_exam
FROM exam_questions eq
JOIN exams e ON eq.exam_id = e.id
JOIN questions q ON eq.question_id = q.id
JOIN topics t ON q.topic_id = t.id
LEFT JOIN exam_topics et ON (et.exam_id = eq.exam_id AND et.topic_id = q.topic_id);

-- 17. Dados iniciais para tags
INSERT INTO question_tags (name) VALUES 
    ('Básico'),
    ('Intermediário'),
    ('Avançado'),
    ('Prático'),
    ('Teórico'),
    ('Conceitual'),
    ('Aplicação'),
    ('Análise'),
    ('Síntese'),
    ('Avaliação')
ON CONFLICT (name) DO NOTHING; 