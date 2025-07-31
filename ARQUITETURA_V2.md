# Arquitetura do Sistema eSimulate v2

## 🎯 **Visão Geral**
O eSimulate v2 é um sistema completo para criação, realização e gerenciamento de provas e simulados, com níveis de acesso diferenciados e análise detalhada de desempenho. O sistema foi desenvolvido seguindo Clean Architecture, Clean Code, princípios de segurança e conformidade com a LGPD.

## 👥 **Níveis de Acesso**

### 🔐 **Usuário (user)**
- Aplicar provas disponíveis
- Visualizar histórico de tentativas
- Analisar desempenho por domínios
- Marcar questões para revisão durante a prova

### ✍️ **Redator (redator)**
- Criar e gerenciar provas
- Cadastrar questões e alternativas
- Definir domínios e pesos
- Configurar áreas de conhecimento

### ⚙️ **Administrador (admin)**
- Acesso completo ao sistema
- Gerenciar usuários e níveis de acesso
- Configurar áreas de conhecimento
- Visualizar relatórios gerais

## 🏗️ **Camadas da Arquitetura**

### 1. Domain (Domínio)
- Define as entidades centrais do sistema
- Contém interfaces (contracts) para repositórios
- Implementa regras de negócio básicas

**Entidades Principais:**
- `User` - Usuários com níveis de acesso
- `Area` - Áreas de conhecimento (TI, OAB, CRC, etc.)
- `Exam` - Provas/exames com tempo e pontuação
- `Domain` - Domínios da prova com pesos
- `Question` - Questões com múltiplas alternativas
- `UserExam` - Aplicação de provas por usuários
- `UserAnswer` - Respostas dos usuários
- `DomainPerformance` - Desempenho por domínio

### 2. Usecase (Casos de Uso)
- Implementa as regras de negócio do sistema
- Orquestra operações entre entidades e repositórios
- Responsável por validações e segurança

**Use Cases Principais:**
- `AreaUsecase` - Gerenciamento de áreas
- `ExamUsecase` - Gerenciamento de provas
- `UserExamUsecase` - Aplicação de provas
- `QuestionUsecase` - Gerenciamento de questões

### 3. Infra (Infraestrutura)
- Implementa os repositórios de persistência (PostgreSQL)
- Pode ser expandida para outros serviços

### 4. Interfaces (Interface de Entrada)
- Handlers HTTP para cada recurso
- Middlewares de autenticação (JWT), CORS e tratamento de erros
- Rotas organizadas e agrupadas por recurso

## 📊 **Estrutura de Provas**

### 🎯 **Áreas de Conhecimento**
- TI - Tecnologia da Informação
- OAB - Ordem dos Advogados do Brasil
- CRC - Conselho Regional de Contabilidade
- Concurso Petrobras
- Vestibular
- Concurso Público

### 📝 **Exemplo: AWS Certified Cloud Practitioner**
```
Prova: AWS Certified Cloud Practitioner (CLF-C02)
Tempo: 90 minutos
Pontuação mínima: 70%

Domínios:
- Cloud Concepts (24%)
- Security and Compliance (30%)
- Cloud Technology and Services (34%)
- Billing, Pricing, and Support (12%)
```

### 🧩 **Estrutura de Questões**
- **Enunciado:** Texto da questão
- **Alternativas:** 5 ou mais opções
- **Explicação:** Justificativa da resposta correta
- **Domínio:** Pertence a um domínio específico
- **Dificuldade:** 1=fácil, 2=médio, 3=difícil
- **Revisão:** Pode ser marcada para revisão

## 🔄 **Fluxo de Aplicação de Provas**

### 1. **Início da Prova**
- Usuário seleciona uma prova disponível
- Sistema inicia cronômetro
- Questões são apresentadas sequencialmente

### 2. **Durante a Prova**
- Usuário responde questões
- Pode marcar para revisão
- Pode pular e voltar depois
- Sistema salva respostas automaticamente

### 3. **Revisão**
- Questões marcadas são destacadas
- Usuário pode alterar respostas
- Questões não respondidas são identificadas

### 4. **Finalização**
- Sistema calcula pontuação total
- Analisa desempenho por domínios
- Identifica áreas que precisam melhorar
- Gera relatório detalhado

## 📈 **Análise de Desempenho**

### 🎯 **Métricas por Domínio**
- **Questões respondidas:** Quantidade total
- **Respostas corretas:** Quantidade de acertos
- **Percentual de acerto:** (Corretas/Total) * 100
- **Necessita melhorar:** Se < 70% no domínio

