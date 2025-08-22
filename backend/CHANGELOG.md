# Changelog

Todas as mudanças notáveis neste projeto serão documentadas neste arquivo.

O formato é baseado em [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
e este projeto adere ao [Versionamento Semântico](https://semver.org/spec/v2.0.0.html).

## [2.0.0] - 2024-01-15

### 🚀 Adicionado
- **Sistema de Relacionamentos N:N**: Implementação completa de relacionamentos entre exames-tópicos e exames-questões
- **Sistema de Tags**: Sistema flexível de categorização para questões com relacionamento N:N
- **Importação de Questões**: Endpoint para importar estruturas completas via JSON
- **Validação de Negócio**: Validação hierárquica exame-questão implementada na camada de aplicação
- **Paginação Universal**: Todos os endpoints de listagem agora suportam paginação
- **Filtros Avançados**: Filtros específicos por entidade (área, status, dificuldade, tags)
- **Documentação Completa**: README, API Documentation, Requirements e Changelog

### 🔄 Modificado
- **Arquitetura de Tópicos**: Tópicos agora são entidades independentes (não pertencem a exame específico)
- **Estrutura de Exames**: Exames agora têm relacionamento N:N com tópicos incluindo configurações de peso
- **Validação**: Movida do banco de dados (triggers) para a camada de aplicação (use cases)
- **Tags Simplificadas**: Removido campo 'description', tags agora têm apenas nome
- **Schema do Banco**: Atualizado para v2.3 com novas tabelas e relacionamentos

### 🗑️ Removido
- **Campos Obsoletos**: Removidos campos de tópicos que eram específicos de exames
- **Triggers de Validação**: Removidas validações do banco, movidas para aplicação
- **Campo Description**: Removido das tags para simplificar o modelo

### 🐛 Corrigido
- **Proxy Issues**: Corrigidos problemas de roteamento no proxy do frontend
- **Null Pointer Errors**: Corrigidos erros de ponteiro nulo em repositórios
- **SQL Parameter Errors**: Corrigidos problemas de contagem de parâmetros SQL
- **NULL Handling**: Implementado tratamento adequado de campos NULL no banco
- **Git Issues**: Removido node_modules do controle de versão

## [1.0.0] - 2024-01-01

### 🚀 Inicial
- **CRUD Básico**: Implementação de áreas, exames, tópicos e questões
- **Autenticação**: Sistema JWT para autenticação
- **Clean Architecture**: Estrutura baseada em Clean Architecture
- **PostgreSQL**: Integração com banco PostgreSQL
- **API REST**: Endpoints básicos para todas as entidades

---

## Detalhes das Versões

### [2.0.0] - Relacionamentos N:N e Sistema de Tags

#### 🏗️ Mudanças na Arquitetura

**Antes (v1.0):**
```
Área → Exame → Tópico → Questão
```

**Depois (v2.0):**
```
Área → Exame ←→ Tópico (N:N com configurações)
Tópico → Questão ←→ Exame (N:N com validação)
Questão ←→ Tag (N:N para categorização)
```

#### 📊 Novas Tabelas

1. **exam_topics** - Relacionamento N:N exames-tópicos
   ```sql
   - exam_id, topic_id (PK composta)
   - questions_count, weight_percentage
   - difficulty distribution (easy/medium/hard %)
   - order_index
   ```

2. **exam_questions** - Relacionamento N:N exames-questões
   ```sql
   - exam_id, question_id (PK composta)
   - order_index
   ```

3. **question_tags** - Sistema de tags
   ```sql
   - id, name (único)
   - created_at
   ```

4. **question_tag_associations** - Relacionamento N:N questões-tags
   ```sql
   - question_id, tag_id (PK composta)
   ```

#### 🔧 Novos Endpoints

**Relacionamentos Exame-Tópico:**
- `GET /exams/{id}/topics` - Listar tópicos do exame
- `POST /exams/{id}/topics/{topic_id}` - Associar tópico
- `PUT /exams/{id}/topics/{topic_id}` - Atualizar associação
- `DELETE /exams/{id}/topics/{topic_id}` - Desassociar
- `POST /exams/{id}/topics/bulk` - Associação em lote

**Relacionamentos Exame-Questão:**
- `GET /exams/{id}/questions` - Listar questões do exame
- `POST /exams/{id}/questions/{question_id}` - Associar questão (com validação)
- `DELETE /exams/{id}/questions/{question_id}` - Desassociar
- `POST /exams/{id}/questions/bulk` - Associação em lote
- `PUT /exams/{id}/questions/reorder` - Reordenar questões
- `GET /exams/{id}/questions/available` - Questões disponíveis

**Sistema de Tags:**
- `GET /tags` - Listar todas
- `POST /tags` - Criar nova
- `PUT /tags/{id}` - Atualizar
- `DELETE /tags/{id}` - Deletar
- `GET /tags/search/{term}` - Buscar por nome
- `GET /questions/{id}/tags` - Tags da questão
- `PUT /questions/{id}/tags` - Atualizar tags
- `POST /questions/{id}/tags/bulk` - Associar múltiplas
- `DELETE /questions/{id}/tags/{tag_id}` - Remover tag

**Importação:**
- `POST /questions/import` - Importar estrutura completa via JSON

#### ⚖️ Regras de Negócio Implementadas

1. **Validação Hierárquica Exame-Questão:**
   - Questão deve ter tópico atribuído
   - Tópico da questão deve estar associado ao exame
   - **Exceção:** Exame com único tópico (peso 100%) aceita qualquer questão

2. **Distribuição de Pesos:**
   - Soma dos pesos dos tópicos por exame deve ser 100%
   - Distribuição de dificuldade por tópico deve somar 100%

3. **Tags Simplificadas:**
   - Nome único obrigatório
   - Sem campo description
   - Relacionamento N:N flexível

#### 🔧 Melhorias Técnicas

1. **Paginação Universal:**
   ```http
   GET /entity/paginated?page=1&pageSize=10&filter=value
   ```

2. **Filtros Específicos:**
   - Áreas: busca por nome
   - Exames: área, status (ativo/inativo)
   - Tópicos: exame
   - Questões: exame, tópico, dificuldade, tags

3. **Validação na Aplicação:**
   - Movida do banco (triggers) para use cases
   - Mensagens de erro específicas
   - Melhor testabilidade

4. **View Otimizada:**
   ```sql
   CREATE VIEW v_exam_questions_with_topics AS
   -- Join otimizado para consultas complexas
   ```

#### 📦 Estrutura de Importação JSON

```json
{
  "area": {
    "name": "Tecnologia",
    "description": "Área de tecnologia"
  },
  "exam": {
    "title": "Prova de JavaScript",
    "max_time_minutes": 120,
    "passing_score": 70.0
  },
  "topic": {
    "name": "JavaScript Básico",
    "weight_percentage": 100,
    "questions_count": 10,
    "difficulty_easy_percentage": 30,
    "difficulty_medium_percentage": 50,
    "difficulty_hard_percentage": 20
  },
  "questions": [
    {
      "text": "O que é JavaScript?",
      "difficulty_level": "easy",
      "explanation": "JavaScript é uma linguagem de programação",
      "tags": ["básico", "linguagem"],
      "options": [
        {"text": "Uma linguagem", "is_correct": true},
        {"text": "Um banco de dados", "is_correct": false}
      ]
    }
  ]
}
```

#### 🧪 Casos de Teste Implementados

1. **Validação Hierárquica:**
   ```
   ✅ Questão de tópico associado → Permitido
   ❌ Questão de tópico não associado → Rejeitado
   ✅ Exame com 1 tópico 100% → Qualquer questão permitida
   ```

2. **Importação JSON:**
   ```
   ✅ Estrutura válida → Sucesso
   ❌ Dados inválidos → Erro específico
   ✅ Entidades existentes → Reutilizadas
   ```

3. **Paginação e Filtros:**
   ```
   ✅ Página inválida → Página 1
   ✅ PageSize > 100 → Limitado a 100
   ✅ Filtros combinados → Funcionando
   ```

### [1.0.0] - Versão Inicial

#### 🏗️ Arquitetura Base
- **Clean Architecture** com separação clara de camadas
- **Domain-Driven Design** com entidades bem definidas
- **Repository Pattern** para acesso a dados
- **Use Case Pattern** para regras de negócio

#### 📊 Entidades Básicas
1. **Áreas** - Categorização principal
2. **Exames** - Avaliações por área
3. **Tópicos** - Temas específicos (pertenciam a exames)
4. **Questões** - Perguntas com opções múltiplas

#### 🔧 Funcionalidades Básicas
- CRUD completo para todas as entidades
- Autenticação JWT
- Validações básicas
- API REST padronizada

---

## 🚀 Próximas Versões

### [2.1.0] - Planejado
- [ ] Interface web completa
- [ ] Aplicação de simulados
- [ ] Correção automática
- [ ] Relatórios básicos

### [3.0.0] - Futuro
- [ ] Análise de desempenho
- [ ] Integração com LMS
- [ ] IA para geração de questões
- [ ] Banco de questões público

---

## 📝 Convenções

### Tipos de Mudança
- **🚀 Adicionado** - Para novas funcionalidades
- **🔄 Modificado** - Para mudanças em funcionalidades existentes
- **🗑️ Removido** - Para funcionalidades removidas
- **🐛 Corrigido** - Para correções de bugs
- **🔒 Segurança** - Para correções de vulnerabilidades

### Versionamento
- **MAJOR** (X.0.0) - Mudanças incompatíveis na API
- **MINOR** (x.Y.0) - Novas funcionalidades compatíveis
- **PATCH** (x.y.Z) - Correções de bugs compatíveis
