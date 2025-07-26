import { Component, OnInit } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { SimuladosService } from '../simulados.service';

@Component({
  selector: 'app-simulado-execucao',
  templateUrl: './simulado-execucao.component.html',
  styleUrls: ['./simulado-execucao.component.scss'],
})
export class SimuladoExecucaoComponent implements OnInit {
  simulado: any = null;
  questoes: any[] = [];
  respostas: number[] = [];
  atual = 0;
  loading = true;
  finalizado = false;
  resultado: any = null;

  constructor(private route: ActivatedRoute, private simuladosService: SimuladosService) {}

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
}
