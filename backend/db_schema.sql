-- Script de criação do banco de dados para o eSimulate (PostgreSQL)

-- 1. Usuários
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(150) UNIQUE NOT NULL,
    password_hash VARCHAR(255), -- pode ser NULL para login social
    google_id VARCHAR(100),
    facebook_id VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 2. Temas e Subtemas
CREATE TABLE themes (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    parent_id INTEGER REFERENCES themes(id) ON DELETE CASCADE -- permite subtemas
);

-- 3. Questões
CREATE TABLE questions (
    id SERIAL PRIMARY KEY,
    theme_id INTEGER REFERENCES themes(id) ON DELETE SET NULL,
    statement TEXT NOT NULL, -- enunciado da questão
    explanation TEXT, -- explicação geral da questão
    created_by INTEGER REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 4. Opções de resposta
CREATE TABLE options (
    id SERIAL PRIMARY KEY,
    question_id INTEGER REFERENCES questions(id) ON DELETE CASCADE,
    text TEXT NOT NULL,
    is_correct BOOLEAN NOT NULL,
    explanation TEXT -- explicação do porquê está certa ou errada
);

-- 5. Simulados (quizzes)
CREATE TABLE quizzes (
    id SERIAL PRIMARY KEY,
    title VARCHAR(150) NOT NULL,
    description TEXT,
    theme_id INTEGER REFERENCES themes(id) ON DELETE SET NULL,
    created_by INTEGER REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 6. Relação entre simulados e questões
CREATE TABLE quiz_questions (
    quiz_id INTEGER REFERENCES quizzes(id) ON DELETE CASCADE,
    question_id INTEGER REFERENCES questions(id) ON DELETE CASCADE,
    PRIMARY KEY (quiz_id, question_id)
);

-- 7. Simulados realizados por usuário
CREATE TABLE user_quiz (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    quiz_id INTEGER REFERENCES quizzes(id) ON DELETE CASCADE,
    started_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    finished_at TIMESTAMP
);

-- 8. Respostas do usuário em cada simulado
CREATE TABLE user_answers (
    id SERIAL PRIMARY KEY,
    user_quiz_id INTEGER REFERENCES user_quiz(id) ON DELETE CASCADE,
    question_id INTEGER REFERENCES questions(id) ON DELETE CASCADE,
    option_id INTEGER REFERENCES options(id) ON DELETE SET NULL,
    answered_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Índices e constraints adicionais podem ser criados conforme a necessidade de performance e integridade.
