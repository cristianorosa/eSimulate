import { ComponentFixture, TestBed } from '@angular/core/testing';
import { LoginComponent } from './login.component';
import { AuthService } from '../../core/auth.service';
import { Router } from '@angular/router';
import { MatSnackBarModule } from '@angular/material/snack-bar';
import { FormsModule } from '@angular/forms';
import { of, throwError } from 'rxjs';

class MockAuthService {
  login(email: string, password: string) {
    if (email === 'user@test.com' && password === '123456') {
      return of({ token: 'fake-jwt' });
    }
    return throwError({ error: { message: 'Usuário ou senha inválidos' } });
  }
}

class MockRouter {
  navigate(path: string[]) {}
}

describe('LoginComponent', () => {
  let component: LoginComponent;
  let fixture: ComponentFixture<LoginComponent>;
  let authService: AuthService;
  let router: Router;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [LoginComponent],
      imports: [MatSnackBarModule, FormsModule],
      providers: [
        { provide: AuthService, useClass: MockAuthService },
        { provide: Router, useClass: MockRouter },
      ],
    }).compileComponents();
    fixture = TestBed.createComponent(LoginComponent);
    component = fixture.componentInstance;
    authService = TestBed.inject(AuthService);
    router = TestBed.inject(Router);
    fixture.detectChanges();
  });

  it('deve autenticar com credenciais válidas', () => {
    spyOn(authService, 'login').and.callThrough();
    spyOn(router, 'navigate');
    component.email = 'user@test.com';
    component.password = '123456';
    component.onSubmit();
    expect(authService.login).toHaveBeenCalledWith('user@test.com', '123456');
  });

  it('deve exibir erro com credenciais inválidas', () => {
    component.email = 'user@test.com';
    component.password = 'errada';
    component.onSubmit();
    expect(component.error).toBe('Usuário ou senha inválidos');
  });
});
