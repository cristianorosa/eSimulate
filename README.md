# eSimulate

Sistema para simulados de questões objetivas

## Objetivo
O eSimulate é um sistema para criação, realização e gerenciamento de simulados, voltado para aprendizado e preparação para concursos e provas de conhecimento. O objetivo é proporcionar uma plataforma eficiente para estudantes e concurseiros praticarem questões, acompanharem seu desempenho e evoluírem nos estudos.

## Estrutura do Projeto
O projeto está organizado em duas principais pastas:

- `backend/`: Responsável pela lógica de negócio, API, banco de dados e autenticação.
- `frontend/`: Responsável pela interface do usuário, onde os simulados são realizados e os resultados visualizados.

## Backend
- Linguagem: Go (Golang), utilizando os recursos mais atuais.
- Arquitetura: Clean Architecture, com aplicação de Clean Code e os Design Patterns mais adequados para cada contexto.
- Comentários no código em português, facilitando o entendimento e manutenção.
- Banco de dados inicial: PostgreSQL (com possibilidade de alteração futura).
- Principais funcionalidades:
  - Cadastro de questões e suas opções de resposta, incluindo explicação do porquê das opções estarem certas ou erradas.
  - Cadastro de usuários, login tradicional e login via conta Google e/ou Facebook.
  - Realização de simulados por temas e subtemas (cada tema pode conter subdivisões).
  - Histórico dos simulados realizados por cada usuário.
  - Geração de gráficos e relatórios de desempenho para acompanhamento da evolução.

## Frontend
- Interface para realização dos simulados, visualização de resultados, gráficos e relatórios.
- Cadastro e autenticação de usuários.
- Navegação por temas e subtemas.

## Documentação Detalhada
- [Arquitetura do Sistema](./ARQUITETURA.md)

## Como começar
1. Clone o repositório:
   ```bash
   git clone https://github.com/cristianorosa/eSimulate
   ```
2. Acesse a pasta do projeto:
   ```bash
   cd eSimulate
   ```
3. Estrutura inicial:
   - `backend/` e `frontend/` já criados para facilitar o desenvolvimento modular.

## Próximos Passos
- Definir frameworks e bibliotecas para backend e frontend.
- Iniciar a implementação do backend em Go, seguindo Clean Architecture.
- Implementar as primeiras funcionalidades: cadastro de questões, usuários e autenticação.

---

Este README será atualizado conforme o desenvolvimento avançar e novas informações forem recebidas.
