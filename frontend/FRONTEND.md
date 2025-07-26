# Documentação do Frontend eSimulate

## Visão Geral
O frontend do eSimulate foi desenvolvido em Angular (última versão), utilizando Angular Material para garantir um design moderno, responsivo e elegante, inspirado nas melhores práticas de UX/UI das grandes redes sociais.

## Estrutura de Pastas
- `src/app/core`: Serviços centrais (autenticação, guards, interceptors)
- `src/app/layout`: Layout principal (header, menu lateral, responsividade)
- `src/app/auth`: Login, cadastro e autenticação
- `src/app/dashboard`: Dashboard do usuário
- `src/app/simulados`: Listagem, criação e execução de simulados
- `src/app/questoes`: Listagem, cadastro, detalhamento e edição de questões
- `src/app/perfil`: Perfil do usuário (em expansão)

## Principais Fluxos
- **Autenticação JWT**: Login, cadastro, proteção de rotas e envio automático do token nas requisições.
- **Simulados**: Listagem, execução (questão a questão), envio de respostas e exibição de resultado detalhado.
- **Questões**: Cadastro, listagem, detalhamento e feedback visual.
- **Dashboard**: Resumo de desempenho, histórico e gráficos (em expansão).

## Padrões de Design
- Angular Material para componentes visuais, responsividade e acessibilidade.
- Feedback visual (toasts) para sucesso/erro em todas as ações críticas.
- Navegação fluida, microinterações e layout inspirado em apps modernos.
- SCSS para customização visual.

## Comandos Úteis
- Instalar dependências:
  ```bash
  npm install
  ```
- Rodar o frontend em modo desenvolvimento:
  ```bash
  ng serve --open
  ```
- Gerar novo componente/módulo/serviço:
  ```bash
  ng generate component nome
  ng generate service nome
  ng generate module nome
  ```

## Integração com Backend
- O frontend consome a API REST do backend Go, disponível em `http://localhost:8080`.
- O token JWT é enviado automaticamente nas requisições autenticadas.
- Endpoints e payloads seguem o padrão REST documentado no backend.

## Dicas de Desenvolvimento
- Utilize Angular Material para novos componentes e siga o padrão visual do sistema.
- Sempre forneça feedback visual ao usuário (sucesso, erro, carregando).
- Mantenha o código limpo, modular e documentado.
- Para dúvidas sobre arquitetura, consulte o arquivo `ARQUITETURA.md` na raiz do projeto.

---

Este documento será atualizado conforme o frontend evoluir.
