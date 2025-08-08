import { Component, OnInit } from "@angular/core";
import { CommonModule } from "@angular/common";
import {
  ReactiveFormsModule,
  FormBuilder,
  FormGroup,
  Validators,
} from "@angular/forms";
import { FormsModule } from "@angular/forms";
import { MatCardModule } from "@angular/material/card";
import { MatButtonModule } from "@angular/material/button";
import { MatIconModule } from "@angular/material/icon";
import { MatTableModule } from "@angular/material/table";
import { MatFormFieldModule } from "@angular/material/form-field";
import { MatInputModule } from "@angular/material/input";
import { MatSelectModule } from "@angular/material/select";
import { MatProgressSpinnerModule } from "@angular/material/progress-spinner";
import { MatSnackBar, MatSnackBarModule } from "@angular/material/snack-bar";
import { MatTooltipModule } from "@angular/material/tooltip";
import { MatPaginatorModule, PageEvent } from "@angular/material/paginator";
import { Router } from "@angular/router";

import {
  AdminService,
  Exam,
  Area,
  PaginatedResponse,
} from "../../core/admin.service";
import { MatChipsModule } from "@angular/material/chips";

@Component({
  selector: "app-exams",
  standalone: true,
  imports: [
    CommonModule,
    ReactiveFormsModule,
    FormsModule,
    MatCardModule,
    MatButtonModule,
    MatIconModule,
    MatTableModule,
    MatFormFieldModule,
    MatInputModule,
    MatChipsModule,
    MatSelectModule,
    MatProgressSpinnerModule,
    MatSnackBarModule,
    MatTooltipModule,
    MatPaginatorModule,
  ],
  templateUrl: "./exams.component.html",
  styleUrls: ["./exams.component.scss"],
})
export class ExamsComponent implements OnInit {
  exams: Exam[] = [];
  areas: Area[] = [];
  loading = true;
  saving = false;
  showDialog = false;
  editingExam: Exam | null = null;
  examForm: FormGroup;
  displayedColumns = [
    "id",
    "title",
    "area",
    "max_time",
    "passing_score",
    "is_active",
    "actions",
  ];
  pagination: any = null;
  currentPage = 1;
  currentPageSize = 10;
  selectedAreaId: number | undefined = undefined;
  selectedStatus: boolean | undefined = undefined;

  constructor(
    private adminService: AdminService,
    private snackBar: MatSnackBar,
    private router: Router,
    private fb: FormBuilder
  ) {
    this.examForm = this.fb.group({
      title: ["", [Validators.required, Validators.minLength(3)]],
      description: ["", [Validators.maxLength(500)]],
      area_id: ["", [Validators.required]],
      max_time_minutes: [
        60,
        [Validators.required, Validators.min(1), Validators.max(480)],
      ],
      passing_score: [
        70.0,
        [Validators.required, Validators.min(0), Validators.max(100)],
      ],
      questions_count: [0, [Validators.min(0)]],
      is_active: [true],
      created_by: [1], // TODO: Pegar do usuário autenticado
    });
  }

  ngOnInit() {
    console.log("ExamsComponent: Inicializando...");
    this.loadAreas();
    this.loadExams();
  }

  loadAreas() {
    this.adminService.getAreas().subscribe({
      next: (areas) => {
        this.areas = areas;
      },
      error: (error) => {
        console.error("Erro ao carregar áreas:", error);
        this.snackBar.open("Erro ao carregar áreas", "Fechar", {
          duration: 3000,
        });
      },
    });
  }

  loadExams() {
    this.loading = true;
    console.log("Carregando exames com filtros:", {
      page: this.currentPage,
      pageSize: this.currentPageSize,
      areaId: this.selectedAreaId,
      status: this.selectedStatus,
    });

    this.adminService
      .getExamsPaginated(
        this.currentPage,
        this.currentPageSize,
        this.selectedAreaId
      )
      .subscribe({
        next: (response: PaginatedResponse<Exam>) => {
          console.log("Exams carregados:", response);
          // Filtrar por status no frontend
          let filteredExams = response.data || [];
          if (this.selectedStatus !== undefined) {
            filteredExams = filteredExams.filter(
              (exam) => exam.is_active === this.selectedStatus
            );
          }
          this.exams = filteredExams;
          this.pagination = response.pagination;
          this.loading = false;
        },
        error: (error) => {
          console.error("Erro ao carregar exames:", error);
          this.snackBar.open("Erro ao carregar exames", "Fechar", {
            duration: 3000,
          });
          this.loading = false;
          this.exams = []; // Garantir array vazio em caso de erro
        },
      });
  }

  onPageChange(event: PageEvent) {
    this.currentPage = event.pageIndex + 1;
    this.currentPageSize = event.pageSize;
    this.loadExams();
  }

  onAreaFilterChange() {
    this.currentPage = 1; // Reset para primeira página
    this.loadExams();
  }

  onStatusFilterChange() {
    this.currentPage = 1; // Reset para primeira página
    this.loadExams();
  }

  clearFilters() {
    this.selectedAreaId = undefined;
    this.selectedStatus = undefined;
    this.currentPage = 1;
    this.loadExams();
  }

  getAreaName(areaId: number): string {
    const area = this.areas.find((a) => a.id === areaId);
    return area ? area.name : "N/A";
  }

  openDialog(exam?: Exam) {
    this.editingExam = exam || null;
    if (exam) {
      this.examForm.patchValue({
        title: exam.title,
        description: exam.description || "",
        area_id: exam.area_id,
        max_time_minutes: exam.max_time_minutes,
        passing_score: exam.passing_score,
        questions_count: exam.questions_count || 0,
        is_active: exam.is_active,
        created_by: exam.created_by || 1,
      });
    } else {
      this.examForm.reset({
        max_time_minutes: 60,
        passing_score: 70.0,
        questions_count: 0,
        is_active: true,
        created_by: 1,
      });
    }
    this.showDialog = true;
  }

  closeDialog() {
    this.showDialog = false;
    this.editingExam = null;
    this.examForm.reset();
  }

  saveExam() {
    if (this.examForm.invalid) {
      this.snackBar.open(
        "Por favor, preencha todos os campos obrigatórios",
        "Fechar",
        { duration: 3000 }
      );
      return;
    }

    this.saving = true;
    const formData = this.examForm.value;

    const operation = this.editingExam
      ? this.adminService.updateExam(this.editingExam.id!, formData)
      : this.adminService.createExam(formData);

    operation.subscribe({
      next: () => {
        this.snackBar.open(
          `Exame ${this.editingExam ? "atualizado" : "criado"} com sucesso!`,
          "Fechar",
          { duration: 3000 }
        );
        this.closeDialog();
        this.loadExams();
      },
      error: (error) => {
        console.error("Erro ao salvar exame:", error);
        this.snackBar.open("Erro ao salvar exame", "Fechar", {
          duration: 3000,
        });
        this.saving = false;
      },
    });
  }

  deleteExam(exam: Exam) {
    if (confirm(`Tem certeza que deseja excluir o exame "${exam.title}"?`)) {
      this.adminService.deleteExam(exam.id!).subscribe({
        next: () => {
          this.snackBar.open("Exame excluído com sucesso!", "Fechar", {
            duration: 3000,
          });
          this.loadExams();
        },
        error: (error) => {
          console.error("Erro ao excluir exame:", error);
          this.snackBar.open("Erro ao excluir exame", "Fechar", {
            duration: 3000,
          });
        },
      });
    }
  }
}
