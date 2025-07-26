# Backend eSimulate

## Endpoints planejados

### Autenticação e Usuários
- POST /auth/register — Cadastro de usuário
- POST /auth/login — Login tradicional
- POST /auth/google — Login via Google
- POST /auth/facebook — Login via Facebook
- GET /users/me — Dados do usuário autenticado
- PUT /users/me — Atualizar dados do usuário

### Temas e Subtemas
- GET /themes — Listar temas e subtemas
- POST /themes — Criar tema/subtema
- PUT /themes/{id} — Atualizar tema/subtema
- DELETE /themes/{id} — Remover tema/subtema

### Questões
- GET /questions — Listar questões (com filtros por tema, etc)
- POST /questions — Criar questão
- GET /questions/{id} — Detalhar questão
- PUT /questions/{id} — Atualizar questão
- DELETE /questions/{id} — Remover questão

### Simulados
- POST /quizzes — Criar simulado
- GET /quizzes — Listar simulados
- GET /quizzes/{id} — Detalhar simulado
- POST /quizzes/{id}/start — Iniciar simulado
- POST /quizzes/{id}/answer — Responder questão do simulado
- GET /quizzes/{id}/result — Resultado do simulado

### Histórico e Relatórios
- GET /history — Histórico de simulados realizados pelo usuário
- GET /performance — Gráficos e relatórios de desempenho

---

A estrutura do backend seguirá Clean Architecture, Clean Code e boas práticas, com comentários em português para facilitar o entendimento.
