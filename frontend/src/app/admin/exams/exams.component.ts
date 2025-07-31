import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { MatCardModule } from '@angular/material/card';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatTableModule } from '@angular/material/table';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatSelectModule } from '@angular/material/select';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatSnackBar, MatSnackBarModule } from '@angular/material/snack-bar';
import { MatTooltipModule } from '@angular/material/tooltip';
import { MatPaginatorModule, PageEvent } from '@angular/material/paginator';
import { Router } from '@angular/router';

import { AdminService, Exam, Area, PaginatedResponse } from '../../core/admin.service';

@Component({
  selector: 'app-exams',
  standalone: true,
  imports: [
    CommonModule,
    FormsModule,
    MatCardModule,
    MatButtonModule,
    MatIconModule,
    MatTableModule,
    MatFormFieldModule,
    MatInputModule,
    MatSelectModule,
    MatProgressSpinnerModule,
    MatSnackBarModule,
    MatTooltipModule,
    MatPaginatorModule
  ],
  templateUrl: './exams.component.html',
  styleUrls: ['./exams.component.scss']
})
export class ExamsComponent implements OnInit {
  exams: Exam[] = [];
  areas: Area[] = [];
  loading = true;
  saving = false;
  showDialog = false;
  editingExam: Exam | null = null;
  formExam: Exam = {
    title: '',
    description: '',
    area_id: 0,
    max_time_minutes: 60,
    passing_score: 70.0,
    is_active: true
  };
  displayedColumns = ['id', 'title', 'area', 'max_time', 'passing_score', 'is_active', 'actions'];
  pagination: any = null;
  currentPage = 1;
  currentPageSize = 10;

  constructor(
    private adminService: AdminService,
    private snackBar: MatSnackBar,
    private router: Router
  ) {}

  ngOnInit() {
    console.log('ExamsComponent: Inicializando...');
    this.loadAreas();
    this.loadExams();
  }

  loadAreas() {
    this.adminService.getAreas().subscribe({
      next: (areas) => {
        this.areas = areas;
      },
      error: (error) => {
        console.error('Erro ao carregar áreas:', error);
        this.snackBar.open('Erro ao carregar áreas', 'Fechar', { duration: 3000 });
      }
    });
  }

  loadExams() {
    this.loading = true;
    this.adminService.getExamsPaginated(this.currentPage, this.currentPageSize).subscribe({
      next: (response: PaginatedResponse<Exam>) => {
        console.log('Exams carregados:', response);
        this.exams = response.data;
        this.pagination = response.pagination;
        this.loading = false;
      },
      error: (error) => {
        console.error('Erro ao carregar exames:', error);
        this.snackBar.open('Erro ao carregar exames', 'Fechar', { duration: 3000 });
        this.loading = false;
      }
    });
  }

  onPageChange(event: PageEvent) {
    this.currentPage = event.pageIndex + 1;
    this.currentPageSize = event.pageSize;
    this.loadExams();
  }

  getAreaName(areaId: number): string {
    const area = this.areas.find(a => a.id === areaId);
    return area ? area.name : 'N/A';
  }

  openCreateDialog() {
    this.editingExam = null;
    this.formExam = {
      title: '',
      description: '',
      area_id: 0,
      max_time_minutes: 60,
      passing_score: 70.0,
      is_active: true
    };
    this.showDialog = true;
  }

  editExam(exam: Exam) {
    this.editingExam = exam;
    this.formExam = { ...exam };
    this.showDialog = true;
  }

  closeDialog() {
    this.showDialog = false;
    this.editingExam = null;
    this.formExam = {
      title: '',
      description: '',
      area_id: 0,
      max_time_minutes: 60,
      passing_score: 70.0,
      is_active: true
    };
  }

  saveExam() {
    if (!this.formExam.title.trim()) {
      this.snackBar.open('Título do exame é obrigatório', 'Fechar', { duration: 3000 });
      return;
    }

    if (!this.formExam.area_id) {
      this.snackBar.open('Área é obrigatória', 'Fechar', { duration: 3000 });
      return;
    }

    this.saving = true;
    const operation = this.editingExam
      ? this.adminService.updateExam(this.editingExam.id!, this.formExam)
      : this.adminService.createExam(this.formExam);

    operation.subscribe({
      next: () => {
        this.snackBar.open(
          `Exame ${this.editingExam ? 'atualizado' : 'criado'} com sucesso!`,
          'Fechar',
          { duration: 3000 }
        );
        this.closeDialog();
        this.loadExams();
      },
      error: (error) => {
        console.error('Erro ao salvar exame:', error);
        this.snackBar.open('Erro ao salvar exame', 'Fechar', { duration: 3000 });
        this.saving = false;
      }
    });
  }

  deleteExam(exam: Exam) {
    if (confirm(`Tem certeza que deseja excluir o exame "${exam.title}"?`)) {
      this.adminService.deleteExam(exam.id!).subscribe({
        next: () => {
          this.snackBar.open('Exame excluído com sucesso!', 'Fechar', { duration: 3000 });
          this.loadExams();
        },
        error: (error) => {
          console.error('Erro ao excluir exame:', error);
          this.snackBar.open('Erro ao excluir exame', 'Fechar', { duration: 3000 });
        }
      });
    }
  }
} 