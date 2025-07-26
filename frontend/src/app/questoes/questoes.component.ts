import { Component, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import { QuestoesService } from './questoes.service';

@Component({
  selector: 'app-questoes',
  templateUrl: './questoes.component.html',
  styleUrls: ['./questoes.component.scss'],
})
export class QuestoesComponent implements OnInit {
  questoes: any[] = [];
  loading = true;

  constructor(private questoesService: QuestoesService, private router: Router) {}

  ngOnInit() {
    this.questoesService.getQuestoes().subscribe({
      next: (data) => {
        this.questoes = data;
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
