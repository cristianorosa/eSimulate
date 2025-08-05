-- Script de verificação do schema do banco de dados
-- Execute este script para verificar se todas as tabelas e colunas estão corretas

-- Verificar se a tabela questions tem todas as colunas necessárias
SELECT 
    column_name,
    data_type,
    is_nullable,
    column_default,
    CASE 
        WHEN column_name = 'problem' THEN '✅ Problema separado do enunciado'
        WHEN column_name = 'content_type' THEN '✅ Tipo de conteúdo (text/code)'
        WHEN column_name = 'question_type' THEN '✅ Tipo de questão (objective/multiple_choice)'
        WHEN column_name = 'topic_id' THEN '✅ Referência para tópicos (antigamente domains)'
        ELSE '✅ Coluna padrão'
    END as observacao
FROM information_schema.columns 
WHERE table_name = 'questions' 
ORDER BY ordinal_position;

-- Verificar constraints da tabela questions
SELECT 
    constraint_name,
    constraint_type,
    check_clause
FROM information_schema.check_constraints 
WHERE constraint_name LIKE 'questions_%';

-- Verificar se a tabela options existe e tem a estrutura correta
SELECT 
    column_name,
    data_type,
    is_nullable,
    column_default
FROM information_schema.columns 
WHERE table_name = 'options' 
ORDER BY ordinal_position;

-- Verificar se a tabela topics existe (antigamente domains)
SELECT 
    column_name,
    data_type,
    is_nullable
FROM information_schema.columns 
WHERE table_name = 'topics' 
ORDER BY ordinal_position;

-- Verificar foreign keys
SELECT 
    tc.table_name,
    kcu.column_name,
    ccu.table_name AS foreign_table_name,
    ccu.column_name AS foreign_column_name
FROM information_schema.table_constraints AS tc 
JOIN information_schema.key_column_usage AS kcu
    ON tc.constraint_name = kcu.constraint_name
    AND tc.table_schema = kcu.table_schema
JOIN information_schema.constraint_column_usage AS ccu
    ON ccu.constraint_name = tc.constraint_name
    AND ccu.table_schema = tc.table_schema
WHERE tc.constraint_type = 'FOREIGN KEY' 
    AND tc.table_name IN ('questions', 'options', 'topics')
ORDER BY tc.table_name, kcu.column_name;

-- Verificar índices importantes
SELECT 
    indexname,
    tablename,
    indexdef
FROM pg_indexes 
WHERE tablename IN ('questions', 'options', 'topics')
ORDER BY tablename, indexname;

-- Resumo das tabelas principais
SELECT 
    schemaname,
    tablename,
    tableowner
FROM pg_tables 
WHERE tablename IN ('questions', 'options', 'topics', 'exams', 'areas', 'users', 'roles')
ORDER BY tablename; 