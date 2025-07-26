import { ComponentFixture, TestBed } from '@angular/core/testing';
import { RegisterComponent } from './register.component';
import { AuthService } from '../../core/auth.service';
import { HttpClientTestingModule, HttpTestingController } from '@angular/common/http/testing';
import { Router } from '@angular/router';
import { MatSnackBarModule } from '@angular/material/snack-bar';
import { FormsModule } from '@angular/forms';

class MockRouter {
  navigate(path: string[]) {}
}

describe('RegisterComponent', () => {
  let component: RegisterComponent;
  let fixture: ComponentFixture<RegisterComponent>;
  let httpMock: HttpTestingController;
  let router: Router;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [RegisterComponent],
      imports: [HttpClientTestingModule, MatSnackBarModule, FormsModule],
      providers: [
        { provide: Router, useClass: MockRouter },
        AuthService,
      ],
    }).compileComponents();
    fixture = TestBed.createComponent(RegisterComponent);
    component = fixture.componentInstance;
    httpMock = TestBed.inject(HttpTestingController);
    router = TestBed.inject(Router);
    fixture.detectChanges();
  });

  it('deve cadastrar usuário com dados válidos', () => {
    component.name = 'Novo Usuário';
    component.email = 'novo@teste.com';
    component.password = '123456';
    component.onSubmit();
    const req = httpMock.expectOne(r => r.url.includes('/auth/register'));
    expect(req.request.method).toBe('POST');
    req.flush({});
    expect(component.success).toBeTrue();
  });

  it('deve exibir erro ao falhar cadastro', () => {
    component.name = 'Novo Usuário';
    component.email = 'novo@teste.com';
    component.password = '123456';
    component.onSubmit();
    const req = httpMock.expectOne(r => r.url.includes('/auth/register'));
    req.flush({ message: 'Erro ao cadastrar' }, { status: 400, statusText: 'Bad Request' });
    expect(component.error).toBe('Erro ao cadastrar');
  });
});
