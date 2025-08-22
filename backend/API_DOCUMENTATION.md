# API Documentation - eSimulate

Documentação completa da API REST do sistema eSimulate.

## 📋 Índice

- [Autenticação](#autenticação)
- [Formato de Resposta](#formato-de-resposta)
- [Endpoints](#endpoints)
  - [Áreas](#áreas)
  - [Exames](#exames)
  - [Tópicos](#tópicos)
  - [Questões](#questões)
  - [Tags](#tags)
  - [Relacionamentos](#relacionamentos)
- [Códigos de Status](#códigos-de-status)
- [Exemplos de Uso](#exemplos-de-uso)

## 🔐 Autenticação

A API utiliza JWT (JSON Web Tokens) para autenticação.

### Login

```http
POST /auth/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "password"
}
```

**Resposta:**
```json
{
  "message": "Login successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "name": "João Silva",
      "email": "user@example.com",
      "role": "admin"
    }
  }
}
```

### Usando o Token

Inclua o token no header de todas as requisições protegidas:

```http
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

## 📤 Formato de Resposta

Todas as respostas seguem o padrão:

### Sucesso
```json
{
  "message": "Operação realizada com sucesso",
  "data": { /* dados da resposta */ }
}
```

### Erro
```json
{
  "error": "Descrição do erro",
  "details": "Detalhes adicionais (opcional)"
}
```

### Paginação
```json
{
  "message": "Dados recuperados com sucesso",
  "data": [/* array de itens */],
  "pagination": {
    "page": 1,
    "pageSize": 10,
    "total": 150,
    "totalPages": 15
  }
}
```

## 🚀 Endpoints

### Áreas

#### Listar Todas as Áreas
```http
GET /areas
Authorization: Bearer {token}
```

**Resposta:**
```json
{
  "message": "Áreas recuperadas com sucesso",
  "data": [
    {
      "id": 1,
      "name": "Tecnologia",
      "description": "Área de tecnologia e programação",
      "created_at": "2024-01-15T10:00:00Z",
      "updated_at": "2024-01-15T10:00:00Z"
    }
  ]
}
```

#### Listar Áreas Paginadas
```http
GET /areas/paginated?page=1&pageSize=10&search=tech
Authorization: Bearer {token}
```

**Parâmetros:**
- `page` (int): Página atual (padrão: 1)
- `pageSize` (int): Itens por página (padrão: 10)
- `search` (string): Termo de busca (opcional)

#### Criar Área
```http
POST /areas
Content-Type: application/json
Authorization: Bearer {token}

{
  "name": "Matemática",
  "description": "Área de ciências exatas"
}
```

#### Atualizar Área
```http
PUT /areas/{id}
Content-Type: application/json
Authorization: Bearer {token}

{
  "name": "Matemática Avançada",
  "description": "Área de matemática de nível superior"
}
```

#### Deletar Área
```http
DELETE /areas/{id}
Authorization: Bearer {token}
```

### Exames

#### Listar Todos os Exames
```http
GET /exams
Authorization: Bearer {token}
```

#### Listar Exames Paginados
```http
GET /exames/paginated?page=1&pageSize=10&areaId=1&status=active
Authorization: Bearer {token}
```

**Parâmetros:**
- `page` (int): Página atual
- `pageSize` (int): Itens por página
- `areaId` (int): Filtrar por área (opcional)
- `status` (string): `active`, `inactive`, `all` (opcional)

#### Listar Exames por Área
```http
GET /exams/area/{area_id}
Authorization: Bearer {token}
```

#### Criar Exame
```http
POST /exams
Content-Type: application/json
Authorization: Bearer {token}

{
  "area_id": 1,
  "title": "Prova de JavaScript",
  "description": "Avaliação de conhecimentos básicos",
  "max_time_minutes": 120,
  "passing_score": 70.0,
  "is_active": true
}
```

#### Atualizar Exame
```http
PUT /exams/{id}
Content-Type: application/json
Authorization: Bearer {token}

{
  "title": "Prova de JavaScript Avançado",
  "max_time_minutes": 180,
  "passing_score": 75.0
}
```

### Tópicos

#### Listar Todos os Tópicos
```http
GET /topics
Authorization: Bearer {token}
```

#### Listar Tópicos Paginados
```http
GET /topics/paginated?page=1&pageSize=10&examId=1
Authorization: Bearer {token}
```

**Parâmetros:**
- `examId` (int): Filtrar por exame (opcional)

#### Listar Tópicos por Exame
```http
GET /topics/exam/{exam_id}
Authorization: Bearer {token}
```

#### Criar Tópico
```http
POST /topics
Content-Type: application/json
Authorization: Bearer {token}

{
  "name": "JavaScript Básico",
  "description": "Conceitos fundamentais da linguagem"
}
```

### Questões

#### Listar Questões Paginadas
```http
GET /questions/paginated?page=1&pageSize=10&examId=1&topicId=2&difficulty=easy
Authorization: Bearer {token}
```

**Parâmetros:**
- `examId` (int): Filtrar por exame (opcional)
- `topicId` (int): Filtrar por tópico (opcional)
- `difficulty` (string): `easy`, `medium`, `hard` (opcional)

#### Listar Questões por Tópico
```http
GET /questions/topic/{topic_id}
Authorization: Bearer {token}
```

#### Criar Questão
```http
POST /questions
Content-Type: application/json
Authorization: Bearer {token}

{
  "topic_id": 1,
  "text": "O que é uma variável em JavaScript?",
  "difficulty_level": "easy",
  "explanation": "Uma variável é um container para armazenar dados",
  "is_active": true,
  "options": [
    {
      "text": "Um container para dados",
      "is_correct": true
    },
    {
      "text": "Uma função",
      "is_correct": false
    },
    {
      "text": "Um objeto",
      "is_correct": false
    },
    {
      "text": "Um método",
      "is_correct": false
    }
  ]
}
```

#### Importar Questões
```http
POST /questions/import
Content-Type: application/json
Authorization: Bearer {token}

{
  "area": {
    "name": "Tecnologia",
    "description": "Área de tecnologia"
  },
  "exam": {
    "title": "Prova de JavaScript",
    "description": "Avaliação básica",
    "max_time_minutes": 120,
    "passing_score": 70.0
  },
  "topic": {
    "name": "JavaScript Básico",
    "description": "Conceitos fundamentais",
    "weight_percentage": 100,
    "questions_count": 5,
    "difficulty_easy_percentage": 40,
    "difficulty_medium_percentage": 40,
    "difficulty_hard_percentage": 20
  },
  "questions": [
    {
      "text": "O que é JavaScript?",
      "difficulty_level": "easy",
      "explanation": "JavaScript é uma linguagem de programação",
      "tags": ["básico", "linguagem"],
      "options": [
        {
          "text": "Uma linguagem de programação",
          "is_correct": true
        },
        {
          "text": "Um banco de dados",
          "is_correct": false
        }
      ]
    }
  ]
}
```

### Tags

#### Listar Todas as Tags
```http
GET /tags
Authorization: Bearer {token}
```

**Resposta:**
```json
{
  "message": "Tags recuperadas com sucesso",
  "data": [
    {
      "id": 1,
      "name": "JavaScript",
      "created_at": "2024-01-15T10:00:00Z"
    },
    {
      "id": 2,
      "name": "Básico",
      "created_at": "2024-01-15T10:00:00Z"
    }
  ]
}
```

#### Criar Tag
```http
POST /tags
Content-Type: application/json
Authorization: Bearer {token}

{
  "name": "React"
}
```

#### Atualizar Tag
```http
PUT /tags/{id}
Content-Type: application/json
Authorization: Bearer {token}

{
  "name": "React.js"
}
```

#### Deletar Tag
```http
DELETE /tags/{id}
Authorization: Bearer {token}
```

#### Buscar Tags
```http
GET /tags/search/{term}
Authorization: Bearer {token}
```

### Relacionamentos

#### Exame-Tópico

##### Listar Tópicos do Exame
```http
GET /exams/{exam_id}/topics
Authorization: Bearer {token}
```

**Resposta:**
```json
{
  "message": "Tópicos do exame recuperados com sucesso",
  "data": [
    {
      "exam_id": 1,
      "topic_id": 1,
      "topic_name": "JavaScript Básico",
      "questions_count": 10,
      "weight_percentage": 60.0,
      "order_index": 1,
      "difficulty_easy_percentage": 30.0,
      "difficulty_medium_percentage": 50.0,
      "difficulty_hard_percentage": 20.0
    }
  ]
}
```

##### Associar Tópico ao Exame
```http
POST /exams/{exam_id}/topics/{topic_id}
Content-Type: application/json
Authorization: Bearer {token}

{
  "questions_count": 15,
  "weight_percentage": 70.0,
  "order_index": 1,
  "difficulty_easy_percentage": 40.0,
  "difficulty_medium_percentage": 40.0,
  "difficulty_hard_percentage": 20.0
}
```

##### Atualizar Associação Exame-Tópico
```http
PUT /exams/{exam_id}/topics/{topic_id}
Content-Type: application/json
Authorization: Bearer {token}

{
  "questions_count": 20,
  "weight_percentage": 80.0
}
```

##### Desassociar Tópico do Exame
```http
DELETE /exams/{exam_id}/topics/{topic_id}
Authorization: Bearer {token}
```

##### Associação em Lote
```http
POST /exams/{exam_id}/topics/bulk
Content-Type: application/json
Authorization: Bearer {token}

{
  "topics": [
    {
      "topic_id": 1,
      "questions_count": 10,
      "weight_percentage": 50.0,
      "order_index": 1
    },
    {
      "topic_id": 2,
      "questions_count": 10,
      "weight_percentage": 50.0,
      "order_index": 2
    }
  ]
}
```

#### Exame-Questão

##### Listar Questões do Exame
```http
GET /exams/{exam_id}/questions
Authorization: Bearer {token}
```

##### Associar Questão ao Exame
```http
POST /exams/{exam_id}/questions/{question_id}
Authorization: Bearer {token}
```

**Nota:** A associação é validada automaticamente:
- Questão deve ter tópico atribuído
- Tópico da questão deve estar associado ao exame
- Exceção: Exames com 1 tópico (100% peso) aceitam qualquer questão

##### Listar Questões Disponíveis
```http
GET /exams/{exam_id}/questions/available
Authorization: Bearer {token}
```

##### Associação em Lote
```http
POST /exams/{exam_id}/questions/bulk
Content-Type: application/json
Authorization: Bearer {token}

{
  "question_ids": [1, 2, 3, 4, 5]
}
```

##### Reordenar Questões
```http
PUT /exams/{exam_id}/questions/reorder
Content-Type: application/json
Authorization: Bearer {token}

{
  "questions": [
    {"question_id": 1, "order_index": 1},
    {"question_id": 2, "order_index": 2},
    {"question_id": 3, "order_index": 3}
  ]
}
```

#### Questão-Tag

##### Listar Tags da Questão
```http
GET /questions/{question_id}/tags
Authorization: Bearer {token}
```

##### Atualizar Tags da Questão
```http
PUT /questions/{question_id}/tags
Content-Type: application/json
Authorization: Bearer {token}

{
  "tag_ids": [1, 2, 3]
}
```

##### Associar Múltiplas Tags
```http
POST /questions/{question_id}/tags/bulk
Content-Type: application/json
Authorization: Bearer {token}

{
  "tag_ids": [4, 5, 6]
}
```

##### Remover Tag da Questão
```http
DELETE /questions/{question_id}/tags/{tag_id}
Authorization: Bearer {token}
```

## 📊 Códigos de Status

### Sucesso
- `200 OK` - Operação realizada com sucesso
- `201 Created` - Recurso criado com sucesso
- `204 No Content` - Operação realizada, sem conteúdo de retorno

### Erro do Cliente
- `400 Bad Request` - Dados inválidos ou malformados
- `401 Unauthorized` - Token de autenticação inválido ou ausente
- `403 Forbidden` - Acesso negado (permissões insuficientes)
- `404 Not Found` - Recurso não encontrado
- `409 Conflict` - Conflito (ex: nome já existe)
- `422 Unprocessable Entity` - Dados válidos, mas regras de negócio violadas

### Erro do Servidor
- `500 Internal Server Error` - Erro interno do servidor
- `503 Service Unavailable` - Serviço temporariamente indisponível

## 🔧 Exemplos de Uso

### Cenário Completo: Criar Prova de JavaScript

#### 1. Criar Área
```bash
curl -X POST http://localhost:8080/areas \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Tecnologia",
    "description": "Área de tecnologia e programação"
  }'
```

#### 2. Criar Exame
```bash
curl -X POST http://localhost:8080/exams \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "area_id": 1,
    "title": "Prova de JavaScript",
    "description": "Avaliação de conhecimentos em JavaScript",
    "max_time_minutes": 120,
    "passing_score": 70.0
  }'
```

#### 3. Criar Tópicos
```bash
# Tópico 1
curl -X POST http://localhost:8080/topics \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "JavaScript Básico",
    "description": "Conceitos fundamentais"
  }'

# Tópico 2
curl -X POST http://localhost:8080/topics \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "JavaScript Avançado",
    "description": "Conceitos avançados"
  }'
```

#### 4. Associar Tópicos ao Exame
```bash
# Associar tópico 1 (60% da prova)
curl -X POST http://localhost:8080/exams/1/topics/1 \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "questions_count": 12,
    "weight_percentage": 60.0,
    "order_index": 1,
    "difficulty_easy_percentage": 40.0,
    "difficulty_medium_percentage": 40.0,
    "difficulty_hard_percentage": 20.0
  }'

# Associar tópico 2 (40% da prova)
curl -X POST http://localhost:8080/exams/1/topics/2 \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "questions_count": 8,
    "weight_percentage": 40.0,
    "order_index": 2,
    "difficulty_easy_percentage": 30.0,
    "difficulty_medium_percentage": 50.0,
    "difficulty_hard_percentage": 20.0
  }'
```

#### 5. Criar Tags
```bash
curl -X POST http://localhost:8080/tags \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name": "variáveis"}'

curl -X POST http://localhost:8080/tags \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name": "funções"}'
```

#### 6. Importar Questões via JSON
```bash
curl -X POST http://localhost:8080/questions/import \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d @prova_javascript.json
```

#### 7. Verificar Questões do Exame
```bash
curl -X GET http://localhost:8080/exams/1/questions \
  -H "Authorization: Bearer $TOKEN"
```

### Exemplo de Filtros Avançados

```bash
# Buscar questões de JavaScript básico, dificuldade fácil, com tag "variáveis"
curl -X GET "http://localhost:8080/questions/paginated?page=1&pageSize=5&topicId=1&difficulty=easy&tags=variáveis" \
  -H "Authorization: Bearer $TOKEN"

# Buscar exames ativos da área de tecnologia
curl -X GET "http://localhost:8080/exams/paginated?areaId=1&status=active&page=1&pageSize=10" \
  -H "Authorization: Bearer $TOKEN"
```

## 🔍 Dicas de Uso

### Performance
- Use paginação para listas grandes
- Filtre resultados sempre que possível
- Cache tokens JWT no cliente

### Validação
- Sempre valide dados no frontend antes de enviar
- Trate todos os códigos de erro apropriadamente
- Use a validação automática de exame-questão

### Segurança
- Nunca exponha tokens no código cliente
- Implemente refresh tokens para sessões longas
- Valide permissões no backend

### Debugging
- Use logs detalhados para troubleshooting
- Monitore códigos de status HTTP
- Implemente retry logic para falhas temporárias