### 📊 **Relatórios Disponíveis**
- **Histórico de tentativas:** Todas as provas realizadas
- **Progresso temporal:** Evolução ao longo do tempo
- **Análise por domínio:** Pontos fortes e fracos
- **Comparação:** Desempenho vs. média

## 🔐 **Segurança**

### 🛡️ **Autenticação e Autorização**
- Senhas com hash seguro (bcrypt)
- Autenticação via JWT para rotas protegidas
- Validação de níveis de acesso
- Middlewares de proteção

### 📋 **LGPD**
- Armazenamento mínimo de dados pessoais
- Possibilidade de exclusão de conta
- Anonimização de dados históricos
- Controle de consentimento

## 🗄️ **Banco de Dados**

### 📋 **Tabelas Principais**
```sql
-- Usuários com níveis de acesso
users (id, name, email, password_hash, role, created_at, updated_at)

-- Áreas de conhecimento
areas (id, name, description, created_at)

-- Provas/exames
exams (id, title, description, area_id, max_time_minutes, passing_score, is_active, created_by, created_at, updated_at)

-- Domínios da prova
domains (id, exam_id, name, description, weight_percentage, order_index, created_at)

-- Questões
questions (id, exam_id, domain_id, statement, explanation, difficulty_level, created_by, is_active, created_at, updated_at)

-- Opções de resposta
options (id, question_id, text, is_correct, explanation, order_index)

-- Aplicação de provas
user_exams (id, user_id, exam_id, started_at, finished_at, total_score, passed, time_spent_minutes)

-- Respostas dos usuários
user_answers (id, user_exam_id, question_id, option_id, is_correct, is_marked_for_review, answered_at)

-- Desempenho por domínio
domain_performance (id, user_exam_id, domain_id, questions_answered, correct_answers, score_percentage, needs_improvement)
```

## 🚀 **Como Executar**

### 1. **Configuração do Banco**
```bash
# Crie o banco PostgreSQL
createdb esimulate

# Execute o script de criação
psql -d esimulate -f backend/db_schema_v2.sql
```

### 2. **Backend**
```bash
cd backend
go run main.go
```

### 3. **Frontend**
```bash
cd frontend
npm start
```

## 📡 **Endpoints da API**

### 🔐 **Autenticação**
- `POST /auth/login` - Login de usuário
- `POST /auth/register` - Registro de usuário

### 👥 **Usuários**
- `GET /users` - Lista usuários (admin)
- `GET /users/{id}` - Busca usuário por ID

### 📚 **Áreas**
- `GET /areas` - Lista todas as áreas
- `POST /areas/create` - Cria nova área (redator/admin)
- `PUT /areas/update` - Atualiza área (redator/admin)
- `DELETE /areas/delete` - Remove área (admin)
- `GET /areas/detail` - Busca área por ID

### 📝 **Exames**
- `GET /exams` - Lista todos os exames
- `GET /exams/active` - Lista exames ativos
- `POST /exams/create` - Cria novo exame (redator/admin)
- `PUT /exams/update` - Atualiza exame (redator/admin)
- `DELETE /exams/delete` - Remove exame (admin)
- `GET /exams/detail` - Busca exame por ID

### 🎯 **Aplicação de Exames**
- `POST /user-exams/start` - Inicia exame
- `POST /user-exams/submit-answer` - Submete resposta
- `POST /user-exams/finish` - Finaliza exame
- `GET /user-exams/detail` - Busca exame aplicado
- `GET /user-exams/list` - Lista exames do usuário

### 📊 **Histórico e Performance**
- `GET /history` - Histórico de tentativas
- `GET /performance` - Relatório de desempenho

## 🧪 **Testes**

### 📋 **Usuários de Teste**
```bash
# Administrador
Email: admin@esimulate.com
Senha: password

# Redator
Email: redator@esimulate.com
Senha: password

# Usuário
Email: user@esimulate.com
Senha: password
```

## 🔄 **Expansão Futura**

### 🚀 **Funcionalidades Planejadas**
- **Simulados em tempo real:** Múltiplos usuários simultaneamente
- **Gamificação:** Pontos, badges, rankings
- **Relatórios avançados:** Gráficos e análises detalhadas
- **Integração com LMS:** Moodle, Canvas, etc.
- **API pública:** Para integração com outros sistemas
- **Mobile app:** Aplicativo nativo para iOS/Android

### 🔧 **Melhorias Técnicas**
- **Cache Redis:** Para melhor performance
- **Microserviços:** Separação por domínio
- **Event sourcing:** Para auditoria completa
- **Kubernetes:** Orquestração de containers
- **CI/CD:** Pipeline automatizado

---

**Versão:** 2.0  
**Data:** 2024  
**Autor:** Equipe eSimulate 