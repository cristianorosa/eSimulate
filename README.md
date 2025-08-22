# 🎓 eSimulate - Sistema de Simulados Educacionais

Sistema completo para criação, gestão e aplicação de simulados educacionais com arquitetura moderna e relacionamentos flexíveis.

![Version](https://img.shields.io/badge/version-2.0.0-blue.svg)
![Go](https://img.shields.io/badge/Go-1.21+-00ADD8.svg)
![Angular](https://img.shields.io/badge/Angular-17-DD0031.svg)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-14+-336791.svg)
![License](https://img.shields.io/badge/license-MIT-green.svg)

## 📋 Índice

- [Visão Geral](#visão-geral)
- [Características](#características)
- [Arquitetura](#arquitetura)
- [Tecnologias](#tecnologias)
- [Instalação Rápida](#instalação-rápida)
- [Documentação](#documentação)
- [Screenshots](#screenshots)
- [Roadmap](#roadmap)
- [Contribuição](#contribuição)

## 🎯 Visão Geral

O **eSimulate** é uma plataforma completa para criação e aplicação de simulados educacionais, projetada para instituições de ensino, professores e estudantes. 

### Principais Funcionalidades

🏢 **Gestão Hierárquica**
- Áreas de conhecimento
- Exames por área
- Tópicos independentes e reutilizáveis
- Questões com múltiplas opções

🔗 **Relacionamentos Flexíveis**
- Exames ↔ Tópicos (N:N) com configurações específicas
- Questões ↔ Exames (N:N) com validação inteligente
- Questões ↔ Tags (N:N) para categorização

📥 **Importação Inteligente**
- Upload de questões via JSON
- Criação automática de estruturas
- Validação e rollback automático

🎨 **Interface Moderna**
- Design responsivo com Angular Material
- Filtros avançados e paginação
- Experiência de usuário intuitiva

## ✨ Características

### 🚀 Tecnologia de Ponta
- **Backend**: Go com Clean Architecture
- **Frontend**: Angular 17 com TypeScript
- **Banco**: PostgreSQL com relacionamentos otimizados
- **API**: REST completa e documentada

### 🔧 Flexibilidade Total
- **Tópicos Reutilizáveis**: Um tópico pode ser usado em múltiplos exames
- **Distribuição Inteligente**: Configure peso e dificuldade por tópico
- **Tags Flexíveis**: Categorize questões com sistema de tags simples
- **Validação Hierárquica**: Regras de negócio implementadas na aplicação

### 📊 Gestão Avançada
- **Paginação Universal**: Todas as listas são paginadas
- **Filtros Específicos**: Busque por área, status, dificuldade, tags
- **Importação em Lote**: Importe provas completas via JSON
- **Validação Inteligente**: Questões só podem ser associadas a exames se o tópico permitir

## 🏗️ Arquitetura

### Estrutura Geral
```
eSimulate/
├── backend/          # API Go com Clean Architecture
├── frontend/         # SPA Angular 17
├── docs/            # Documentação adicional
└── scripts/         # Scripts de deploy e utilitários
```

### Modelo de Dados
```
Área (1) ──→ (N) Exame
           ↗
Tópico (N) ←──→ (N) Exame (com peso, dificuldade, quantidade)
  │
  ↓ (1:N)
Questão (N) ←──→ (N) Exame (com validação hierárquica)
  │
  ↓ (N:N)
Tag (categorização flexível)
```

### Clean Architecture
```
┌─────────────────┐    ┌─────────────────┐
│   Frontend      │    │   Backend       │
│   Angular 17    │◄──►│   Go API        │
└─────────────────┘    └─────────────────┘
                              │
                    ┌─────────────────┐
                    │  PostgreSQL     │
                    │  Database       │
                    └─────────────────┘
```

## 🛠️ Tecnologias

### Backend
- **Go 1.21+** - Performance e concorrência
- **PostgreSQL 14+** - Banco relacional robusto
- **JWT** - Autenticação segura
- **Clean Architecture** - Código maintível e testável

### Frontend
- **Angular 17** - Framework moderno
- **TypeScript 5** - Tipagem estática
- **Angular Material** - Componentes de UI
- **RxJS** - Programação reativa

### DevOps
- **Docker** - Containerização
- **Git** - Controle de versão
- **GitHub Actions** - CI/CD (planejado)

## 🚀 Instalação Rápida

### Pré-requisitos
- **Go 1.21+**
- **Node.js 18+**
- **PostgreSQL 14+**
- **Git**

### 1. Clonar Repositório
```bash
git clone <repository-url>
cd eSimulate
```

### 2. Configurar Banco de Dados
```bash
# Criar banco
createdb esimulate

# Executar schema
psql -d esimulate -f backend/db_schema_v2.sql
```

### 3. Configurar Backend
```bash
cd backend

# Instalar dependências
go mod tidy

# Configurar variáveis de ambiente
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=postgres
export DB_PASSWORD=password
export DB_NAME=esimulate
export JWT_SECRET=your-secret-key

# Executar
go run main.go
```

### 4. Configurar Frontend
```bash
cd frontend

# Instalar dependências
npm install

# Executar com proxy
npm start
```

### 5. Acessar Sistema
- **Frontend**: http://localhost:4200
- **Backend**: http://localhost:8080
- **API Docs**: http://localhost:8080/api (em breve)

## 📚 Documentação

### Documentação Técnica
- [**Backend README**](backend/README.md) - Arquitetura e APIs
- [**Frontend README**](frontend/README.md) - Interface e componentes
- [**API Documentation**](backend/API_DOCUMENTATION.md) - Endpoints completos
- [**Requirements**](backend/REQUIREMENTS.md) - Requisitos funcionais
- [**Changelog**](backend/CHANGELOG.md) - Histórico de versões

### Guias de Uso
- [**Guia de Instalação**](docs/INSTALLATION.md) - Instalação detalhada (em breve)
- [**Guia do Usuário**](docs/USER_GUIDE.md) - Como usar o sistema (em breve)
- [**Guia de Desenvolvimento**](docs/DEVELOPMENT.md) - Para desenvolvedores (em breve)

## 🖼️ Screenshots

### Dashboard Principal
```
┌─────────────────────────────────────────────┐
│  🎓 eSimulate                    👤 Admin   │
├─────────────────────────────────────────────┤
│                                             │
│  📊 Visão Geral                            │
│  ┌─────────┐ ┌─────────┐ ┌─────────┐       │
│  │ 15 Áreas│ │ 45 Exam.│ │ 320 Quest│      │
│  └─────────┘ └─────────┘ └─────────┘       │
│                                             │
│  🚀 Ações Rápidas                          │
│  [Criar Área] [Novo Exame] [Import JSON]   │
│                                             │
└─────────────────────────────────────────────┘
```

### Gestão de Questões
```
┌─────────────────────────────────────────────┐
│  ❓ Questões                  [+ Nova] [📥]  │
├─────────────────────────────────────────────┤
│  🔍 Filtros                                 │
│  [Exame: JavaScript ▼] [Tópico: Básico ▼]  │
│  [Dificuldade: Todas ▼] [Tags: react,js]   │
├─────────────────────────────────────────────┤
│  ID │ Questão           │ Tópico │ Diff │ Tags│
│  001 │ O que é JS?      │ Básico │ 🟢   │ js  │
│  002 │ Como criar var?  │ Básico │ 🟡   │ var │
│  003 │ Async/await?     │ Avançd │ 🔴   │async│
├─────────────────────────────────────────────┤
│                    ← 1 2 3 ... 15 →         │
└─────────────────────────────────────────────┘
```

## 🗺️ Roadmap

### ✅ Versão 2.0 (Atual)
- [x] Sistema completo de relacionamentos N:N
- [x] Interface administrativa completa
- [x] Sistema de tags flexível
- [x] Importação de questões via JSON
- [x] Validação de regras de negócio

### 🚧 Versão 2.1 (Em Desenvolvimento)
- [ ] Aplicação de simulados online
- [ ] Correção automática
- [ ] Relatórios básicos de desempenho
- [ ] Sistema de usuários estudantes

### 🔮 Versão 3.0 (Planejado)
- [ ] Dashboard com gráficos avançados
- [ ] Análise de desempenho por tópico
- [ ] Integração com sistemas LMS
- [ ] API pública para terceiros
- [ ] Mobile app (React Native)

### 🌟 Versão 4.0 (Futuro)
- [ ] IA para geração automática de questões
- [ ] Banco de questões público
- [ ] Sistema de gamificação
- [ ] Análise preditiva de desempenho
- [ ] Integração com videoconferência

## 🤝 Contribuição

Contribuições são muito bem-vindas! Este projeto segue as melhores práticas de desenvolvimento colaborativo.

### Como Contribuir

1. **Fork** o repositório
2. **Clone** seu fork localmente
3. **Crie** uma branch para sua feature (`git checkout -b feature/nova-funcionalidade`)
4. **Implemente** suas mudanças seguindo os padrões do projeto
5. **Teste** suas mudanças
6. **Commit** com mensagens descritivas (`git commit -m 'feat: adicionar funcionalidade X'`)
7. **Push** para sua branch (`git push origin feature/nova-funcionalidade`)
8. **Abra** um Pull Request

### Padrões de Commit
```
feat: nova funcionalidade
fix: correção de bug
docs: atualização de documentação
style: mudanças de formatação
refactor: refatoração de código
test: adição/modificação de testes
chore: tarefas de manutenção
```

### Áreas de Contribuição
- 🐛 **Bug Reports** - Encontrou um problema? Reporte!
- 💡 **Feature Requests** - Tem uma ideia? Compartilhe!
- 📝 **Documentação** - Ajude a melhorar a documentação
- 🧪 **Testes** - Mais testes são sempre bem-vindos
- 🎨 **UI/UX** - Melhorias na interface
- 🔧 **Performance** - Otimizações são valiosas

## 📞 Suporte

### Comunidade
- **Issues**: Reporte bugs e solicite features
- **Discussions**: Tire dúvidas e compartilhe ideias
- **Wiki**: Documentação colaborativa (em breve)

### Documentação
- **README Técnico**: Documentação detalhada por módulo
- **API Docs**: Especificação completa das APIs
- **Tutoriais**: Guias passo-a-passo (em desenvolvimento)

## 📄 Licença

Este projeto está licenciado sob a **MIT License** - veja o arquivo [LICENSE](LICENSE) para detalhes.

### Resumo da Licença
- ✅ **Uso comercial** permitido
- ✅ **Modificação** permitida
- ✅ **Distribuição** permitida
- ✅ **Uso privado** permitido
- ❌ **Garantia** não fornecida
- ❌ **Responsabilidade** do autor limitada

---

## 🌟 Agradecimentos

Agradecimentos especiais a todos os contribuidores e à comunidade open source pelas ferramentas e bibliotecas que tornam este projeto possível.

### Tecnologias Utilizadas
- [Go](https://golang.org/) - Linguagem de programação
- [Angular](https://angular.io/) - Framework frontend
- [PostgreSQL](https://www.postgresql.org/) - Banco de dados
- [Angular Material](https://material.angular.io/) - Componentes UI
- [RxJS](https://rxjs.dev/) - Programação reativa

---

<div align="center">

**🎓 eSimulate - Transformando a educação através da tecnologia**

[![Feito com ❤️](https://img.shields.io/badge/Feito%20com-❤️-red.svg)](https://github.com/your-username/eSimulate)
[![Contribuidores](https://img.shields.io/badge/Contribuidores-Bem%20vindos-brightgreen.svg)](CONTRIBUTING.md)

[⬆️ Voltar ao topo](#-esimulate---sistema-de-simulados-educacionais)

</div>