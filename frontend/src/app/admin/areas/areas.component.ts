import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { MatCardModule } from '@angular/material/card';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatTableModule } from '@angular/material/table';
import { MatDialogModule } from '@angular/material/dialog';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatSnackBar, MatSnackBarModule } from '@angular/material/snack-bar';
import { MatTooltipModule } from '@angular/material/tooltip';
import { MatPaginatorModule, PageEvent } from '@angular/material/paginator';
import { Router } from '@angular/router';

import { AdminService, Area, PaginatedResponse } from '../../core/admin.service';

@Component({
  selector: 'app-areas',
  standalone: true,
  imports: [
    CommonModule,
    FormsModule,
    MatCardModule,
    MatButtonModule,
    MatIconModule,
    MatTableModule,
    MatDialogModule,
    MatFormFieldModule,
    MatInputModule,
    MatProgressSpinnerModule,
    MatSnackBarModule,
    MatTooltipModule,
    MatPaginatorModule
  ],
  templateUrl: './areas.component.html',
  styleUrls: ['./areas.component.scss']
})
export class AreasComponent implements OnInit {
  areas: Area[] = [];
  loading = true;
  saving = false;
  showDialog = false;
  editingArea: Area | null = null;
  formArea: Area = { name: '', description: '' };
  displayedColumns = ['id', 'name', 'description', 'created_at', 'actions'];
  pagination: any = null;
  currentPage = 1;
  currentPageSize = 10;

  constructor(
    private adminService: AdminService,
    private snackBar: MatSnackBar,
    private router: Router
  ) {}

  ngOnInit() {
    console.log('AreasComponent: Inicializando...');
    this.loadAreas();
  }

  loadAreas() {
    this.loading = true;
    this.adminService.getAreasPaginated(this.currentPage, this.currentPageSize).subscribe({
      next: (response: PaginatedResponse<Area>) => {
        console.log('Areas carregadas:', response);
        this.areas = response.data;
        this.pagination = response.pagination;
        this.loading = false;
      },
      error: (error) => {
        console.error('Erro ao carregar áreas:', error);
        this.snackBar.open('Erro ao carregar áreas', 'Fechar', { duration: 3000 });
        this.loading = false;
      }
    });
  }

  onPageChange(event: PageEvent) {
    this.currentPage = event.pageIndex + 1;
    this.currentPageSize = event.pageSize;
    this.loadAreas();
  }

  openCreateDialog() {
    this.editingArea = null;
    this.formArea = { name: '', description: '' };
    this.showDialog = true;
  }

  editArea(area: Area) {
    this.editingArea = area;
    this.formArea = { ...area };
    this.showDialog = true;
  }

  closeDialog() {
    this.showDialog = false;
    this.editingArea = null;
    this.formArea = { name: '', description: '' };
  }

  saveArea() {
    if (!this.formArea.name.trim()) {
      this.snackBar.open('Nome da área é obrigatório', 'Fechar', { duration: 3000 });
      return;
    }

    this.saving = true;
    const operation = this.editingArea
      ? this.adminService.updateArea(this.editingArea.id!, this.formArea)
      : this.adminService.createArea(this.formArea);

    operation.subscribe({
      next: () => {
        this.snackBar.open(
          `Área ${this.editingArea ? 'atualizada' : 'criada'} com sucesso!`,
          'Fechar',
          { duration: 3000 }
        );
        this.closeDialog();
        this.loadAreas();
      },
      error: (error) => {
        console.error('Erro ao salvar área:', error);
        this.snackBar.open('Erro ao salvar área', 'Fechar', { duration: 3000 });
        this.saving = false;
      }
    });
  }

  deleteArea(area: Area) {
    if (confirm(`Tem certeza que deseja excluir a área "${area.name}"?`)) {
      this.adminService.deleteArea(area.id!).subscribe({
        next: () => {
          this.snackBar.open('Área excluída com sucesso!', 'Fechar', { duration: 3000 });
          this.loadAreas();
        },
        error: (error) => {
          console.error('Erro ao excluir área:', error);
          this.snackBar.open('Erro ao excluir área', 'Fechar', { duration: 3000 });
        }
      });
    }
  }
} 