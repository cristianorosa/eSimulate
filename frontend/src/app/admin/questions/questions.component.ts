import { Component, OnInit } from "@angular/core";
import { CommonModule } from "@angular/common";
import { MatTableModule } from "@angular/material/table";
import { MatButtonModule } from "@angular/material/button";
import { MatIconModule } from "@angular/material/icon";
import { MatCardModule } from "@angular/material/card";
import { MatFormFieldModule } from "@angular/material/form-field";
import { MatInputModule } from "@angular/material/input";
import { MatSelectModule } from "@angular/material/select";
import { MatSnackBarModule } from "@angular/material/snack-bar";
import { MatProgressSpinnerModule } from "@angular/material/progress-spinner";
import { MatDialogModule } from "@angular/material/dialog";
import { MatChipsModule } from "@angular/material/chips";
import { MatTooltipModule } from "@angular/material/tooltip";
import { MatPaginatorModule, PageEvent } from "@angular/material/paginator";
import { FormsModule } from "@angular/forms";
import { ReactiveFormsModule, FormBuilder, FormGroup, FormArray, Validators } from "@angular/forms";
import { MatCheckboxModule } from "@angular/material/checkbox";

import { AdminService, Question, Exam, Topic, Option } from "../../core/admin.service";
import { MatSnackBar } from "@angular/material/snack-bar";

interface QuestionWithDetails extends Question {
  exam_name?: string;
  topic_name?: string;
  options_count?: number;
}

@Component({
  selector: "app-questions",
  standalone: true,
  imports: [
    CommonModule,
    MatTableModule,
    MatButtonModule,
    MatIconModule,
    MatCardModule,
    MatFormFieldModule,
    MatInputModule,
    MatSelectModule,
    MatSnackBarModule,
    MatProgressSpinnerModule,
    MatDialogModule,
    MatChipsModule,
    MatTooltipModule,
    MatPaginatorModule,
    FormsModule,
    ReactiveFormsModule,
    MatCheckboxModule,
  ],
  templateUrl: "./questions.component.html",
  styleUrls: ["./questions.component.scss"],
})
export class QuestionsComponent implements OnInit {
  questions: QuestionWithDetails[] = [];
  exams: Exam[] = [];
  topics: Topic[] = [];
  loading = false;
  saving = false;
  showDialog = false;
  editingQuestion: QuestionWithDetails | null = null;
  questionForm: FormGroup;

  // Cache de tópicos para buscar exam_id
  private topicsCache: Topic[] = [];

  displayedColumns: string[] = [
    "id",
    "statement",
    "exam_name",
    "topic_name",
    "question_type",
    "difficulty_level",
    "is_active",
    "actions",
  ];

  pagination: any = null;
  currentPage = 1;
  currentPageSize = 10;
  selectedExamId: number | undefined = undefined;
  selectedTopicId: number | undefined = undefined;

  constructor(
    private adminService: AdminService,
    private snackBar: MatSnackBar,
    private fb: FormBuilder
  ) {
    this.questionForm = this.fb.group({
      exam_id: ["", Validators.required], // Mantido para UX (filtro de tópicos)
      topic_id: ["", Validators.required],
      statement: ["", [Validators.required, Validators.minLength(10)]],
      problem: ["", [Validators.required, Validators.minLength(10)]],
      content_type: ["text", Validators.required],
      explanation: ["", [Validators.required, Validators.minLength(10)]],
      question_type: ["objective", Validators.required],
      difficulty_level: [
        1,
        [Validators.required, Validators.min(1), Validators.max(5)],
      ],
      is_active: [true],
      options: this.fb.array([]),
    });
  }

  ngOnInit() {
    this.loadExams();
    this.loadAllTopics(); // Carregar todos os tópicos para o cache
    this.loadQuestions();
  }

  loadAllTopics() {
    this.adminService.getTopics().subscribe({
      next: (topics) => {
        this.topicsCache = topics;
      },
      error: (error) => {
        console.error("Erro ao carregar todos os tópicos:", error);
      },
    });
  }

