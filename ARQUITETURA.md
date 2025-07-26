# Arquitetura do Sistema eSimulate

## Visão Geral
O eSimulate é um sistema para criação, realização e gerenciamento de simulados, voltado para aprendizado e preparação para concursos e provas de conhecimento. O backend foi desenvolvido em Go, seguindo Clean Architecture, Clean Code, princípios de segurança e conformidade com a LGPD.

## Camadas da Arquitetura

### 1. Domain (Domínio)
- Define as entidades centrais do sistema (Usuário, Tema, Questão, Simulado, etc).
- Contém interfaces (contracts) para repositórios, facilitando testes e desacoplamento.

### 2. Usecase (Casos de Uso)
- Implementa as regras de negócio do sistema.
- Orquestra operações entre entidades e repositórios.
- Responsável por validações, segurança (hash de senha, validação de entrada) e lógica de negócio.

### 3. Infra (Infraestrutura)
- Implementa os repositórios de persistência (PostgreSQL).
- Pode ser expandida para outros serviços (cache, email, etc).

### 4. Interfaces (Interface de Entrada)
- Handlers HTTP para cada recurso (usuário, tema, questão, simulado, etc).
- Middlewares de autenticação (JWT), CORS e tratamento de erros.
- Rotas organizadas e agrupadas por recurso.

## Segurança
- Senhas são armazenadas com hash seguro (bcrypt).
- Autenticação via JWT para rotas protegidas.
- Validação de entrada em todos os endpoints.
- Estrutura para login social (Google, Facebook).
- Middlewares para proteção de rotas e tratamento de autorização.

## LGPD
- Armazenamento mínimo de dados pessoais.
- Comentários e pontos de extensão para consentimento, anonimização e exclusão de dados.
- Possibilidade de exclusão de conta e dados do usuário.

## Banco de Dados
- Script SQL para criação das tabelas em PostgreSQL disponível em `backend/db_schema.sql`.
- Estrutura relacional cobre usuários, temas, questões, opções, simulados, histórico e desempenho.

## Testes
- Estrutura de testes automatizados para todos os endpoints.
- Mock de repositórios para facilitar testes unitários.

## Como Executar
1. Configure o PostgreSQL e crie o banco de dados `esimulate`.
2. Execute o script `backend/db_schema.sql` para criar as tabelas.
3. Compile e execute o backend:
   ```bash
   cd backend
   go run main.go
   ```
4. O servidor estará disponível em `http://localhost:8080`.

## Consumo dos Endpoints
- Todos os endpoints estão documentados no arquivo `backend/README.md`.
- Utilize ferramentas como Postman, Insomnia ou cURL para testar as rotas.
- Para rotas protegidas, obtenha um token JWT via login e envie no header `Authorization: Bearer <token>`.

## Expansão
- O sistema está pronto para expansão de funcionalidades, integração com frontend e melhorias contínuas.
- Novos repositórios, serviços e integrações podem ser adicionados facilmente devido à separação de camadas.

---

Este documento será atualizado conforme o sistema evoluir.
