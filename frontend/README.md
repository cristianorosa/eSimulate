# eSimulate Frontend

Interface web moderna para o sistema de simulados educacionais eSimulate, construída com Angular 17.

## 📋 Índice

- [Visão Geral](#visão-geral)
- [Tecnologias](#tecnologias)
- [Estrutura do Projeto](#estrutura-do-projeto)
- [Funcionalidades](#funcionalidades)
- [Instalação](#instalação)
- [Desenvolvimento](#desenvolvimento)
- [Componentes](#componentes)
- [Estilização](#estilização)

## 🎯 Visão Geral

O frontend do eSimulate é uma Single Page Application (SPA) que fornece interface intuitiva para:

- **Gestão Administrativa**: CRUD completo de áreas, exames, tópicos e questões
- **Relacionamentos N:N**: Interface para associar tópicos a exames e questões a exames
- **Sistema de Tags**: Categorização flexível de questões
- **Importação**: Upload e processamento de arquivos JSON
- **Filtros Avançados**: Busca e filtros em tempo real
- **Paginação**: Navegação eficiente em grandes listas

## 🚀 Tecnologias

### Core
- **Angular 17** - Framework principal
- **TypeScript 5.0+** - Linguagem de desenvolvimento
- **RxJS 7** - Programação reativa

### UI/UX
- **Angular Material 17** - Componentes de interface
- **SCSS** - Estilização avançada
- **Responsive Design** - Compatível com dispositivos móveis

### Ferramentas
- **Angular CLI** - Ferramenta de desenvolvimento
- **Webpack** - Bundler de módulos
- **ESLint** - Linting de código
- **Prettier** - Formatação de código

## 📁 Estrutura do Projeto

```
frontend/src/
├── app/
│   ├── admin/                    # Módulo administrativo
│   │   ├── areas/               # Gestão de áreas
│   │   ├── exams/               # Gestão de exames
│   │   ├── topics/              # Gestão de tópicos
│   │   └── questions/           # Gestão de questões
│   ├── core/                    # Serviços principais
│   │   ├── admin.service.ts     # Service para APIs admin
│   │   ├── auth.service.ts      # Autenticação
│   │   └── guards/              # Guards de rota
│   ├── shared/                  # Componentes compartilhados
│   │   ├── components/          # Componentes reutilizáveis
│   │   └── models/              # Interfaces TypeScript
│   ├── dashboard/               # Dashboard principal
│   └── auth/                    # Módulo de autenticação
├── assets/                      # Recursos estáticos
├── environments/                # Configurações de ambiente
└── styles.scss                 # Estilos globais
```

## ✨ Funcionalidades

### 🏢 Gestão de Áreas
- ✅ Listagem com paginação e busca
- ✅ Criação e edição via modal
- ✅ Exclusão com confirmação
- ✅ Filtro por nome em tempo real

### 📚 Gestão de Exames
- ✅ CRUD completo com formulários reativos
- ✅ Filtros por área e status (ativo/inativo)
- ✅ Configuração de tempo e nota mínima
- ✅ Associação com tópicos (N:N)
- ✅ Gestão de questões do exame

### 📖 Gestão de Tópicos
- ✅ Tópicos independentes (não vinculados a exame específico)
- ✅ Filtro por exame associado
- ✅ Descrição opcional
- ✅ Reutilização entre múltiplos exames

### ❓ Gestão de Questões
- ✅ Editor de questões com múltiplas opções
- ✅ Níveis de dificuldade (fácil/médio/difícil)
- ✅ Filtros por exame, tópico e dificuldade
- ✅ Sistema de tags para categorização
- ✅ **Importação via JSON** com upload de arquivo
- ✅ Status ativo/inativo com chips coloridos

### 🏷️ Sistema de Tags
- ✅ Tags simples (apenas nome)
- ✅ Associação N:N com questões
- ✅ Busca e filtro por tags
- ✅ Gestão via interface intuitiva

### 📥 Importação de Questões
- ✅ Upload de arquivo JSON via interface
- ✅ Validação de formato no frontend
- ✅ Processamento no backend
- ✅ Feedback de sucesso/erro
- ✅ Criação automática de estruturas

## 🛠️ Instalação

### Pré-requisitos
- **Node.js 18+**
- **npm 9+** ou **yarn 1.22+**
- **Angular CLI 17+**

### Passos

```bash
# Clonar repositório
git clone <repository-url>
cd eSimulate/frontend

# Instalar dependências
npm install

# Instalar Angular CLI globalmente (se necessário)
npm install -g @angular/cli@17

# Verificar instalação
ng version
```

## 🚀 Desenvolvimento

### Servidor de Desenvolvimento

```bash
# Iniciar servidor (http://localhost:4200)
ng serve

# Com proxy para backend
ng serve --proxy-config proxy.conf.json

# Modo de produção
ng serve --configuration production
```

### Scripts Disponíveis

```bash
# Desenvolvimento
npm start              # ng serve com proxy
npm run dev            # ng serve sem proxy

# Build
npm run build          # Build de produção
npm run build:dev      # Build de desenvolvimento

# Testes
npm test               # Testes unitários
npm run test:coverage  # Testes com coverage
npm run e2e            # Testes end-to-end

# Linting
npm run lint           # ESLint
npm run lint:fix       # ESLint com correção automática

# Formatação
npm run format         # Prettier
npm run format:check   # Verificar formatação
```

### Proxy Configuration

O arquivo `proxy.conf.json` redireciona chamadas da API para o backend:

```json
{
  "/api/*": {
    "target": "http://localhost:8080",
    "secure": false,
    "changeOrigin": true,
    "logLevel": "debug"
  }
}
```

## 🧩 Componentes Principais

### AdminService

Serviço central para comunicação com APIs administrativas:

```typescript
// Áreas
getAreas(): Observable<Area[]>
getAreasPaginated(page: number, pageSize: number, search?: string): Observable<PaginatedResponse<Area>>

// Exames  
getExams(): Observable<Exam[]>
getExamsPaginated(page: number, pageSize: number, areaId?: number, status?: string): Observable<PaginatedResponse<Exam>>

// Tópicos
getTopics(): Observable<Topic[]>
getTopicsPaginated(page: number, pageSize: number, examId?: number): Observable<PaginatedResponse<Topic>>

// Questões
getQuestions(): Observable<Question[]>
getQuestionsPaginated(page: number, pageSize: number, examId?: number, topicId?: number): Observable<PaginatedResponse<Question>>
importQuestions(data: any): Observable<any>
```

### Componentes de Gestão

#### AreasComponent
```typescript
// Funcionalidades principais
- loadAreas(): void                    // Carregar lista
- onSearchChange(term: string): void   // Busca em tempo real
- openDialog(area?: Area): void        // Modal criar/editar
- deleteArea(id: number): void         // Exclusão
- clearFilters(): void                 // Limpar filtros
```

#### ExamsComponent
```typescript
// Funcionalidades principais
- loadExams(): void                           // Carregar lista
- onAreaFilterChange(areaId: number): void    // Filtro por área
- onStatusFilterChange(status: string): void  // Filtro por status
- openDialog(exam?: Exam): void               // Modal criar/editar
- saveExam(): void                            // Salvar (create/update)
```

#### TopicsComponent
```typescript
// Funcionalidades principais
- loadTopics(): void                          // Carregar lista paginada
- onPageChange(event: PageEvent): void        // Mudança de página
- onExamFilterChange(examId: number): void    // Filtro por exame
- openDialog(topic?: Topic): void             // Modal criar/editar
```

#### QuestionsComponent
```typescript
// Funcionalidades principais
- loadQuestions(): void                                    // Carregar lista paginada
- onExamFilterChange(examId: number): void                 // Filtro por exame
- onTopicFilterChange(topicId: number): void               // Filtro por tópico
- importQuestions(): void                                  // Abrir seletor de arquivo
- onFileSelected(event: any): void                         // Processar arquivo JSON
- processImportData(data: any): void                       // Enviar para backend
```

## 🎨 Estilização

### Arquitetura CSS

O projeto utiliza uma arquitetura CSS modular e escalável:

```scss
// styles.scss - Estilos globais
$primary-color: #1976d2;
$secondary-color: #424242;
$success-color: #4caf50;
$warning-color: #ff9800;
$error-color: #f44336;

// Mixins globais
@mixin card-style {
  background: white;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
  padding: 24px;
}

@mixin responsive-grid {
  display: grid;
  gap: 16px;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
}
```

### Componentes de Estilo

#### Layout Padrão das Telas Admin
```scss
.admin-container {
  padding: 24px;
  max-width: 1200px;
  margin: 0 auto;

  .header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 24px;
    padding-bottom: 16px;
    border-bottom: 1px solid #e0e0e0;
  }

  .content {
    .filters-section {
      margin-bottom: 24px;
      padding: 16px;
      background: #f5f5f5;
      border-radius: 8px;
    }
  }
}
```

#### Sistema de Chips para Status
```scss
.status-active {
  background-color: #e8f5e8 !important;
  color: #2e7d32 !important;
}

.status-inactive {
  background-color: #ffebee !important;
  color: #c62828 !important;
}
```

#### Responsividade
```scss
@media (max-width: 768px) {
  .admin-container {
    padding: 16px;
    
    .header {
      flex-direction: column;
      gap: 16px;
    }
    
    .filters-section {
      .form-row {
        flex-direction: column;
      }
    }
  }
}
```

### Material Design Customization

```scss
// Customização de componentes Material
.mat-mdc-form-field {
  width: 100%;
  
  .mat-mdc-form-field-subscript-wrapper {
    height: 0 !important; // Remove espaço extra
  }
}

.mat-mdc-dialog-container {
  .mat-mdc-dialog-surface {
    border-radius: 12px;
    padding: 0;
  }
}

.mat-mdc-table {
  .mat-mdc-header-cell {
    font-weight: 600;
    color: #424242;
  }
  
  .mat-mdc-cell {
    border-bottom: 1px solid #e0e0e0;
  }
}
```

## 📱 Funcionalidades por Tela

### Dashboard
- ✅ Visão geral do sistema
- ✅ Navegação rápida para seções
- ✅ Estatísticas básicas
- ⏳ Gráficos de desempenho

### Áreas
- ✅ Grid responsiva com paginação
- ✅ Busca em tempo real
- ✅ Modal de criação/edição
- ✅ Confirmação de exclusão

### Exames
- ✅ Filtros por área e status
- ✅ Formulário reativo completo
- ✅ Validações de tempo e nota
- ✅ Chips coloridos para status

### Tópicos
- ✅ Filtro por exame associado
- ✅ Layout compacto otimizado
- ✅ Header fixo fora do scroll
- ✅ Paginação integrada

### Questões
- ✅ Filtros múltiplos (exame, tópico, dificuldade)
- ✅ Editor de opções dinâmico
- ✅ **Importação de JSON** com drag-and-drop
- ✅ Sistema de tags integrado
- ✅ Chips coloridos para status e dificuldade

## 🔧 Configuração de Ambiente

### Development
```typescript
// environment.ts
export const environment = {
  production: false,
  apiUrl: 'http://localhost:8080/api',
  enableDebugTools: true
};
```

### Production
```typescript
// environment.prod.ts
export const environment = {
  production: true,
  apiUrl: '/api',
  enableDebugTools: false
};
```

## 🧪 Testes

### Estrutura de Testes
```
src/
├── app/
│   ├── admin/
│   │   ├── areas/
│   │   │   ├── areas.component.spec.ts
│   │   │   └── areas.component.ts
│   │   └── ...
│   └── core/
│       ├── admin.service.spec.ts
│       └── admin.service.ts
```

### Executar Testes

```bash
# Testes unitários
ng test

# Testes com coverage
ng test --code-coverage

# Testes em modo watch
ng test --watch

# Testes end-to-end
ng e2e
```

## 📚 Padrões de Desenvolvimento

### Convenções de Nomenclatura
- **Componentes**: PascalCase (`AreasComponent`)
- **Serviços**: PascalCase com sufixo Service (`AdminService`)
- **Interfaces**: PascalCase (`Area`, `Exam`)
- **Métodos**: camelCase (`loadAreas`, `onPageChange`)
- **Propriedades**: camelCase (`currentPage`, `selectedAreaId`)

### Estrutura de Componente
```typescript
@Component({
  selector: 'app-areas',
  templateUrl: './areas.component.html',
  styleUrls: ['./areas.component.scss']
})
export class AreasComponent implements OnInit, OnDestroy {
  // Propriedades públicas
  areas: Area[] = [];
  loading = false;
  
  // Propriedades de paginação
  currentPage = 1;
  currentPageSize = 10;
  totalItems = 0;
  
  // Propriedades de filtro
  searchTerm = '';
  
  // Propriedades privadas
  private destroy$ = new Subject<void>();
  
  constructor(
    private adminService: AdminService,
    private snackBar: MatSnackBar
  ) {}
  
  ngOnInit(): void {
    this.loadAreas();
  }
  
  ngOnDestroy(): void {
    this.destroy$.next();
    this.destroy$.complete();
  }
  
  // Métodos públicos
  loadAreas(): void { /* ... */ }
  onSearchChange(term: string): void { /* ... */ }
  
  // Métodos privados
  private showSuccessMessage(message: string): void { /* ... */ }
}
```

## 🚀 Deploy

### Build de Produção
```bash
# Build otimizado
ng build --configuration production

# Arquivos gerados em dist/
ls dist/esimulate-frontend/
```

### Configuração do Servidor
```nginx
# nginx.conf
server {
    listen 80;
    server_name localhost;
    root /usr/share/nginx/html;
    
    # Angular routing
    location / {
        try_files $uri $uri/ /index.html;
    }
    
    # API proxy
    location /api/ {
        proxy_pass http://backend:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

## 🤝 Contribuição

### Workflow
1. Fork do repositório
2. Criar branch para feature (`git checkout -b feature/nova-funcionalidade`)
3. Implementar mudanças seguindo padrões
4. Executar testes (`npm test`)
5. Commit com mensagem descritiva
6. Push e criação de Pull Request

### Padrões de Commit
```
feat: adicionar importação de questões
fix: corrigir paginação em tópicos
refactor: simplificar sistema de tags
docs: atualizar README
style: ajustar espaçamento em cards
test: adicionar testes para AdminService
```

## 📄 Licença

Este projeto está sob a licença MIT. Veja o arquivo `LICENSE` para mais detalhes.