import { Component, OnInit, ViewChild, ElementRef } from "@angular/core";
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
import { MatAutocompleteModule, MatAutocompleteTrigger } from "@angular/material/autocomplete";
import { Observable, startWith, map, debounceTime, distinctUntilChanged, switchMap } from 'rxjs';

import { AdminService, Question, Exam, Topic, Option } from "../../core/admin.service";
import { MatSnackBar } from "@angular/material/snack-bar";

interface QuestionWithDetails extends Question {
  exam_name?: string;
  topic_name?: string;
  options_count?: number;
}

interface Tag {
  id?: number;
  name: string;
  created_at?: string;
}

interface OptionGridItem {
  id?: number;
  text: string;
  is_correct: boolean;
  order_index: number;
  isEditing?: boolean;
  isNew?: boolean;
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
    MatAutocompleteModule,
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
  
  // Para autocomplete de tópicos
  filteredTopics: Observable<Topic[]> = new Observable();
  allTopics: Topic[] = [];
  
  // Para tags
  availableTags: Tag[] = [];
  selectedTags: Tag[] = [];
  tagInput = '';
  suggestedTags: string[] = [];
  loadingTags = false;
  showTagSuggestions = false;

  // Para grid de opções
  optionsData: OptionGridItem[] = [];
  optionsDisplayedColumns: string[] = ['text', 'actions'];
  editingOptionId: number | null = null;
  showOptionDialog = false;
  editingOption: OptionGridItem | null = null;
  optionDialogTitle = '';
  optionText = '';

