import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { MatCardModule } from '@angular/material/card';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatTableModule } from '@angular/material/table';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatSelectModule } from '@angular/material/select';
import { MatSnackBarModule } from '@angular/material/snack-bar';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatDialogModule } from '@angular/material/dialog';
import { MatChipsModule } from '@angular/material/chips';
import { MatCheckboxModule } from '@angular/material/checkbox';
import { MatTooltipModule } from '@angular/material/tooltip';
import { FormsModule, ReactiveFormsModule, FormBuilder, FormGroup, Validators, FormArray } from '@angular/forms';
import { AdminService, Question, Exam, Domain, Option } from '../../core/admin.service';
import { MatSnackBar } from '@angular/material/snack-bar';

interface QuestionWithDetails extends Question {
  exam_name?: string;
  domain_name?: string;
  options_count?: number;
}

@Component({
  selector: 'app-questions',
  standalone: true,
  imports: [
    CommonModule, 
    MatCardModule, 
    MatButtonModule, 
    MatIconModule,
    MatTableModule,
    MatFormFieldModule,
    MatInputModule,
    MatSelectModule,
    MatSnackBarModule,
    MatProgressSpinnerModule,
    MatDialogModule,
    MatChipsModule,
    MatCheckboxModule,
    MatTooltipModule,
    FormsModule,
    ReactiveFormsModule
  ],
  templateUrl: './questions.component.html',
  styleUrls: ['./questions.component.scss']
})
export class QuestionsComponent implements OnInit {
  questions: QuestionWithDetails[] = [];
  exams: Exam[] = [];
  domains: Domain[] = [];
  loading = false;
  showDialog = false;
  editingQuestion: QuestionWithDetails | null = null;
  questionForm: FormGroup;

  displayedColumns: string[] = ['id', 'statement', 'exam_name', 'domain_name', 'difficulty_level', 'is_active', 'options_count', 'actions'];

  constructor(
    private adminService: AdminService,
    private snackBar: MatSnackBar,
    private fb: FormBuilder
  ) {
    this.questionForm = this.fb.group({
      exam_id: ['', Validators.required],
      domain_id: ['', Validators.required],
      statement: ['', [Validators.required, Validators.minLength(10)]],
      explanation: ['', [Validators.required, Validators.minLength(10)]],
      difficulty_level: [1, [Validators.required, Validators.min(1), Validators.max(5)]],
      is_active: [true],
      options: this.fb.array([])
    });
  }

  ngOnInit() {
    this.loadExams();
    this.loadQuestions();
  }

  loadQuestions() {
    this.loading = true;
    this.adminService.getQuestions().subscribe({
      next: (questions) => {
        this.questions = questions.map(question => ({
          ...question,
          exam_name: this.getExamName(question.exam_id),
          domain_name: this.getDomainName(question.domain_id),
          options_count: question.options?.length || 0
        }));
        this.loading = false;
      },
      error: (error) => {
        console.error('Erro ao carregar questões:', error);
        this.snackBar.open('Erro ao carregar questões', 'Fechar', { duration: 3000 });
        this.loading = false;
      }
    });
  }

  loadExams() {
    this.adminService.getExams().subscribe({
      next: (exams) => {
        this.exams = exams;
      },
      error: (error) => {
        console.error('Erro ao carregar exames:', error);
      }
    });
  }

  loadDomains(examId: number) {
    this.adminService.getDomains(examId).subscribe({
      next: (domains) => {
        this.domains = domains;
      },
      error: (error) => {
        console.error('Erro ao carregar domínios:', error);
      }
    });
  }

  onExamChange() {
    const examId = this.questionForm.get('exam_id')?.value;
    if (examId) {
      this.loadDomains(examId);
      this.questionForm.get('domain_id')?.setValue('');
    } else {
      this.domains = [];
    }
  }

  get optionsArray() {
    return this.questionForm.get('options') as FormArray;
  }

  addOption() {
    const optionGroup = this.fb.group({
      text: ['', [Validators.required, Validators.minLength(3)]],
      is_correct: [false],
      explanation: [''],
      order_index: [this.optionsArray.length]
    });
    this.optionsArray.push(optionGroup);
  }

