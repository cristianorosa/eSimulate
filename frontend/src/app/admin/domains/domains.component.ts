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
import { AdminService, Domain, Exam } from '../../core/admin.service';
import { MatSnackBar } from '@angular/material/snack-bar';

interface DomainWithExam extends Domain {
  exam_name?: string;
  updated_at?: string;
}

@Component({
  selector: 'app-domains',
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
  templateUrl: './domains.component.html',
  styleUrls: ['./domains.component.scss']
})
export class DomainsComponent implements OnInit {
  domains: DomainWithExam[] = [];
  exams: Exam[] = [];
  loading = false;
  showDialog = false;
  editingDomain: DomainWithExam | null = null;
  domainForm: FormGroup;

  displayedColumns: string[] = ['id', 'name', 'description', 'exam_name', 'created_at', 'actions'];

  constructor(
    private adminService: AdminService,
    private snackBar: MatSnackBar,
    private fb: FormBuilder
  ) {
    this.domainForm = this.fb.group({
      name: ['', [Validators.required, Validators.minLength(3)]],
      description: ['', [Validators.required, Validators.minLength(10)]],
      exam_id: ['', Validators.required],
      weight_percentage: [0, [Validators.required, Validators.min(0), Validators.max(100)]],
      order_index: [0, [Validators.required, Validators.min(0)]]
    });
  }

  ngOnInit() {
    this.loadExams();
  }

  loadDomains(examId?: number) {
    this.loading = true;
    this.adminService.getDomains(examId || 0).subscribe({
      next: (domains) => {
        // Verificar se domains não é null ou undefined
        if (domains && Array.isArray(domains)) {
          this.domains = domains.map(domain => ({
            ...domain,
            exam_name: this.getExamName(domain.exam_id)
          }));
        } else {
          this.domains = [];
        }
        this.loading = false;
      },
      error: (error) => {
        console.error('Erro ao carregar domínios:', error);
        this.snackBar.open('Erro ao carregar domínios', 'Fechar', { duration: 3000 });
        this.loading = false;
        this.domains = [];
      }
    });
  }

  loadExams() {
    this.adminService.getExams().subscribe({
      next: (exams) => {
        this.exams = exams;
        // Carregar domínios após carregar exames
        this.loadDomains();
      },
      error: (error) => {
        console.error('Erro ao carregar exames:', error);
      }
    });
  }

  openDialog(domain?: DomainWithExam) {
    this.editingDomain = domain || null;
    if (domain) {
      this.domainForm.patchValue({
        name: domain.name,
        description: domain.description,
        exam_id: domain.exam_id,
        weight_percentage: domain.weight_percentage,
        order_index: domain.order_index
      });
    } else {
      this.domainForm.reset({
        weight_percentage: 0,
        order_index: 0
      });
    }
    this.showDialog = true;
  }

  closeDialog() {
    this.showDialog = false;
    this.editingDomain = null;
    this.domainForm.reset();
  }

  onSubmit() {
    if (this.domainForm.valid) {
      const domainData = this.domainForm.value;
      
      if (this.editingDomain) {
        // Atualizar domínio existente
        this.adminService.updateDomain(this.editingDomain.id!, domainData).subscribe({
          next: () => {
            this.snackBar.open('Domínio atualizado com sucesso!', 'Fechar', { duration: 3000 });
            this.closeDialog();
            this.loadDomains();
          },
          error: (error) => {
            console.error('Erro ao atualizar domínio:', error);
            this.snackBar.open('Erro ao atualizar domínio', 'Fechar', { duration: 3000 });
          }
        });
      } else {
        // Criar novo domínio
        this.adminService.createDomain(domainData).subscribe({
          next: () => {
            this.snackBar.open('Domínio criado com sucesso!', 'Fechar', { duration: 3000 });
            this.closeDialog();
            this.loadDomains();
          },
          error: (error) => {
            console.error('Erro ao criar domínio:', error);
            this.snackBar.open('Erro ao criar domínio', 'Fechar', { duration: 3000 });
          }
        });
      }
    }
  }

  deleteDomain(domain: DomainWithExam) {
    if (confirm(`Tem certeza que deseja excluir o domínio "${domain.name}"?`)) {
      this.adminService.deleteDomain(domain.id!).subscribe({
        next: () => {
          this.snackBar.open('Domínio excluído com sucesso!', 'Fechar', { duration: 3000 });
          this.loadDomains();
        },
        error: (error) => {
          console.error('Erro ao excluir domínio:', error);
          this.snackBar.open('Erro ao excluir domínio', 'Fechar', { duration: 3000 });
        }
      });
    }
  }

  getExamName(examId: number): string {
    const exam = this.exams.find(e => e.id === examId);
    return exam ? exam.title : 'N/A';
  }
} 