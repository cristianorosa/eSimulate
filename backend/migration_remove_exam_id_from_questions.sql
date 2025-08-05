-- Script de migração para remover exam_id da tabela questions
-- Execute este script no banco de dados para atualizar a estrutura

-- 1. Verificar se a coluna exam_id existe
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.columns 
               WHERE table_name = 'questions' AND column_name = 'exam_id') THEN
        -- Remover foreign key constraint se existir
        IF EXISTS (SELECT 1 FROM information_schema.table_constraints 
                   WHERE constraint_name = 'fk_questions_exam') THEN
            ALTER TABLE questions DROP CONSTRAINT fk_questions_exam;
        END IF;
        
        -- Remover índice se existir
        DROP INDEX IF EXISTS idx_questions_exam;
        
        -- Remover a coluna exam_id
        ALTER TABLE questions DROP COLUMN exam_id;
        
        RAISE NOTICE 'Coluna exam_id removida da tabela questions';
    ELSE
        RAISE NOTICE 'Coluna exam_id não existe na tabela questions';
    END IF;
END $$;

-- 2. Verificar estrutura final
SELECT column_name, data_type, is_nullable, column_default 
FROM information_schema.columns 
WHERE table_name = 'questions' 
ORDER BY ordinal_position;

-- 3. Verificar foreign keys restantes
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
    AND tc.table_name = 'questions'
ORDER BY kcu.column_name; 