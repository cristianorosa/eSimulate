import { Component, OnInit } from '@angular/core';
import { DashboardService } from './dashboard.service';
import { AuthService } from '../core/auth.service';
import { take } from 'rxjs/operators';
import { CommonModule } from '@angular/common';
import { MatCardModule } from '@angular/material/card';
import { MatIconModule } from '@angular/material/icon';
import { MatButtonModule } from '@angular/material/button';
import { RouterModule } from '@angular/router';

@Component({
  selector: 'app-dashboard',
  templateUrl: './dashboard.component.html',
  styleUrls: ['./dashboard.component.scss'],
  standalone: true,
  imports: [
    CommonModule,
    MatCardModule,
    MatIconModule,
    MatButtonModule,
    RouterModule,
  ]
})
export class DashboardComponent implements OnInit {
  user: any = null;
  performance: any = null;
  history: any[] = [];
  loading = true;

  constructor(private dashboard: DashboardService, private auth: AuthService) {}

  ngOnInit() {
    const user = this.auth.user();
    console.log('DashboardComponent: Usuário atual:', user);
    
    if (user && user.id) {
      this.user = user;
      console.log('DashboardComponent: Buscando performance para usuário:', user.id);
      
      this.dashboard.getPerformance(user.id).subscribe({
        next: (perf) => {
          console.log('DashboardComponent: Performance recebida:', perf);
          this.performance = perf;
        },
        error: (err) => {
          console.error('DashboardComponent: Erro ao buscar performance:', err);
          this.performance = null;
        },
      });
      
      console.log('DashboardComponent: Buscando histórico para usuário:', user.id);
      this.dashboard.getHistory(user.id).subscribe({
        next: (hist: any[]) => {
          console.log('DashboardComponent: Histórico recebido:', hist);
          this.history = hist || [];
        },
        error: (err) => {
          console.error('DashboardComponent: Erro ao buscar histórico:', err);
          this.history = [];
        },
        complete: () => {
          console.log('DashboardComponent: Carregamento concluído');
          this.loading = false;
        },
      });
    } else {
      console.log('DashboardComponent: Nenhum usuário encontrado');
      this.loading = false;
    }
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
}
