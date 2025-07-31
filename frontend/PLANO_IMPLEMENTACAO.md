# 📋 Plano de Implementação - Frontend eSimulate v2

## 🎯 **Objetivo**
Implementar as novas funcionalidades do eSimulate v2 no frontend Angular, incluindo níveis de acesso, gerenciamento de provas, aplicação de exames e análise de desempenho.

## 🏗️ **Estrutura de Componentes**

### 🔐 **Autenticação e Autorização**
```
auth/
├── login/
├── register/
└── role-guard.ts (novo)
```

### 📚 **Gerenciamento de Conteúdo (Redator/Admin)**
```
admin/
├── areas/
│   ├── area-list/
│   ├── area-form/
│   └── area-detail/
├── exams/
│   ├── exam-list/
│   ├── exam-form/
│   ├── exam-detail/
│   └── exam-domains/
├── questions/
│   ├── question-list/
│   ├── question-form/
│   └── question-detail/
└── users/
    ├── user-list/
    ├── user-form/
    └── user-detail/
```

### 🎯 **Aplicação de Provas (Usuário)**
```
exams/
├── exam-list/
├── exam-detail/
├── exam-execution/
│   ├── question-viewer/
│   ├── answer-form/
│   └── review-panel/
└── exam-result/
    ├── score-summary/
    ├── domain-analysis/
    └── improvement-suggestions/
```

### 📊 **Análise e Relatórios**
```
analytics/
├── dashboard/
├── performance/
├── history/
└── reports/
```

## 🚀 **Fases de Implementação**

### 📋 **Fase 1: Estrutura Base**
- [ ] Criar estrutura de pastas
- [ ] Implementar role-guard para níveis de acesso
- [ ] Atualizar AuthService para incluir role
- [ ] Criar interfaces TypeScript para novas entidades

### 📋 **Fase 2: Gerenciamento de Áreas**
- [ ] Componente area-list
- [ ] Componente area-form
- [ ] Componente area-detail
- [ ] Service para áreas
- [ ] Rotas protegidas para redator/admin

### 📋 **Fase 3: Gerenciamento de Exames**
- [ ] Componente exam-list
- [ ] Componente exam-form
- [ ] Componente exam-detail
- [ ] Componente exam-domains
- [ ] Service para exames
- [ ] Upload de questões em lote

### 📋 **Fase 4: Gerenciamento de Questões**
- [ ] Componente question-list
- [ ] Componente question-form
- [ ] Componente question-detail
- [ ] Editor de alternativas
- [ ] Service para questões

### 📋 **Fase 5: Aplicação de Provas**
- [ ] Componente exam-execution
- [ ] Componente question-viewer
- [ ] Componente answer-form
- [ ] Componente review-panel
- [ ] Timer de prova
- [ ] Salvamento automático

### 📋 **Fase 6: Resultados e Análise**
- [ ] Componente exam-result
- [ ] Componente score-summary
- [ ] Componente domain-analysis
- [ ] Componente improvement-suggestions
- [ ] Gráficos de desempenho

### 📋 **Fase 7: Relatórios e Dashboard**
- [ ] Dashboard administrativo
- [ ] Relatórios de performance
- [ ] Histórico de tentativas
- [ ] Análise comparativa

## 🔧 **Interfaces TypeScript**

### 👥 **User Interface**
```typescript
interface User {
  id: number;
  name: string;
  email: string;
  role: 'user' | 'redator' | 'admin';
  createdAt: Date;
  updatedAt: Date;
}
```

### 📚 **Area Interface**
```typescript
interface Area {
  id: number;
  name: string;
  description: string;
  createdAt: Date;
}
```

### 📝 **Exam Interface**
```typescript
interface Exam {
  id: number;
  title: string;
  description: string;
  areaId: number;
  maxTimeMinutes: number;
  passingScore: number;
  isActive: boolean;
  createdBy: number;
  createdAt: Date;
  updatedAt: Date;
  domains?: Domain[];
  questions?: Question[];
}
```

### 🧩 **Domain Interface**
```typescript
interface Domain {
  id: number;
  examId: number;
  name: string;
  description: string;
  weightPercentage: number;
  orderIndex: number;
  createdAt: Date;
}
```

### ❓ **Question Interface**
```typescript
interface Question {
  id: number;
  examId: number;
  domainId: number;
  statement: string;
  explanation: string;
  difficultyLevel: number;
  createdBy: number;
  isActive: boolean;
  createdAt: Date;
  updatedAt: Date;
  options?: Option[];
}
```

### ✅ **Option Interface**
```typescript
interface Option {
  id: number;
  questionId: number;
  text: string;
  isCorrect: boolean;
  explanation: string;
  orderIndex: number;
}
```