  removeOption(index: number) {
    this.optionsArray.removeAt(index);
    // Reordenar os índices
    this.optionsArray.controls.forEach((control, i) => {
      control.patchValue({ order_index: i });
    });
  }

  openDialog(question?: QuestionWithDetails) {
    this.editingQuestion = question || null;
    if (question) {
      // Carregar domínios do exame
      this.loadDomains(question.exam_id);
      
      this.questionForm.patchValue({
        exam_id: question.exam_id,
        domain_id: question.domain_id,
        statement: question.statement,
        explanation: question.explanation,
        difficulty_level: question.difficulty_level,
        is_active: question.is_active
      });

      // Limpar e recriar opções
      this.optionsArray.clear();
      if (question.options) {
        question.options.forEach(option => {
          const optionGroup = this.fb.group({
            text: [option.text, [Validators.required, Validators.minLength(3)]],
            is_correct: [option.is_correct],
            explanation: [option.explanation || ''],
            order_index: [option.order_index || 0]
          });
          this.optionsArray.push(optionGroup);
        });
      }
    } else {
      this.questionForm.reset({
        difficulty_level: 1,
        is_active: true
      });
      this.optionsArray.clear();
      this.domains = [];
    }
    this.showDialog = true;
  }

  closeDialog() {
    this.showDialog = false;
    this.editingQuestion = null;
    this.questionForm.reset();
    this.optionsArray.clear();
  }

  onSubmit() {
    if (this.questionForm.valid && this.optionsArray.length >= 2) {
      const questionData = this.questionForm.value;
      
      if (this.editingQuestion) {
        // Atualizar questão existente
        this.adminService.updateQuestion(this.editingQuestion.id!, questionData).subscribe({
          next: () => {
            this.snackBar.open('Questão atualizada com sucesso!', 'Fechar', { duration: 3000 });
            this.closeDialog();
            this.loadQuestions();
          },
          error: (error) => {
            console.error('Erro ao atualizar questão:', error);
            this.snackBar.open('Erro ao atualizar questão', 'Fechar', { duration: 3000 });
          }
        });
      } else {
        // Criar nova questão
        this.adminService.createQuestion(questionData).subscribe({
          next: () => {
            this.snackBar.open('Questão criada com sucesso!', 'Fechar', { duration: 3000 });
            this.closeDialog();
            this.loadQuestions();
          },
          error: (error) => {
            console.error('Erro ao criar questão:', error);
            this.snackBar.open('Erro ao criar questão', 'Fechar', { duration: 3000 });
          }
        });
      }
    } else if (this.optionsArray.length < 2) {
      this.snackBar.open('Adicione pelo menos 2 opções', 'Fechar', { duration: 3000 });
    }
  }

  deleteQuestion(question: QuestionWithDetails) {
    if (confirm(`Tem certeza que deseja excluir a questão "${question.statement.substring(0, 50)}..."?`)) {
      this.adminService.deleteQuestion(question.id!).subscribe({
        next: () => {
          this.snackBar.open('Questão excluída com sucesso!', 'Fechar', { duration: 3000 });
          this.loadQuestions();
        },
        error: (error) => {
          console.error('Erro ao excluir questão:', error);
          this.snackBar.open('Erro ao excluir questão', 'Fechar', { duration: 3000 });
        }
      });
    }
  }

  getExamName(examId: number): string {
    const exam = this.exams.find(e => e.id === examId);
    return exam ? exam.title : 'N/A';
  }

  getDomainName(domainId: number): string {
    const domain = this.domains.find(d => d.id === domainId);
    return domain ? domain.name : 'N/A';
  }

  getDifficultyLabel(level: number): string {
    const labels = ['', 'Fácil', 'Médio', 'Difícil', 'Muito Difícil', 'Expert'];
    return labels[level] || 'N/A';
  }

  getDifficultyColor(level: number): string {
    const colors = ['', 'success', 'primary', 'accent', 'warn', 'error'];
    return colors[level] || 'primary';
  }
} 