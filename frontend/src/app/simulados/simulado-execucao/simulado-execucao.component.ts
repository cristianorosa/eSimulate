import { Component, OnInit } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { SimuladosService } from '../simulados.service';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { MatCardModule } from '@angular/material/card';
import { MatRadioModule } from '@angular/material/radio';
import { MatButtonModule } from '@angular/material/button';
import { MatDividerModule } from '@angular/material/divider';
import { MatIconModule } from '@angular/material/icon';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';

@Component({
  selector: 'app-simulado-execucao',
  templateUrl: './simulado-execucao.component.html',
  styleUrls: ['./simulado-execucao.component.scss'],
  standalone: true,
  imports: [
    CommonModule,
    FormsModule,
    MatCardModule,
    MatRadioModule,
    MatButtonModule,
    MatDividerModule,
    MatIconModule,
    MatProgressSpinnerModule,
  ]
})
export class SimuladoExecucaoComponent implements OnInit {
  simulado: any = null;
  questoes: any[] = [];
  respostas: number[] = [];
  atual = 0;
  loading = true;
  finalizado = false;
  resultado: any = null;
  String = String; // Para usar no template

  constructor(
    private route: ActivatedRoute, 
    private simuladosService: SimuladosService,
    private router: Router
  ) {}

  ngOnInit() {
    const id = this.route.snapshot.paramMap.get('id');
    if (id) {
      this.simuladosService.getSimuladoById(+id).subscribe({
        next: (data) => {
          this.simulado = data;
          this.questoes = data.questions || [];
          this.respostas = new Array(this.questoes.length).fill(null);
          this.loading = false;
        },
        error: () => {
          this.simulado = null;
          this.loading = false;
        },
      });
    } else {
      this.loading = false;
    }
  }

  responder(opcaoId: number) {
    this.respostas[this.atual] = opcaoId;
  }

  proxima() {
    if (this.atual < this.questoes.length - 1) {
      this.atual++;
    }
  }

  anterior() {
    if (this.atual > 0) {
      this.atual--;
    }
  }

  finalizar() {
    this.finalizado = true;
    this.simuladosService.submitRespostas(this.simulado.id, this.respostas).subscribe({
      next: (res) => {
        this.resultado = res;
      },
      error: () => {
        this.resultado = { erro: 'Erro ao processar resultado.' };
      },
    });
  }

  get allQuestionsAnswered(): boolean {
    return this.respostas.every(resposta => resposta !== undefined && resposta !== null);
  }
}
