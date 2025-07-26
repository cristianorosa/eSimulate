import { Component, OnInit } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { QuestoesService } from '../questoes.service';

@Component({
  selector: 'app-questao-detalhe',
  templateUrl: './questao-detalhe.component.html',
  styleUrls: ['./questao-detalhe.component.scss'],
})
export class QuestaoDetalheComponent implements OnInit {
  questao: any = null;
  loading = true;

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
