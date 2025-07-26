# Instruções para Rodar o Sistema eSimulate

## Pré-requisitos
- **Go** (versão 1.21 ou superior)
- **Node.js** (versão 18 ou superior) e **npm**
- **PostgreSQL** (versão 12 ou superior)
- **Git**

## 1. Clonar o repositório
```bash
git clone https://github.com/cristianorosa/eSimulate.git
cd eSimulate
```

## 2. Configurar o Banco de Dados
- Crie um banco de dados PostgreSQL chamado `esimulate`.
- Crie um usuário e senha (ou use o padrão do seu ambiente).
- Atualize as configurações de conexão no arquivo `backend/main.go` se necessário:
  - Host, porta, usuário, senha, nome do banco.
- Execute o script de criação das tabelas:
```bash
psql -U seu_usuario -d esimulate -f backend/db_schema.sql
```

## 3. Rodar o Backend
```bash
cd backend
go run main.go
```
- O backend estará disponível em `http://localhost:8080`.
- O backend já cria as tabelas e o usuário admin automaticamente se necessário.

## 4. Rodar o Frontend
```bash
cd frontend
npm install
ng serve --open
```
- O frontend estará disponível em `http://localhost:4200`.
- O frontend já está configurado para consumir a API do backend em `http://localhost:8080`.

## 5. Variáveis de Ambiente (opcional)
- Para produção, configure variáveis de ambiente para as conexões do banco e secrets do JWT.
- Ajuste o CORS no backend se for rodar frontend e backend em domínios diferentes.

## 6. Testes Automatizados
- **Backend:**
  ```bash
  cd backend
  go test ./...
  ```
- **Frontend:**
  ```bash
  cd frontend
  npm run test
  npx cypress open
  ```
- Os pipelines de CI/CD já executam todos os testes automaticamente a cada push/pull request.

## 7. Deploy (manual)
- Para produção, faça build do frontend (`ng build --configuration production`) e sirva os arquivos estáticos.
- O backend pode ser compilado com `go build` e rodado em qualquer servidor Linux/Windows/Mac.
- Configure variáveis de ambiente e secrets adequados para produção.
- Configure HTTPS e variáveis de CORS conforme necessário.

## 8. O que depende de você
- Instalar e configurar o PostgreSQL localmente.
- Garantir que as portas 8080 (backend) e 4200 (frontend) estejam livres.
- Ajustar variáveis de ambiente para produção (se necessário).
- Configurar domínio, HTTPS e deploy em ambiente de produção.
- Gerenciar secrets e tokens de autenticação social (Google/Facebook) se for usar login social.
- Manter backups do banco de dados.

## 9. O que já está pronto e automatizado
- Estrutura de código, testes, pipelines de CI/CD, lint, cobertura de testes, documentação, loaders, feedback visual, responsividade, acessibilidade básica.

---

Para dúvidas, consulte os arquivos README.md, ARQUITETURA.md e FRONTEND.md no projeto.