### 🎯 **UserExam Interface**
```typescript
interface UserExam {
  id: number;
  userId: number;
  examId: number;
  startedAt: Date;
  finishedAt?: Date;
  totalScore?: number;
  passed?: boolean;
  timeSpentMinutes?: number;
  answers?: UserAnswer[];
  domainPerformance?: DomainPerformance[];
}
```

### 📝 **UserAnswer Interface**
```typescript
interface UserAnswer {
  id: number;
  userExamId: number;
  questionId: number;
  optionId?: number;
  isCorrect?: boolean;
  isMarkedForReview: boolean;
  answeredAt: Date;
}
```

### 📊 **DomainPerformance Interface**
```typescript
interface DomainPerformance {
  id: number;
  userExamId: number;
  domainId: number;
  questionsAnswered: number;
  correctAnswers: number;
  scorePercentage: number;
  needsImprovement: boolean;
}
```

## 🔐 **Role Guard**

### 🛡️ **RoleGuard Implementation**
```typescript
@Injectable({
  providedIn: 'root'
})
export class RoleGuard implements CanActivate {
  constructor(
    private authService: AuthService,
    private router: Router
  ) {}

  canActivate(
    route: ActivatedRouteSnapshot,
    state: RouterStateSnapshot
  ): boolean {
    const requiredRoles = route.data['roles'] as string[];
    const user = this.authService.getCurrentUser();
    
    if (!user) {
      this.router.navigate(['/login']);
      return false;
    }

    if (requiredRoles && !requiredRoles.includes(user.role)) {
      this.router.navigate(['/dashboard']);
      return false;
    }

    return true;
  }
}
```

## 📱 **Componentes Principais**

### 🎯 **Exam Execution Component**
```typescript
@Component({
  selector: 'app-exam-execution',
  template: `
    <div class="exam-container">
      <!-- Timer -->
      <div class="timer">
        <mat-icon>schedule</mat-icon>
        <span>{{ timeRemaining }}</span>
      </div>

      <!-- Question Viewer -->
      <app-question-viewer
        [question]="currentQuestion"
        [userAnswer]="currentAnswer"
        (answerSubmitted)="onAnswerSubmitted($event)"
        (markedForReview)="onMarkedForReview($event)">
      </app-question-viewer>

      <!-- Navigation -->
      <div class="navigation">
        <button mat-button (click)="previousQuestion()">
          <mat-icon>chevron_left</mat-icon>
          Anterior
        </button>
        <button mat-button (click)="nextQuestion()">
          Próxima
          <mat-icon>chevron_right</mat-icon>
        </button>
      </div>

      <!-- Review Panel -->
      <app-review-panel
        [questions]="questions"
        [userAnswers]="userAnswers"
        (questionSelected)="goToQuestion($event)">
      </app-review-panel>

      <!-- Finish Button -->
      <button mat-raised-button color="primary" (click)="finishExam()">
        Finalizar Prova
      </button>
    </div>
  `
})
export class ExamExecutionComponent {
  @Input() examId: number;
  @Input() userExamId: number;

  currentQuestion: Question;
  currentAnswer: UserAnswer;
  questions: Question[] = [];
  userAnswers: UserAnswer[] = [];
  timeRemaining: string;

  constructor(
    private examService: ExamService,
    private timerService: TimerService
  ) {}

  ngOnInit() {
    this.loadExam();
    this.startTimer();
  }

  onAnswerSubmitted(answer: UserAnswer) {
    this.examService.submitAnswer(this.userExamId, answer);
  }

  onMarkedForReview(questionId: number) {
    this.examService.markForReview(this.userExamId, questionId);
  }

  finishExam() {
    this.examService.finishExam(this.userExamId).subscribe(result => {
      this.router.navigate(['/exam-result', result.id]);
    });
  }
}
```

### 📊 **Exam Result Component**
```typescript
@Component({
  selector: 'app-exam-result',
  template: `
    <div class="result-container">
      <!-- Score Summary -->
      <app-score-summary
        [userExam]="userExam"
        [exam]="exam">
      </app-score-summary>

      <!-- Domain Analysis -->
      <app-domain-analysis
        [domainPerformance]="domainPerformance">
      </app-domain-analysis>

      <!-- Improvement Suggestions -->
      <app-improvement-suggestions
        [domainPerformance]="domainPerformance">
      </app-improvement-suggestions>
    </div>
  `
})
export class ExamResultComponent {
  @Input() userExamId: number;

  userExam: UserExam;
  exam: Exam;
  domainPerformance: DomainPerformance[];

  constructor(
    private examService: ExamService,
    private performanceService: PerformanceService
  ) {}

  ngOnInit() {
    this.loadExamResult();
  }

  loadExamResult() {
    this.examService.getUserExam(this.userExamId).subscribe(userExam => {
      this.userExam = userExam;
      this.loadExamDetails();
      this.loadDomainPerformance();
    });
  }
}
```

## 🎨 **Design System**