  displayedColumns: string[] = [
    "id",
    "statement",
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

  // Para importação de questões
  @ViewChild('fileInput') fileInput!: ElementRef;
  importing = false;

  // Para controle do autocomplete de tópicos
  @ViewChild(MatAutocompleteTrigger) autocompleteTrigger!: MatAutocompleteTrigger;

  constructor(
    private adminService: AdminService,
    private snackBar: MatSnackBar,
    private fb: FormBuilder
  ) {
    this.questionForm = this.fb.group({
      topic_name: ["", Validators.required], // Novo campo para autocomplete
      topic_id: [""], // Será preenchido automaticamente
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
      // Removido FormArray de options - agora gerenciado pela grid
    });
  }

  ngOnInit() {
    this.loadExams();
    this.loadAllTopics();
    // Não carregamos mais todas as tags de uma vez
    this.loadQuestions();
    this.setupTopicAutocomplete();
  }

  loadAllTopics() {
    this.adminService.getTopics().subscribe({
      next: (topics) => {
        this.topicsCache = topics;
        this.allTopics = topics;
      },
      error: (error) => {
        console.error("Erro ao carregar todos os tópicos:", error);
      },
    });
  }

  searchTags(query: string) {
    if (query.length < 2) {
      this.suggestedTags = [];
      return;
    }

    this.loadingTags = true;
    this.adminService.searchTags(query.toLowerCase(), 10).subscribe({
      next: (response: any) => {
        const allSuggestions = response?.data || [];
        
        // Filtrar tags que já foram selecionadas
        const selectedTagNames = this.selectedTags.map(tag => tag.name.toLowerCase());
        this.suggestedTags = allSuggestions.filter((tagName: string) => 
          !selectedTagNames.includes(tagName.toLowerCase())
        );
        
        this.loadingTags = false;
      },
      error: (error) => {
        console.error("Erro ao buscar tags:", error);
        this.suggestedTags = [];
        this.loadingTags = false;
      },
    });
  }

  setupTopicAutocomplete() {
    this.filteredTopics = this.questionForm.get('topic_name')!.valueChanges.pipe(
      startWith(''),
      debounceTime(300),
      distinctUntilChanged(),
      map(value => this.filterTopics(value || ''))
    );
  }

  filterTopics(value: string): Topic[] {
    const filterValue = value.toLowerCase();
    return this.allTopics.filter(topic => 
      topic.name.toLowerCase().includes(filterValue)
    );
  }

  onTopicSelected(topic: Topic) {
    this.questionForm.patchValue({
      topic_name: topic.name,
      topic_id: topic.id
    });
  }

  onTopicInputFocus() {
    // Método pode ser usado para lógicas futuras quando o campo ganha foco
  }

  onTopicInputBlur() {
    // Fechar o autocomplete após um pequeno delay para permitir seleção
    setTimeout(() => {
      if (this.autocompleteTrigger && this.autocompleteTrigger.panelOpen) {
        this.autocompleteTrigger.closePanel();
      }
    }, 200);
  }

  onTopicAutocompleteClosed() {
    // Evento chamado quando o autocomplete é fechado
    // Pode ser usado para lógicas adicionais se necessário
  }

  async onTopicInput() {
    const topicName = this.questionForm.get('topic_name')?.value;
    if (topicName) {
      // Verificar se o tópico já existe
      const existingTopic = this.allTopics.find(t => t.name.toLowerCase() === topicName.toLowerCase());
      if (existingTopic) {
        this.questionForm.patchValue({
          topic_id: existingTopic.id
        });
      } else {
        // Limpar topic_id se não encontrar o tópico
        this.questionForm.patchValue({
          topic_id: null
        });
      }
    }
  }

  async createTopicIfNotExists(topicName: string): Promise<number> {
    // Verificar se o tópico já existe
    const existingTopic = this.allTopics.find(t => t.name.toLowerCase() === topicName.toLowerCase());
    if (existingTopic && existingTopic.id) {
      return existingTopic.id;
    }

    // Criar novo tópico
    try {
      const newTopic = await this.adminService.createTopic({ name: topicName }).toPromise();
      if (newTopic && newTopic.id) {
        this.allTopics.push(newTopic);
        this.snackBar.open(`Tópico "${topicName}" criado com sucesso!`, 'Fechar', { duration: 3000 });
        return newTopic.id;
      }
      throw new Error('Erro ao criar tópico');
    } catch (error) {
      console.error('Erro ao criar tópico:', error);
      this.snackBar.open('Erro ao criar tópico', 'Fechar', { duration: 3000 });
      throw error;
    }
  }

  // Métodos para gerenciar tags
  addTag(tag: Tag) {
    if (!this.selectedTags.find(t => t.name === tag.name)) {
      this.selectedTags.push(tag);
      // Atualizar sugestões após adicionar tag
      this.updateSuggestions();
    }
  }

  removeTag(tag: Tag) {
    this.selectedTags = this.selectedTags.filter(t => t.name !== tag.name);
    // Atualizar sugestões após remover tag
    this.updateSuggestions();
  }

  // Método para atualizar sugestões baseado no input atual
  updateSuggestions() {
    const currentQuery = this.tagInput.trim();
    if (currentQuery.length >= 2) {
      this.searchTags(currentQuery);
    }
  }

  async createTag(tagName: string): Promise<Tag> {
    try {
      const response = await this.adminService.createTag({ name: tagName }).toPromise();
      const newTag = response?.data;
      this.availableTags.push(newTag);
      this.snackBar.open(`Tag "${tagName}" criada com sucesso!`, 'Fechar', { duration: 3000 });
      return newTag;
    } catch (error: any) {
      console.error('Erro ao criar tag:', error);
      
      // Se o erro for de tag já existente, tentar buscar a tag existente
      if (error.error && typeof error.error === 'string' && error.error.includes('already exists')) {
        try {
          // Recarregar tags e tentar encontrar a tag existente
          const response = await this.adminService.getTags().toPromise();
          this.availableTags = (response as any)?.data || response || []; // Suporta tanto formato {data: []} quanto [] direto
          const existingTag = this.availableTags.find(t => t.name.toLowerCase() === tagName.toLowerCase());
          
          if (existingTag) {
            this.snackBar.open(`Tag "${tagName}" já existe e foi encontrada!`, 'Fechar', { duration: 2000 });
            return existingTag;
          }
        } catch (fetchError) {
          console.error('Erro ao buscar tags existentes:', fetchError);
        }
      }
      
      this.snackBar.open('Erro ao criar tag', 'Fechar', { duration: 3000 });
      throw error;
    }
  }

  // Método não é mais necessário, mas mantido para compatibilidade futura
  async findOrCreateTag(tagName: string) {
    try {
      // Usar o novo endpoint que busca ou cria a tag automaticamente
      const response = await this.adminService.findOrCreateTag({ name: tagName }).toPromise();
      const tag = response?.data;
      
      // Adicionar a tag à lista de tags disponíveis se não estiver lá
      if (!this.availableTags.find(t => t.id === tag.id)) {
        this.availableTags.push(tag);
      }
      
      this.addTag(tag);
      this.snackBar.open(`Tag "${tagName}" adicionada com sucesso!`, 'Fechar', { duration: 2000 });
      this.tagInput = '';
    } catch (error) {
      console.error('Erro ao buscar ou criar tag:', error);
      this.snackBar.open('Erro ao processar tag', 'Fechar', { duration: 3000 });
      this.tagInput = '';
    }
  }

  onTagInputKeydown(event: KeyboardEvent) {
    if (event.key === 'Enter' && this.tagInput.trim()) {
      event.preventDefault();
      const tagName = this.tagInput.trim().toLowerCase(); // Sempre em caixa baixa
      
      // Verificar se a tag já está selecionada
      if (this.selectedTags.find(t => t.name.toLowerCase() === tagName)) {
        this.snackBar.open(`Tag "${tagName}" já foi adicionada`, 'Fechar', { duration: 2000 });
        this.tagInput = '';
        this.suggestedTags = [];
        this.showTagSuggestions = false;
        return;
      }
      
      // Criar tag (que agora retorna a existente se já existir)
      this.createTag(tagName).then(tag => {
        this.addTag(tag);
        this.tagInput = '';
        this.suggestedTags = [];
        this.showTagSuggestions = false;
      }).catch(error => {
        console.error('Erro ao criar tag:', error);
      });
    } else if (event.key === 'Escape') {
      // Limpar sugestões ao pressionar Escape
      this.suggestedTags = [];
      this.showTagSuggestions = false;
    }
  }

  onTagInputChange() {
    const query = this.tagInput.trim();
    if (query.length >= 2) {
      this.searchTags(query);
      this.showTagSuggestions = true;
    } else {
      this.suggestedTags = [];
      this.showTagSuggestions = false;
    }
  }

  selectSuggestedTag(tagName: string) {
    // Verificar se a tag já está selecionada
    if (this.selectedTags.find(t => t.name.toLowerCase() === tagName.toLowerCase())) {
      this.snackBar.open(`Tag "${tagName}" já foi adicionada`, 'Fechar', { duration: 2000 });
      this.tagInput = '';
      this.suggestedTags = [];
      this.showTagSuggestions = false;
      return;
    }

    // Criar tag (que agora retorna a existente se já existir)
    this.createTag(tagName).then(tag => {
      this.addTag(tag);
      this.tagInput = '';
      this.suggestedTags = [];
      this.showTagSuggestions = false;
    }).catch(error => {
      console.error('Erro ao criar tag:', error);
    });
  }

  onTagInputFocus() {
    if (this.suggestedTags.length > 0) {
      this.showTagSuggestions = true;
    }
  }

  onTagInputBlur() {
    // Usar setTimeout para permitir que o clique no botão seja processado antes
    setTimeout(() => {
      this.showTagSuggestions = false;
    }, 200);
  }

  associateTagsToQuestion(questionId: number, tags: Tag[]) {
    const tagIds = tags.map(tag => tag.id).filter(id => id !== undefined);
    
    if (tagIds.length === 0) {
      return;
    }

    this.adminService.associateQuestionTags(questionId, tagIds).subscribe({
      next: () => {
        this.snackBar.open(`Tags associadas com sucesso!`, 'Fechar', { duration: 3000 });
      },
      error: (error) => {
        console.error('Erro ao associar tags:', error);
        this.snackBar.open('Erro ao associar tags', 'Fechar', { duration: 3000 });
      }
    });
  }

  loadQuestionTags(questionId: number) {
    this.adminService.getQuestionTags(questionId).subscribe({
      next: (response: any) => {
        if (response?.data && Array.isArray(response.data)) {
          this.selectedTags = response.data;
        }
      },
      error: (error) => {
        console.error('Erro ao carregar tags da questão:', error);
      }
    });
  }

  loadQuestions() {
    this.loading = true;
    console.log('Carregando questões com paginação:', { page: this.currentPage, pageSize: this.currentPageSize, examId: this.selectedExamId, topicId: this.selectedTopicId });
    
    this.adminService.getQuestionsPaginated(this.currentPage, this.currentPageSize, this.selectedExamId, this.selectedTopicId).subscribe({
      next: (response) => {
        console.log('Questões carregadas:', response);
        // Garantir que data seja sempre um array
        const questionsData = response?.data || [];
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

  // Métodos para gerenciar a grid de opções
  openAddOptionDialog() {
    this.optionDialogTitle = 'Adicionar Nova Opção';
    this.optionText = '';
    this.editingOption = null;
    this.showOptionDialog = true;
  }

  openEditOptionDialog(option: OptionGridItem, index: number) {
    this.optionDialogTitle = 'Editar Opção';
    this.optionText = option.text;
    this.editingOption = option;
    this.editingOptionId = index;
    this.showOptionDialog = true;
  }

  closeOptionDialog() {
    this.showOptionDialog = false;
    this.editingOption = null;
    this.editingOptionId = null;
    this.optionText = '';
  }

  saveOptionFromDialog() {
    if (!this.optionText.trim()) {
      this.snackBar.open('O texto da opção é obrigatório', 'Fechar', { duration: 3000 });
      return;
    }

    if (this.editingOption) {
      // Editando opção existente
      this.editingOption.text = this.optionText.trim();
      this.optionsData = [...this.optionsData];
    } else {
      // Adicionando nova opção
      const newOption: OptionGridItem = {
        text: this.optionText.trim(),
        is_correct: false,
        order_index: this.optionsData.length,
        isEditing: false,
        isNew: false
      };
      this.optionsData = [...this.optionsData, newOption];
    }

    this.closeOptionDialog();
  }

  removeOption(index: number) {
    this.optionsData = this.optionsData.filter((_, i) => i !== index);
    // Reordenar os índices
    this.optionsData.forEach((option, i) => {
      option.order_index = i;
    });
    this.editingOptionId = null;
  }

  toggleCorrectOption(option: OptionGridItem, index: number) {
    const questionType = this.questionForm.get("question_type")?.value;
    
    if (questionType === "objective") {
      // Para questões objetivas, apenas uma opção pode ser correta
      this.optionsData.forEach((opt, i) => {
        opt.is_correct = i === index;
      });
    } else {
      // Para múltipla escolha, pode ter várias corretas
      option.is_correct = !option.is_correct;
    }
    
    this.optionsData = [...this.optionsData];
  }

  openDialog(question?: QuestionWithDetails) {
    this.editingQuestion = question || null;
    this.selectedTags = []; // Limpar tags selecionadas
    
    if (question) {
      // Encontrar o nome do tópico
      const topic = this.allTopics.find(t => t.id === question.topic_id);
      const topicName = topic ? topic.name : '';

      this.questionForm.patchValue({
        topic_name: topicName,
        topic_id: question.topic_id,
        statement: question.statement,
        problem: question.problem,
        content_type: question.content_type,
        explanation: question.explanation,
        question_type: question.question_type,
        difficulty_level: question.difficulty_level,
        is_active: question.is_active,
      });

      // Carregar tags da questão
      if (question.id) {
        this.loadQuestionTags(question.id);
      }

      // Carregar opções na grid
      this.optionsData = question.options ? question.options.map((option, index) => ({
        id: option.id,
        text: option.text,
        is_correct: option.is_correct,
        order_index: option.order_index || index,
        isEditing: false,
        isNew: false
      })) : [];
    } else {
      this.questionForm.reset({
        topic_name: "",
        topic_id: "",
        content_type: "text",
        question_type: "objective",
        difficulty_level: 1,
        is_active: true,
      });
      this.optionsData = [];
    }
    this.showDialog = true;
  }

  closeDialog() {
    this.showDialog = false;
    this.editingQuestion = null;
    this.questionForm.reset();
    this.optionsData = [];
    this.editingOptionId = null;
    this.closeOptionDialog(); // Fechar popup de opções também
  }

  async onSubmit() {
    // Verificar se há opções em edição
    if (this.editingOptionId !== null) {
      this.snackBar.open('Salve ou cancele a edição da opção antes de continuar', 'Fechar', { duration: 3000 });
      return;
    }

    if (this.questionForm.valid && this.optionsData.length >= 2) {
      // Validar alternativas corretas baseado no tipo de questão
      const questionType = this.questionForm.get("question_type")?.value;
      const correctOptions = this.optionsData.filter(option => option.is_correct).length;

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

      try {
        // Garantir que temos um topic_id válido
        const topicName = questionData.topic_name;
        if (!questionData.topic_id && topicName) {
          // Criar tópico se não existir
          const topicId = await this.createTopicIfNotExists(topicName);
          questionData.topic_id = topicId;
        }

        // Remover topic_name do payload e adicionar opções
        const { topic_name, ...questionPayload } = questionData;
        
        // Adicionar opções da grid ao payload
        questionPayload.options = this.optionsData.map(option => ({
          id: option.id,
          text: option.text.trim(),
          is_correct: option.is_correct,
          order_index: option.order_index
        }));

        if (this.editingQuestion) {
          // Atualizar questão existente
          this.adminService
            .updateQuestion(this.editingQuestion.id!, questionPayload)
            .subscribe({
              next: () => {
                this.snackBar.open("Questão atualizada com sucesso!", "Fechar", {
                  duration: 3000,
                });
                
                // Associar tags à questão atualizada
                if (this.selectedTags.length > 0 && this.editingQuestion?.id) {
                  this.associateTagsToQuestion(this.editingQuestion.id, this.selectedTags);
                }
                
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
            next: (response: any) => {
              this.snackBar.open("Questão criada com sucesso!", "Fechar", {
                duration: 3000,
              });
              
              // Associar tags à questão
              if (this.selectedTags.length > 0 && response.data?.id) {
                this.associateTagsToQuestion(response.data.id, this.selectedTags);
              }
              
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
      } catch (error) {
        console.error("Erro ao processar questão:", error);
        this.snackBar.open("Erro ao processar questão", "Fechar", {
          duration: 3000,
        });
        this.saving = false;
      }
    } else {
      let message = "Por favor, corrija os seguintes problemas:\n";
      
      if (this.questionForm.invalid) {
        message += "- Preencha todos os campos obrigatórios\n";
      }
      
      if (this.optionsData.length < 2) {
        message += "- Adicione pelo menos 2 opções de resposta\n";
      }
      
      // Verificar se há opções sem texto
      const emptyOptions = this.optionsData.filter(opt => !opt.text.trim()).length;
      if (emptyOptions > 0) {
        message += `- ${emptyOptions} opção(ões) estão sem texto\n`;
      }

      this.snackBar.open(message, "Fechar", {
        duration: 5000,
      });
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

  // Métodos para importação de questões
  importQuestions() {
    this.fileInput.nativeElement.click();
  }

  onFileSelected(event: any) {
    const file = event.target.files[0];
    if (!file) return;

    if (file.type !== 'application/json' && !file.name.endsWith('.json')) {
      this.snackBar.open('Por favor, selecione um arquivo JSON válido', 'Fechar', { duration: 3000 });
      return;
    }

    const reader = new FileReader();
    reader.onload = (e: any) => {
      try {
        const jsonData = JSON.parse(e.target.result);
        this.processImportData(jsonData);
      } catch (error) {
        this.snackBar.open('Erro ao processar arquivo JSON', 'Fechar', { duration: 3000 });
        console.error('Erro ao processar JSON:', error);
      }
    };
    reader.readAsText(file);
  }

  private async processImportData(data: any) {
    this.importing = true;
    
    try {
      // Validar estrutura do JSON
      if (!data.exam || !data.area || !data.topics) {
        throw new Error('Estrutura JSON inválida. Deve conter exam, area e topics.');
      }

      // Enviar dados para o backend
      const response = await this.adminService.importQuestions(data).toPromise();
      
      this.snackBar.open(
        `Importação concluída! ${response.result.questions_created} questões importadas com sucesso.`, 
        'Fechar', 
        { duration: 5000 }
      );
      
      // Recarregar dados
      this.loadQuestions();
      this.loadExams();
      this.loadAllTopics();
      
    } catch (error) {
      console.error('Erro na importação:', error);
      this.snackBar.open(
        `Erro na importação: ${error instanceof Error ? error.message : 'Erro desconhecido'}`, 
        'Fechar', 
        { duration: 5000 }
      );
    } finally {
      this.importing = false;
      // Limpar input
      this.fileInput.nativeElement.value = '';
    }
  }


}