  loadQuestions() {
    this.loading = true;
    console.log('Carregando questões com paginação:', { page: this.currentPage, pageSize: this.currentPageSize, examId: this.selectedExamId, topicId: this.selectedTopicId });
    
    this.adminService.getQuestionsPaginated(this.currentPage, this.currentPageSize, this.selectedExamId, this.selectedTopicId).subscribe({
      next: (response) => {
        console.log('Questões carregadas:', response);
        // Garantir que data seja sempre um array
        const questionsData = response.data || [];
        this.questions = questionsData.map((question) => ({
          ...question,
          exam_id: question.exam_id || 0,
          exam_name: question.exam_title || "N/A",
          topic_name: question.topic_name || "N/A",
          options_count: question.options?.length || 0,
        }));
        this.pagination = response.pagination;
        console.log('Paginação:', this.pagination);
        this.loading = false;
      },
      error: (error) => {
        console.error("Erro ao carregar questões:", error);
        this.snackBar.open("Erro ao carregar questões", "Fechar", {
          duration: 3000,
        });
        this.loading = false;
        this.questions = []; // Garantir que questions seja um array vazio em caso de erro
      },
    });
  }

  onPageChange(event: PageEvent) {
    this.currentPage = event.pageIndex + 1;
    this.currentPageSize = event.pageSize;
    this.loadQuestions();
  }

  onExamFilterChange() {
    this.currentPage = 1; // Reset para primeira página
    this.loadTopics(); // Recarregar tópicos baseado no exame selecionado
    this.loadQuestions();
  }

  onTopicFilterChange() {
    this.currentPage = 1; // Reset para primeira página
    this.loadQuestions();
  }

  clearFilters() {
    this.selectedExamId = undefined;
    this.selectedTopicId = undefined;
    this.currentPage = 1;
    this.loadTopics(); // Recarregar todos os tópicos
    this.loadQuestions();
  }

  loadExams() {
    this.adminService.getExams().subscribe({
      next: (exams) => {
        this.exams = exams;
      },
      error: (error) => {
        console.error("Erro ao carregar exames:", error);
      },
    });
  }

  loadTopics() {
    this.adminService.getTopics(this.selectedExamId).subscribe({
      next: (topics) => {
        this.topics = topics;
      },
      error: (error) => {
        console.error("Erro ao carregar tópicos:", error);
        this.snackBar.open("Erro ao carregar tópicos", "Fechar", {
          duration: 3000,
        });
      },
    });
  }

  onExamChange() {
    // Reset topic_id quando exame muda
    this.questionForm.patchValue({ topic_id: null });
    this.loadTopics();
  }

  get optionsArray() {
    return this.questionForm.get("options") as FormArray;
  }

