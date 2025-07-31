import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { MatCardModule } from '@angular/material/card';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatChipsModule } from '@angular/material/chips';
import { MatDividerModule } from '@angular/material/divider';
import { MatProgressBarModule } from '@angular/material/progress-bar';
import { MatTableModule } from '@angular/material/table';
import { MatPaginatorModule } from '@angular/material/paginator';
import { MatSortModule } from '@angular/material/sort';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatSelectModule } from '@angular/material/select';
import { MatOptionModule } from '@angular/material/core';
import { FormsModule } from '@angular/forms';
import { AuthService } from '../core/auth.service';
import { HttpClient } from '@angular/common/http';

@Component({
  selector: 'app-historico',
  standalone: true,
  imports: [
    CommonModule,
    MatCardModule,
    MatButtonModule,
    MatIconModule,
    MatChipsModule,
    MatDividerModule,
    MatProgressBarModule,
    MatTableModule,
    MatPaginatorModule,
    MatSortModule,
    MatFormFieldModule,
    MatSelectModule,
    MatOptionModule,
    FormsModule
  ],
  templateUrl: './historico.component.html',
  styleUrls: ['./historico.component.scss']
})
export class HistoricoComponent implements OnInit {
  history: any[] = [];
  loading = false;
  displayedColumns: string[] = ['quiz', 'score', 'time', 'date', 'actions'];
  
  // Filtros
  selectedPeriod = 'all';
  selectedQuiz = 'all';
  
  periods = [
    { value: 'all', label: 'Todos os períodos' },
    { value: 'week', label: 'Última semana' },
    { value: 'month', label: 'Último mês' },
    { value: 'year', label: 'Último ano' }
  ];

  constructor(
    private auth: AuthService,
    private http: HttpClient
  ) {}

  ngOnInit() {
    this.loadHistory();
  }

  loadHistory() {
    this.loading = true;
    const user = this.auth.user();
    
    if (user) {
      // Simular dados de histórico (implementar chamada real para API)
      setTimeout(() => {
        this.history = [
          {
            id: 1,
            quiz_title: 'Simulado de Matemática',
            score: 85,
            total_questions: 20,
            correct_answers: 17,
            time_spent: '45:30',
            completed_at: '2024-01-15T10:30:00Z',
            status: 'completed'
          },
          {
            id: 2,
            quiz_title: 'Simulado de Português',
            score: 72,
            total_questions: 15,
            correct_answers: 11,
            time_spent: '32:15',
            completed_at: '2024-01-12T14:20:00Z',
            status: 'completed'
          },
          {
            id: 3,
            quiz_title: 'Simulado de História',
            score: 90,
            total_questions: 25,
            correct_answers: 23,
            time_spent: '58:45',
            completed_at: '2024-01-10T09:15:00Z',
            status: 'completed'
          },
          {
            id: 4,
            quiz_title: 'Simulado de Geografia',
            score: 68,
            total_questions: 18,
            correct_answers: 12,
            time_spent: '41:20',
            completed_at: '2024-01-08T16:45:00Z',
            status: 'completed'
          }
        ];
        this.loading = false;
      }, 1000);
    }
  }

  getScoreColor(score: number): string {
    if (score >= 80) return 'success';
    if (score >= 60) return 'warning';
    return 'error';
  }

  getScoreLabel(score: number): string {
    if (score >= 80) return 'Excelente';
    if (score >= 60) return 'Bom';
    return 'Precisa Melhorar';
  }

  formatDate(dateString: string): string {
    const date = new Date(dateString);
    return date.toLocaleDateString('pt-BR', {
      day: '2-digit',
      month: '2-digit',
      year: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    });
  }

  viewDetails(item: any) {
    // Implementar visualização detalhada
    console.log('Ver detalhes:', item);
  }

  retakeQuiz(item: any) {
    // Implementar retomada do simulado
    console.log('Retomar simulado:', item);
  }

  getFilteredHistory() {
    let filtered = this.history || [];
    
    // Aplicar filtros se necessário
    if (this.selectedPeriod !== 'all') {
      // Implementar filtro por período
    }
    
    return filtered;
  }

  getStats() {
    const history = this.getFilteredHistory();
    if (history.length === 0) return { total: 0, average: 0, best: 0, totalTime: '0:00' };
    
    const scores = history.map(h => h.score);
    const totalTime = history.reduce((acc, h) => {
      const [min, sec] = h.time_spent.split(':').map(Number);
      return acc + min * 60 + sec;
    }, 0);
    
    return {
      total: history.length,
      average: Math.round(scores.reduce((a, b) => a + b, 0) / scores.length),
      best: Math.max(...scores),
      totalTime: `${Math.floor(totalTime / 60)}:${(totalTime % 60).toString().padStart(2, '0')}`
    };
  }
} 