### 🎨 **Cores e Temas**
```scss
// Variáveis de cor
$primary-color: #1976d2;
$secondary-color: #dc004e;
$success-color: #4caf50;
$warning-color: #ff9800;
$error-color: #f44336;

// Temas por nível de acesso
.user-theme {
  --primary-color: #1976d2;
  --accent-color: #ff4081;
}

.redator-theme {
  --primary-color: #ff9800;
  --accent-color: #ff5722;
}

.admin-theme {
  --primary-color: #9c27b0;
  --accent-color: #e91e63;
}
```

### 📱 **Componentes Reutilizáveis**
- `app-card` - Cards padronizados
- `app-button` - Botões com variantes
- `app-progress` - Barras de progresso
- `app-timer` - Cronômetro de prova
- `app-question-card` - Card de questão
- `app-answer-option` - Opção de resposta

## 📊 **Serviços**

### 🔧 **ExamService**
```typescript
@Injectable({
  providedIn: 'root'
})
export class ExamService {
  constructor(private http: HttpClient) {}

  // Gerenciamento de exames
  getExams(areaId?: number): Observable<Exam[]> {
    const params = areaId ? { area_id: areaId.toString() } : {};
    return this.http.get<Exam[]>('/api/exams', { params });
  }

  getActiveExams(): Observable<Exam[]> {
    return this.http.get<Exam[]>('/api/exams/active');
  }

  createExam(exam: Exam): Observable<Exam> {
    return this.http.post<Exam>('/api/exams/create', exam);
  }

  // Aplicação de exames
  startExam(userId: number, examId: number): Observable<UserExam> {
    return this.http.post<UserExam>('/api/user-exams/start', { userId, examId });
  }

  submitAnswer(userExamId: number, answer: UserAnswer): Observable<void> {
    return this.http.post<void>('/api/user-exams/submit-answer', answer);
  }

  finishExam(userExamId: number): Observable<UserExam> {
    return this.http.post<UserExam>(`/api/user-exams/finish?id=${userExamId}`, {});
  }

  getUserExam(id: number): Observable<UserExam> {
    return this.http.get<UserExam>(`/api/user-exams/detail?id=${id}`);
  }
}
```

### 📊 **PerformanceService**
```typescript
@Injectable({
  providedIn: 'root'
})
export class PerformanceService {
  constructor(private http: HttpClient) {}

  getDomainPerformance(userExamId: number): Observable<DomainPerformance[]> {
    return this.http.get<DomainPerformance[]>(`/api/performance/domain?user_exam_id=${userExamId}`);
  }

  getUserHistory(userId: number): Observable<UserExam[]> {
    return this.http.get<UserExam[]>(`/api/user-exams/list?user_id=${userId}`);
  }

  getPerformanceReport(userId: number): Observable<any> {
    return this.http.get<any>(`/api/performance?user_id=${userId}`);
  }
}
```

## 🚀 **Cronograma de Implementação**

### 📅 **Semana 1-2: Estrutura Base**
- [ ] Role guard e autorização
- [ ] Interfaces TypeScript
- [ ] Serviços base
- [ ] Componentes de área

### 📅 **Semana 3-4: Gerenciamento de Exames**
- [ ] CRUD de exames
- [ ] Gerenciamento de domínios
- [ ] Upload de questões
- [ ] Validações

### 📅 **Semana 5-6: Aplicação de Provas**
- [ ] Interface de execução
- [ ] Timer e salvamento
- [ ] Revisão de questões
- [ ] Finalização

### 📅 **Semana 7-8: Resultados e Análise**
- [ ] Tela de resultados
- [ ] Análise por domínios
- [ ] Sugestões de melhoria
- [ ] Gráficos

### 📅 **Semana 9-10: Polimento**
- [ ] Testes e correções
- [ ] Otimizações
- [ ] Documentação
- [ ] Deploy

## 🧪 **Testes**

### 🔍 **Testes Unitários**
- [ ] Serviços
- [ ] Guards
- [ ] Componentes
- [ ] Utilitários

### 🎯 **Testes E2E**
- [ ] Fluxo completo de prova
- [ ] Gerenciamento de conteúdo
- [ ] Níveis de acesso
- [ ] Performance

## 📋 **Checklist de Qualidade**

### ✅ **Funcionalidades**
- [ ] Níveis de acesso funcionando
- [ ] CRUD de áreas
- [ ] CRUD de exames
- [ ] Aplicação de provas
- [ ] Análise de resultados
- [ ] Timer de prova
- [ ] Revisão de questões

### ✅ **UX/UI**
- [ ] Design responsivo
- [ ] Acessibilidade
- [ ] Performance
- [ ] Feedback visual
- [ ] Loading states

### ✅ **Técnico**
- [ ] Código limpo
- [ ] Testes cobertos
- [ ] Documentação
- [ ] Performance otimizada
- [ ] Segurança

---

**Status:** Em desenvolvimento  
**Versão:** 2.0  
**Última atualização:** 2024 