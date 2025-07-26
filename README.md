# eSimulate

Sistema completo para simulados de questões objetivas, voltado para aprendizado, concursos e provas de conhecimento.

## Visão Geral
- **Backend:** Go (Clean Architecture, PostgreSQL, JWT, testes automatizados, pipelines CI/CD)
- **Frontend:** Angular (Material, responsivo, feedback visual, testes unitários e e2e, pipelines CI/CD)
- **Documentação:** ARQUITETURA.md, FRONTEND.md, INSTRUCOES_DE_EXECUCAO.md
- **Automação:** Pipelines de CI/CD para frontend e backend, lint, cobertura de testes, build

## Estrutura do Projeto
- `backend/`: API REST em Go, Clean Architecture, integração com PostgreSQL
- `frontend/`: SPA Angular, Angular Material, integração com backend
- `.github/workflows/`: Pipelines de CI/CD para frontend e backend
- `ARQUITETURA.md`: Detalhes da arquitetura do sistema
- `FRONTEND.md`: Documentação do frontend Angular
- `INSTRUCOES_DE_EXECUCAO.md`: Passo a passo para rodar e configurar o sistema

## Principais Funcionalidades
- Cadastro e autenticação de usuários (tradicional e social)
- Cadastro, listagem e detalhamento de questões
- Criação, execução e resultado de simulados
- Dashboard com desempenho, histórico e gráficos
- Feedback visual (toasts, loaders), responsividade e acessibilidade
- Testes unitários e e2e automatizados
- Pipelines de CI/CD, lint, cobertura e build

## Como Rodar o Sistema
Consulte o arquivo [INSTRUCOES_DE_EXECUCAO.md](./INSTRUCOES_DE_EXECUCAO.md) para o passo a passo completo de instalação, configuração e execução do sistema.

## Como Contribuir
- Faça um fork do repositório
- Crie uma branch para sua feature/correção
- Garanta que os testes e lint passam localmente
- Envie um pull request detalhado
- Consulte os arquivos de documentação para padrões e arquitetura

## Documentação Complementar
- [ARQUITETURA.md](./ARQUITETURA.md): Arquitetura do sistema
- [FRONTEND.md](./frontend/FRONTEND.md): Detalhes do frontend Angular
- [INSTRUCOES_DE_EXECUCAO.md](./INSTRUCOES_DE_EXECUCAO.md): Passo a passo para rodar o sistema

---

O sistema está pronto para uso, manutenção e evolução contínua, seguindo as melhores práticas de desenvolvimento, automação e documentação.
