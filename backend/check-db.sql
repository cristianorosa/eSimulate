-- Script para verificar o estado do banco de dados

-- Verificar se as tabelas principais existem
SELECT 'users' as table_name, 
       EXISTS(SELECT FROM information_schema.tables WHERE table_schema = 'public' AND table_name = 'users') as exists
UNION ALL
SELECT 'user_quiz', 
       EXISTS(SELECT FROM information_schema.tables WHERE table_schema = 'public' AND table_name = 'user_quiz')
UNION ALL
SELECT 'user_answers', 
       EXISTS(SELECT FROM information_schema.tables WHERE table_schema = 'public' AND table_name = 'user_answers')
UNION ALL
SELECT 'options', 
       EXISTS(SELECT FROM information_schema.tables WHERE table_schema = 'public' AND table_name = 'options')
UNION ALL
SELECT 'roles', 
       EXISTS(SELECT FROM information_schema.tables WHERE table_schema = 'public' AND table_name = 'roles')
UNION ALL
SELECT 'areas', 
       EXISTS(SELECT FROM information_schema.tables WHERE table_schema = 'public' AND table_name = 'areas')
UNION ALL
SELECT 'exams', 
       EXISTS(SELECT FROM information_schema.tables WHERE table_schema = 'public' AND table_name = 'exams')
UNION ALL
SELECT 'domains', 
       EXISTS(SELECT FROM information_schema.tables WHERE table_schema = 'public' AND table_name = 'domains')
UNION ALL
SELECT 'questions', 
       EXISTS(SELECT FROM information_schema.tables WHERE table_schema = 'public' AND table_name = 'questions')
UNION ALL
SELECT 'user_exams', 
       EXISTS(SELECT FROM information_schema.tables WHERE table_schema = 'public' AND table_name = 'user_exams')
ORDER BY table_name;

-- Verificar dados de usuários
SELECT 'users_count' as info, COUNT(*) as value FROM users
UNION ALL
SELECT 'users_with_roles', COUNT(*) FROM users WHERE role_id IS NOT NULL;

-- Verificar se há dados de performance
SELECT 'user_quiz_count' as info, COUNT(*) as value FROM user_quiz
UNION ALL
SELECT 'user_answers_count', COUNT(*) FROM user_answers; 