  addOption() {
    const optionGroup = this.fb.group({
      text: ["", [Validators.required, Validators.minLength(1)]],
      is_correct: [false],
      order_index: [this.optionsArray.length],
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
      // Carregar tópicos do exame
      if (question.exam_id) {
        this.selectedExamId = question.exam_id;
        this.loadTopics();
      }

      this.questionForm.patchValue({
        exam_id: question.exam_id || "",
        topic_id: question.topic_id,
        statement: question.statement,
        problem: question.problem,
        content_type: question.content_type,
        explanation: question.explanation,
        question_type: question.question_type,
        difficulty_level: question.difficulty_level,
        is_active: question.is_active,
      });

      // Limpar e recriar opções
      this.optionsArray.clear();
      if (question.options) {
        question.options.forEach((option) => {
          const optionGroup = this.fb.group({
            text: [option.text, [Validators.required, Validators.minLength(1)]],
            is_correct: [option.is_correct],
            order_index: [option.order_index || 0],
          });
          this.optionsArray.push(optionGroup);
        });
      }
    } else {
      this.questionForm.reset({
        content_type: "text",
        question_type: "objective",
        difficulty_level: 1,
        is_active: true,
      });
      this.optionsArray.clear();
      this.topics = [];
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
      // Validar alternativas corretas baseado no tipo de questão
      const questionType = this.questionForm.get("question_type")?.value;
      const correctOptions = this.optionsArray.controls.filter(
        (option) => option.get("is_correct")?.value
      ).length;

      if (questionType === "objective" && correctOptions !== 1) {
        this.snackBar.open(
          "Questões objetivas devem ter exatamente uma alternativa correta",
          "Fechar",
          { duration: 3000 }
        );
        return;
      }

      if (questionType === "multiple_choice" && correctOptions < 1) {
        this.snackBar.open(
          "Questões de múltipla escolha devem ter pelo menos uma alternativa correta",
          "Fechar",
          { duration: 3000 }
        );
        return;
      }

      this.saving = true;
      const questionData = this.questionForm.value;

      // Remover exam_id do payload (não é salvo no backend)
      const { exam_id, ...questionPayload } = questionData;

      if (this.editingQuestion) {
        // Atualizar questão existente
        this.adminService
          .updateQuestion(this.editingQuestion.id!, questionPayload)
          .subscribe({
            next: () => {
              this.snackBar.open("Questão atualizada com sucesso!", "Fechar", {
                duration: 3000,
              });
              this.closeDialog();
              this.loadQuestions();
              this.saving = false;
            },
            error: (error) => {
              console.error("Erro ao atualizar questão:", error);
              this.snackBar.open("Erro ao atualizar questão", "Fechar", {
                duration: 3000,
              });
              this.saving = false;
            },
          });
      } else {
        // Criar nova questão
        this.adminService.createQuestion(questionPayload).subscribe({
          next: () => {
            this.snackBar.open("Questão criada com sucesso!", "Fechar", {
              duration: 3000,
            });
            this.closeDialog();
            this.loadQuestions();
            this.saving = false;
          },
          error: (error) => {
            console.error("Erro ao criar questão:", error);
            this.snackBar.open("Erro ao criar questão", "Fechar", {
              duration: 3000,
            });
            this.saving = false;
          },
        });
      }
    }
  }

  deleteQuestion(question: QuestionWithDetails) {
    if (confirm(`Tem certeza que deseja excluir a questão "${question.statement}"?`)) {
      this.adminService.deleteQuestion(question.id!).subscribe({
        next: () => {
          this.snackBar.open("Questão excluída com sucesso!", "Fechar", {
            duration: 3000,
          });
          this.loadQuestions();
        },
        error: (error) => {
          console.error("Erro ao excluir questão:", error);
          this.snackBar.open("Erro ao excluir questão", "Fechar", {
            duration: 3000,
          });
        },
      });
    }
  }

  getExamName(examId: number | undefined): string {
    if (!examId) return "N/A";
    const exam = this.exams.find((e) => e.id === examId);
    return exam ? exam.title : "N/A";
  }

  getTopicName(topicId: number): string {
    const topic = this.topicsCache.find((t) => t.id === topicId);
    return topic ? topic.name : "N/A";
  }

  getDifficultyLabel(level: number): string {
    const labels = ["", "Fácil", "Médio", "Difícil", "Muito Difícil", "Expert"];
    return labels[level] || "N/A";
  }

  getDifficultyColor(level: number): string {
    const colors = ["", "primary", "accent", "warn", "warn", "warn"];
    return colors[level] || "primary";
  }

  getQuestionTypeLabel(type: string): string {
    switch (type) {
      case "objective":
        return "Objetiva";
      case "multiple_choice":
        return "Múltipla Escolha";
      default:
        return "N/A";
    }
  }

  getContentTypeLabel(type: string): string {
    switch (type) {
      case "text":
        return "Texto";
      case "code":
        return "Código";
      default:
        return "N/A";
    }
  }

  private getExamIdFromTopic(topicId: number): number {
    const topic = this.topicsCache.find((t) => t.id === topicId);
    return topic ? topic.exam_id : 0;
  }
}
