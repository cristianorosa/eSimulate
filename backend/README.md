# eSimulate Backend

Sistema de simulados educacionais com relacionamentos N:N entre exames, tópicos e questões.

## 📋 Índice

- [Arquitetura](#arquitetura)
- [Modelo de Dados](#modelo-de-dados)
- [Relacionamentos](#relacionamentos)
- [API Endpoints](#api-endpoints)
- [Regras de Negócio](#regras-de-negócio)
- [Instalação](#instalação)
- [Configuração](#configuração)

## 🏗️ Arquitetura

O sistema segue os princípios de **Clean Architecture**:

```
interfaces/ (HTTP Handlers)
    ↓
usecase/ (Business Logic)
    ↓
domain/ (Entities & Interfaces)
    ↓
infra/ (Database Implementation)
```

### Estrutura de Diretórios

```
backend/
├── domain/           # Entidades e interfaces
├── usecase/          # Regras de negócio
├── interfaces/       # Handlers HTTP
├── infra/           # Implementação PostgreSQL
├── db_schema_v2.sql # Schema completo do banco
└── main.go          # Injeção de dependências
```

## 📊 Modelo de Dados

### Entidades Principais

#### **Áreas**
```sql
CREATE TABLE areas (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
```

#### **Exames**
```sql
CREATE TABLE exams (
    id SERIAL PRIMARY KEY,
    area_id INTEGER NOT NULL REFERENCES areas(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    max_time_minutes INTEGER NOT NULL CHECK (max_time_minutes > 0),
    passing_score DECIMAL(5,2) NOT NULL CHECK (passing_score >= 0 AND passing_score <= 100),
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
```

#### **Tópicos**
```sql
CREATE TABLE topics (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
```

#### **Questões**
```sql
CREATE TABLE questions (
    id SERIAL PRIMARY KEY,
    topic_id INTEGER NOT NULL REFERENCES topics(id) ON DELETE CASCADE,
    text TEXT NOT NULL,
    difficulty_level difficulty_level NOT NULL,
    explanation TEXT,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
```

### Relacionamentos N:N

#### **Exames ↔ Tópicos**
```sql
CREATE TABLE exam_topics (
    exam_id INTEGER NOT NULL REFERENCES exams(id) ON DELETE CASCADE,
    topic_id INTEGER NOT NULL REFERENCES topics(id) ON DELETE CASCADE,
    questions_count INTEGER NOT NULL CHECK (questions_count > 0),
    weight_percentage DECIMAL(5,2) NOT NULL CHECK (weight_percentage > 0 AND weight_percentage <= 100),
    order_index INTEGER NOT NULL DEFAULT 0,
    difficulty_easy_percentage DECIMAL(5,2) DEFAULT 30.0,
    difficulty_medium_percentage DECIMAL(5,2) DEFAULT 50.0,
    difficulty_hard_percentage DECIMAL(5,2) DEFAULT 20.0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (exam_id, topic_id)
);
```

#### **Exames ↔ Questões**
```sql
CREATE TABLE exam_questions (
    exam_id INTEGER NOT NULL REFERENCES exams(id) ON DELETE CASCADE,
    question_id INTEGER NOT NULL REFERENCES questions(id) ON DELETE CASCADE,
    order_index INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (exam_id, question_id)
);
```

#### **Sistema de Tags**
```sql
-- Tags simples (apenas nome)
CREATE TABLE question_tags (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Associação N:N Questões ↔ Tags
CREATE TABLE question_tag_associations (
    question_id INTEGER NOT NULL REFERENCES questions(id) ON DELETE CASCADE,
    tag_id INTEGER NOT NULL REFERENCES question_tags(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (question_id, tag_id)
);
```

## 🔗 Relacionamentos

### Hierarquia do Sistema

```
Área
└── Exame
    ├── Tópicos (N:N com peso e distribuição de dificuldade)
    └── Questões (N:N com validação de hierarquia)
        └── Tags (N:N para categorização)
```

### Regras de Associação

1. **Área → Exames**: 1:N (um área tem vários exames)
2. **Exame → Tópicos**: N:N com configurações específicas
3. **Tópico → Questões**: 1:N (uma questão pertence a um tópico)
4. **Exame → Questões**: N:N com validação hierárquica
5. **Questão → Tags**: N:N para categorização

## ⚖️ Regras de Negócio

### Validação Exame-Questão

**Implementada na camada de negócio** (`ExamQuestionUsecase`):

```go
func validateExamQuestionAssociation(examID, questionID int) error {
    // 1. Questão deve ter tópico atribuído
    // 2. Exame deve ter tópicos configurados
    // 3. Tópico da questão deve estar associado ao exame
    // 4. EXCEÇÃO: Exame com 1 tópico (100% peso) aceita qualquer questão
}
```

### Distribuição de Dificuldade

Cada associação **Exame-Tópico** define:
- **Quantidade de questões** do tópico no exame
- **Peso percentual** do tópico na nota final
- **Distribuição de dificuldade**:
  - % questões fáceis (padrão: 30%)
  - % questões médias (padrão: 50%)
  - % questões difíceis (padrão: 20%)

### Sistema de Tags

- **Tags simples**: apenas nome (sem descrição)
- **Categorização flexível**: questões podem ter múltiplas tags
- **Filtragem avançada**: busca por tags específicas

## 🚀 API Endpoints

### Áreas

```http
GET    /areas                 # Listar todas
GET    /areas/paginated       # Paginado com filtros
POST   /areas                 # Criar nova
PUT    /areas/{id}            # Atualizar
DELETE /areas/{id}            # Deletar
```

### Exames

```http
GET    /exams                 # Listar todos
GET    /exams/paginated       # Paginado com filtros
GET    /exams/area/{area_id}  # Por área
POST   /exams                 # Criar novo
PUT    /exams/{id}            # Atualizar
DELETE /exams/{id}            # Deletar
```

### Tópicos

```http
GET    /topics                # Listar todos
GET    /topics/paginated      # Paginado com filtros
GET    /topics/exam/{exam_id} # Por exame
POST   /topics                # Criar novo
PUT    /topics/{id}           # Atualizar
DELETE /topics/{id}           # Deletar
```

### Questões

```http
GET    /questions             # Listar todas
GET    /questions/paginated   # Paginado com filtros
GET    /questions/topic/{id}  # Por tópico
POST   /questions             # Criar nova
PUT    /questions/{id}        # Atualizar
DELETE /questions/{id}        # Deletar
POST   /questions/import      # Importar JSON
```

### Relacionamentos Exame-Tópico

```http
GET    /exams/{id}/topics                    # Tópicos do exame
POST   /exams/{id}/topics/{topic_id}         # Associar tópico
PUT    /exams/{id}/topics/{topic_id}         # Atualizar associação
DELETE /exams/{id}/topics/{topic_id}         # Desassociar
POST   /exams/{id}/topics/bulk               # Associação em lote
```

### Relacionamentos Exame-Questão

```http
GET    /exams/{id}/questions                 # Questões do exame
POST   /exams/{id}/questions/{question_id}   # Associar questão (com validação)
DELETE /exams/{id}/questions/{question_id}   # Desassociar
POST   /exams/{id}/questions/bulk            # Associação em lote
PUT    /exams/{id}/questions/reorder         # Reordenar questões
GET    /exams/{id}/questions/available       # Questões disponíveis
```

### Sistema de Tags

```http
# Gestão de Tags
GET    /tags                    # Listar todas
POST   /tags                    # Criar nova
PUT    /tags/{id}               # Atualizar
DELETE /tags/{id}               # Deletar
GET    /tags/search/{term}      # Buscar por nome

# Questões-Tags
GET    /questions/{id}/tags     # Tags da questão
PUT    /questions/{id}/tags     # Atualizar tags da questão
POST   /questions/{id}/tags/bulk # Associar múltiplas tags
DELETE /questions/{id}/tags/{tag_id} # Remover tag
```

## 📥 Importação de Questões

### Formato JSON

```json
{
  "area": {
    "name": "Tecnologia",
    "description": "Área de tecnologia e programação"
  },
  "exam": {
    "title": "Prova de JavaScript",
    "description": "Avaliação de conhecimentos em JavaScript",
    "max_time_minutes": 120,
    "passing_score": 70.0
  },
  "topic": {
    "name": "JavaScript Básico",
    "description": "Conceitos fundamentais",
    "weight_percentage": 100,
    "questions_count": 10,
    "difficulty_easy_percentage": 30,
    "difficulty_medium_percentage": 50,
    "difficulty_hard_percentage": 20
  },
  "questions": [
    {
      "text": "O que é uma variável em JavaScript?",
      "difficulty_level": "easy",
      "explanation": "Uma variável é um container para armazenar dados.",
      "tags": ["básico", "variáveis"],
      "options": [
        {
          "text": "Um container para dados",
          "is_correct": true
        },
        {
          "text": "Uma função",
          "is_correct": false
        }
      ]
    }
  ]
}
```

### Endpoint de Importação

```http
POST /questions/import
Content-Type: application/json

{
  "area": {...},
  "exam": {...},
  "topic": {...},
  "questions": [...]
}
```

## 🛠️ Instalação

### Pré-requisitos

- Go 1.21+
- PostgreSQL 14+
- Git

### Configuração do Banco

```bash
# Criar banco
createdb esimulate

# Executar schema
psql -d esimulate -f db_schema_v2.sql
```

### Executar Aplicação

```bash
# Instalar dependências
go mod tidy

# Executar
go run main.go
```

A aplicação estará disponível em `http://localhost:8080`

## ⚙️ Configuração

### Variáveis de Ambiente

```bash
# Banco de dados
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=esimulate

# Servidor
PORT=8080
JWT_SECRET=your-secret-key
```

### Estrutura de Configuração

```go
type Config struct {
    Database DatabaseConfig
    Server   ServerConfig
    JWT      JWTConfig
}
```

## 📈 Recursos Avançados

### Paginação

Todos os endpoints de listagem suportam paginação:

```http
GET /questions/paginated?page=1&pageSize=10&examId=1&topicId=2
```

### Filtros

- **Exames**: por área, status (ativo/inativo)
- **Tópicos**: por exame
- **Questões**: por exame, tópico, dificuldade, tags
- **Tags**: busca por nome

### Estatísticas

```http
GET /questions/{id}/stats     # Estatísticas da questão
GET /exams/{id}/stats         # Estatísticas do exame
GET /tags/{id}/stats          # Estatísticas da tag
```

### View Otimizada

```sql
-- View para consultas otimizadas
CREATE VIEW v_exam_questions_with_topics AS
SELECT 
    eq.exam_id,
    eq.question_id,
    eq.order_index,
    e.title as exam_title,
    q.text as question_text,
    q.difficulty_level,
    t.name as topic_name
FROM exam_questions eq
JOIN exams e ON eq.exam_id = e.id
JOIN questions q ON eq.question_id = q.id
JOIN topics t ON q.topic_id = t.id;
```

## 🧪 Testes

### Executar Testes

```bash
# Todos os testes
go test ./...

# Com coverage
go test -cover ./...

# Testes específicos
go test ./usecase/...
```

### Estrutura de Testes

```
backend/
├── domain/
├── usecase/
│   └── *_test.go
├── infra/
│   └── *_test.go
└── interfaces/
    └── *_test.go
```

## 📚 Documentação Adicional

- **Regras de Negócio**: Consulte `requisitos.md`
- **Migrações**: Scripts em `migration_*.sql`
- **Exemplos**: Dados de teste no schema

## 🤝 Contribuição

1. Fork o projeto
2. Crie sua feature branch (`git checkout -b feature/nova-funcionalidade`)
3. Commit suas mudanças (`git commit -m 'Add nova funcionalidade'`)
4. Push para a branch (`git push origin feature/nova-funcionalidade`)
5. Abra um Pull Request

## 📄 Licença

Este projeto está sob a licença MIT. Veja o arquivo `LICENSE` para mais detalhes.