import { Component, OnInit } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { QuestoesService } from '../questoes.service';
import { CommonModule } from '@angular/common';
import { MatCardModule } from '@angular/material/card';
import { MatListModule } from '@angular/material/list';
import { MatIconModule } from '@angular/material/icon';
import { MatDividerModule } from '@angular/material/divider';
import { MatButtonModule } from '@angular/material/button';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';

@Component({
  selector: 'app-questao-detalhe',
  templateUrl: './questao-detalhe.component.html',
  styleUrls: ['./questao-detalhe.component.scss'],
  standalone: true,
  imports: [
    CommonModule,
    MatCardModule,
    MatListModule,
    MatIconModule,
    MatDividerModule,
    MatButtonModule,
    MatProgressSpinnerModule,
  ]
})
export class QuestaoDetalheComponent implements OnInit {
  questao: any = null;
  loading = true;
  String = String; // Para usar no template

  constructor(private route: ActivatedRoute, private questoesService: QuestoesService) {}

  ngOnInit() {
    const id = this.route.snapshot.paramMap.get('id');
    if (id) {
      this.questoesService.getQuestaoById(+id).subscribe({
        next: (data) => {
          this.questao = data;
          this.loading = false;
        },
        error: () => {
          this.questao = null;
          this.loading = false;
        },
      });
    } else {
      this.loading = false;
    }
  }
}
