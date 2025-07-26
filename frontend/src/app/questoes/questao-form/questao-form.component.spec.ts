import { ComponentFixture, TestBed } from '@angular/core/testing';
import { QuestaoFormComponent } from './questao-form.component';
import { QuestoesService } from '../questoes.service';
import { Router } from '@angular/router';
import { MatSnackBarModule } from '@angular/material/snack-bar';
import { FormsModule } from '@angular/forms';
import { of, throwError } from 'rxjs';
import { LoadingService } from '../../core/loading.service';

class MockQuestoesService {
  createQuestao(data: any) {
    if (data.statement && data.options.length >= 2) {
      return of({});
    }
    return throwError({ error: { message: 'Erro ao cadastrar questão' } });
  }
}
class MockRouter {
  navigate(path: string[]) {}
}
class MockLoadingService {
  show() {}
  hide() {}
}

describe('QuestaoFormComponent', () => {
  let component: QuestaoFormComponent;
  let fixture: ComponentFixture<QuestaoFormComponent>;
  let questoesService: QuestoesService;
  let router: Router;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [QuestaoFormComponent],
      imports: [MatSnackBarModule, FormsModule],
      providers: [
        { provide: QuestoesService, useClass: MockQuestoesService },
        { provide: Router, useClass: MockRouter },
        { provide: LoadingService, useClass: MockLoadingService },
      ],
    }).compileComponents();
    fixture = TestBed.createComponent(QuestaoFormComponent);
    component = fixture.componentInstance;
    questoesService = TestBed.inject(QuestoesService);
    router = TestBed.inject(Router);
    fixture.detectChanges();
  });

  it('deve cadastrar questão com dados válidos', () => {
    component.statement = 'Enunciado de teste';
    component.theme_id = 1;
    component.options = [
      { text: 'A', is_correct: true, explanation: '' },
      { text: 'B', is_correct: false, explanation: '' },
    ];
    component.onSubmit();
    expect(component.success).toBeTrue();
  });

  it('deve exibir erro ao falhar cadastro', () => {
    component.statement = '';
    component.theme_id = 1;
    component.options = [
      { text: '', is_correct: false, explanation: '' },
      { text: '', is_correct: false, explanation: '' },
    ];
    component.onSubmit();
    expect(component.error).toBe('Erro ao cadastrar questão');
  });
});
