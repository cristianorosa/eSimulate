import { Component, OnInit } from '@angular/core';
import { Router, RouterModule } from '@angular/router';
import { QuestoesService } from './questoes.service';
import { CommonModule } from '@angular/common';
import { MatCardModule } from '@angular/material/card';
import { MatListModule } from '@angular/material/list';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatProgressBarModule } from '@angular/material/progress-bar';

@Component({
  selector: 'app-questoes',
  templateUrl: './questoes.component.html',
  styleUrls: ['./questoes.component.scss'],
  standalone: true,
  imports: [
    CommonModule,
    MatCardModule,
    MatListModule,
    MatButtonModule,
    MatIconModule,
    MatProgressBarModule,
    RouterModule,
  ]
})
export class QuestoesComponent implements OnInit {
  questoes: any[] = [];
  loading = true;

  constructor(private questoesService: QuestoesService, private router: Router) {}

  ngOnInit() {
    this.questoesService.getQuestoes().subscribe({
      next: (data) => {
        this.questoes = data || [];
        this.loading = false;
      },
      error: () => {
        this.questoes = [];
        this.loading = false;
      },
    });
  }

  detalharQuestao(id: number) {
    this.router.navigate(['/questoes', id]);
  }
}
