import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ReactiveFormsModule, FormBuilder, FormGroup, Validators } from '@angular/forms';
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
    ReactiveFormsModule,
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
  areaForm: FormGroup;
  displayedColumns = ['id', 'name', 'description', 'created_at', 'actions'];
  pagination: any = null;
  currentPage = 1;
  currentPageSize = 10;
  searchTerm: string = '';

  constructor(
    private adminService: AdminService,
    private snackBar: MatSnackBar,
    private router: Router,
    private fb: FormBuilder
  ) {
    this.areaForm = this.fb.group({
      name: ['', [Validators.required, Validators.minLength(3)]],
      description: ['', [Validators.maxLength(500)]]
    });
  }

  ngOnInit() {
    console.log('AreasComponent: Inicializando...');
    this.loadAreas();
  }

  loadAreas() {
    this.loading = true;
    console.log('Carregando áreas com filtros:', { page: this.currentPage, pageSize: this.currentPageSize, search: this.searchTerm });
    
    this.adminService.getAreasPaginated(this.currentPage, this.currentPageSize).subscribe({
      next: (response: PaginatedResponse<Area>) => {
        console.log('Areas carregadas:', response);
        // Filtrar por termo de busca no frontend
        let filteredAreas = response.data || [];
        if (this.searchTerm.trim()) {
          const searchLower = this.searchTerm.toLowerCase();
          filteredAreas = filteredAreas.filter(area => 
            area.name.toLowerCase().includes(searchLower) ||
            (area.description && area.description.toLowerCase().includes(searchLower))
          );
        }
        this.areas = filteredAreas;
        this.pagination = response.pagination;
        this.loading = false;
      },
      error: (error) => {
        console.error('Erro ao carregar áreas:', error);
        this.snackBar.open('Erro ao carregar áreas', 'Fechar', { duration: 3000 });
        this.loading = false;
        this.areas = []; // Garantir array vazio em caso de erro
      }
    });
  }

  onPageChange(event: PageEvent) {
    this.currentPage = event.pageIndex + 1;
    this.currentPageSize = event.pageSize;
    this.loadAreas();
  }

  onSearchChange() {
    this.currentPage = 1; // Reset para primeira página
    this.loadAreas();
  }

  clearFilters() {
    this.searchTerm = '';
    this.currentPage = 1;
    this.loadAreas();
  }

  openCreateDialog() {
    this.editingArea = null;
    this.areaForm.reset();
    this.showDialog = true;
  }

  editArea(area: Area) {
    this.editingArea = area;
    this.areaForm.patchValue({
      name: area.name,
      description: area.description || ''
    });
    this.showDialog = true;
  }

  closeDialog() {
    this.showDialog = false;
    this.editingArea = null;
    this.areaForm.reset();
  }

  saveArea() {
    if (this.areaForm.invalid) {
      this.snackBar.open('Por favor, preencha todos os campos obrigatórios', 'Fechar', { duration: 3000 });
      return;
    }

    this.saving = true;
    const formData = this.areaForm.value;
    
    const operation = this.editingArea
      ? this.adminService.updateArea(this.editingArea.id!, formData)
      : this.adminService.createArea(formData);

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