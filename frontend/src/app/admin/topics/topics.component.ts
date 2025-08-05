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
import { MatTooltipModule } from '@angular/material/tooltip';
import { FormsModule, ReactiveFormsModule, FormBuilder, FormGroup, Validators } from '@angular/forms';
import { AdminService, Topic, Exam } from '../../core/admin.service';
import { MatSnackBar } from '@angular/material/snack-bar';

interface TopicWithExam extends Topic {
  exam_name?: string;
  updated_at?: string;
}

@Component({
  selector: 'app-topics',
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
    MatTooltipModule,
    FormsModule,
    ReactiveFormsModule
  ],
  templateUrl: './topics.component.html',
  styleUrls: ['./topics.component.scss']
})
export class TopicsComponent implements OnInit {
  topics: TopicWithExam[] = [];
  exams: Exam[] = [];
  loading = false;
  saving = false;
  showDialog = false;
  editingTopic: TopicWithExam | null = null;
  topicForm: FormGroup;

  displayedColumns: string[] = ['id', 'name', 'exam_name', 'questions_count', 'created_at', 'actions'];

  constructor(
    private adminService: AdminService,
    private snackBar: MatSnackBar,
    private fb: FormBuilder
  ) {
    this.topicForm = this.fb.group({
      name: ['', [Validators.required, Validators.minLength(3)]],
      exam_id: ['', Validators.required],
      weight_percentage: [0, [Validators.required, Validators.min(0), Validators.max(100)]],
      order_index: [0, [Validators.required, Validators.min(0)]],
      questions_count: [0, [Validators.required, Validators.min(0)]]
    });
  }

  ngOnInit() {
    this.loadExams();
  }

  loadTopics(examId?: number) {
    this.loading = true;
    
    this.adminService.getTopics(examId || 0).subscribe({
      next: (topics) => {
        // Verificar se topics não é null ou undefined
        if (topics && Array.isArray(topics)) {
          this.topics = topics.map(topic => ({
            ...topic,
            exam_name: this.getExamName(topic.exam_id)
          }));
        } else {
          this.topics = [];
        }
        
        this.loading = false;
      },
      error: (error) => {
        console.error('Erro ao carregar tópicos:', error);
        this.snackBar.open('Erro ao carregar tópicos', 'Fechar', { duration: 3000 });
        this.loading = false;
        this.topics = [];
      }
    });
  }

  loadExams() {
    this.adminService.getExams().subscribe({
      next: (exams) => {
        this.exams = exams;
        // Carregar tópicos após carregar exames
        this.loadTopics();
      },
      error: (error) => {
        console.error('Erro ao carregar exames:', error);
      }
    });
  }

  openDialog(topic?: TopicWithExam) {
    this.editingTopic = topic || null;
    if (topic) {
      this.topicForm.patchValue({
        name: topic.name,
        exam_id: topic.exam_id,
        weight_percentage: topic.weight_percentage,
        order_index: topic.order_index,
        questions_count: topic.questions_count || 0
      });
    } else {
      this.topicForm.reset({
        weight_percentage: 0,
        order_index: 0,
        questions_count: 0
      });
    }
    this.showDialog = true;
  }

  closeDialog() {
    this.showDialog = false;
    this.editingTopic = null;
    this.topicForm.reset();
  }

  onSubmit() {
    if (this.topicForm.valid) {
      this.saving = true;
      const topicData = this.topicForm.value;
      
      if (this.editingTopic) {
        // Atualizar tópico existente
        this.adminService.updateTopic(this.editingTopic.id!, topicData).subscribe({
          next: () => {
            this.snackBar.open('Tópico atualizado com sucesso!', 'Fechar', { duration: 3000 });
            this.closeDialog();
            this.loadTopics();
            this.saving = false;
          },
          error: (error) => {
            console.error('Erro ao atualizar tópico:', error);
            this.snackBar.open('Erro ao atualizar tópico', 'Fechar', { duration: 3000 });
            this.saving = false;
          }
        });
      } else {
        // Criar novo tópico
        this.adminService.createTopic(topicData).subscribe({
          next: () => {
            this.snackBar.open('Tópico criado com sucesso!', 'Fechar', { duration: 3000 });
            this.closeDialog();
            this.loadTopics();
            this.saving = false;
          },
          error: (error) => {
            console.error('Erro ao criar tópico:', error);
            this.snackBar.open('Erro ao criar tópico', 'Fechar', { duration: 3000 });
            this.saving = false;
          }
        });
      }
    }
  }

  deleteTopic(topic: TopicWithExam) {
    if (confirm(`Tem certeza que deseja excluir o tópico "${topic.name}"?`)) {
      this.adminService.deleteTopic(topic.id!).subscribe({
        next: () => {
          this.snackBar.open('Tópico excluído com sucesso!', 'Fechar', { duration: 3000 });
          this.loadTopics();
        },
        error: (error) => {
          console.error('Erro ao excluir tópico:', error);
          this.snackBar.open('Erro ao excluir tópico', 'Fechar', { duration: 3000 });
        }
      });
    }
  }

  getExamName(examId: number): string {
    const exam = this.exams.find(e => e.id === examId);
    return exam ? exam.title : 'N/A';
  }
} 