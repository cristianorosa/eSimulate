import { Component } from '@angular/core';
import { Router } from '@angular/router';
import { QuestoesService } from '../questoes.service';
import { MatSnackBar } from '@angular/material/snack-bar';
import { LoadingService } from '../../core/loading.service';

@Component({
  selector: 'app-questao-form',
  templateUrl: './questao-form.component.html',
  styleUrls: ['./questao-form.component.scss'],
})
export class QuestaoFormComponent {
  statement = '';
  theme_id: number | null = null;
  explanation = '';
  options = [
    { text: '', is_correct: false, explanation: '' },
    { text: '', is_correct: false, explanation: '' },
  ];
  loading = false;
  error = '';
  success = false;

  constructor(
    private questoesService: QuestoesService,
    private router: Router,
    private snackBar: MatSnackBar,
    private loadingService: LoadingService
  ) {}

  addOption() {
    this.options.push({ text: '', is_correct: false, explanation: '' });
  }

  removeOption(i: number) {
    if (this.options.length > 2) {
      this.options.splice(i, 1);
    }
  }

  onSubmit() {
    this.loading = true;
    this.error = '';
    this.success = false;
    this.loadingService.show();
    const data = {
      statement: this.statement,
      theme_id: this.theme_id,
      explanation: this.explanation,
      options: this.options,
      created_by: 1, // Ajustar para pegar o usuário autenticado
    };
    this.questoesService.createQuestao(data).subscribe({
      next: () => {
        this.loading = false;
        this.success = true;
        this.loadingService.hide();
        this.snackBar.open('Questão cadastrada com sucesso!', 'Fechar', { duration: 2500 });
        setTimeout(() => this.router.navigate(['/questoes']), 1500);
      },
      error: (err) => {
        this.loading = false;
        this.loadingService.hide();
        this.error = err.error?.message || 'Erro ao cadastrar questão';
        this.snackBar.open(this.error, 'Fechar', { duration: 3000 });
      },
    });
  }
}
