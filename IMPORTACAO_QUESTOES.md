# Importação de Questões via JSON

## Visão Geral

Esta funcionalidade permite importar um conjunto completo de questões através de um arquivo JSON. O sistema criará automaticamente as áreas, exames, tópicos e questões necessários.

## Estrutura do JSON

O arquivo JSON deve seguir a seguinte estrutura:

```json
{
  "exam": {
    "title": "Título do Exame",
    "description": "Descrição do exame",
    "max_time": 120,
    "passing_score": 70,
    "is_active": true
  },
  "area": {
    "name": "Nome da Área",
    "description": "Descrição da área"
  },
  "topics": [
    {
      "name": "Nome do Tópico",
      "questions_count": 5,
      "questions": [
        {
          "statement": "Enunciado da questão",
          "problem": "Problema ou código da questão",
          "content_type": "text",
          "question_type": "objetiva",
          "difficulty": "medio",
          "is_active": true,
          "options": [
            {
              "text": "Texto da opção",
              "is_correct": true
            }
          ]
        }
      ]
    }
  ]
}
```

## Campos Obrigatórios

### Exam (Exame)
- `title`: Título do exame
- `max_time`: Tempo máximo em minutos (padrão: 120)
- `passing_score`: Pontuação mínima para aprovação (padrão: 70)
- `is_active`: Se o exame está ativo (padrão: true)

### Area (Área)
- `name`: Nome da área
- `description`: Descrição da área (opcional)

### Topics (Tópicos)
- `name`: Nome do tópico
- `questions_count`: Número de questões no tópico (opcional)
- `questions`: Array de questões

### Questions (Questões)
- `statement`: Enunciado da questão
- `problem`: Problema ou código da questão
- `content_type`: Tipo de conteúdo ("text" ou "code")
- `question_type`: Tipo de questão ("objetiva" ou "multipla_escolha")
- `difficulty`: Nível de dificuldade ("facil", "medio", "dificil")
- `is_active`: Se a questão está ativa (padrão: true)
- `options`: Array de opções

### Options (Opções)
- `text`: Texto da opção
- `is_correct`: Se a opção está correta

## Tipos de Questão

### Objetiva
- Apenas uma opção pode estar correta
- `question_type`: "objetiva"

### Múltipla Escolha
- Uma ou mais opções podem estar corretas
- `question_type`: "multipla_escolha"

## Níveis de Dificuldade

- `facil`: Nível 1
- `medio`: Nível 3
- `dificil`: Nível 5

## Como Usar

1. Acesse a tela de "Gerenciar Questões"
2. Clique no botão "Importar Questões"
3. Selecione o arquivo JSON
4. O sistema processará automaticamente:
   - Criará a área se não existir
   - Criará o exame se não existir
   - Criará os tópicos se não existirem
   - Criará todas as questões e opções

## Validações

O sistema valida:
- Estrutura do JSON
- Campos obrigatórios
- Tipos de questão (objetiva deve ter apenas uma opção correta)
- Pelo menos uma opção correta por questão
- Formato do arquivo (.json)

## Exemplo Completo

Veja o arquivo `question-import-example.json` para um exemplo completo de importação.

## Observações

- Se uma área, exame ou tópico já existir (mesmo nome), o sistema usará o existente
- Questões são sempre criadas como novas
- O sistema recarrega automaticamente os dados após a importação
- Um resumo da importação é exibido ao final do processo 