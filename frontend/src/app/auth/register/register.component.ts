import { Component } from '@angular/core';
import { Router } from '@angular/router';
import { AuthService } from '../../core/auth.service';
import { HttpClient } from '@angular/common/http';
import { MatSnackBar } from '@angular/material/snack-bar';
import { LoadingService } from '../../core/loading.service';

@Component({
  selector: 'app-register',
  templateUrl: './register.component.html',
  styleUrls: ['./register.component.scss'],
})
export class RegisterComponent {
  name = '';
  email = '';
  password = '';
  loading = false;
  error = '';
  success = false;

  constructor(
    private http: HttpClient,
    private router: Router,
    private auth: AuthService,
    private snackBar: MatSnackBar,
    private loadingService: LoadingService
  ) {}

  onSubmit() {
    this.loading = true;
    this.error = '';
    this.success = false;
    this.loadingService.show();
    this.http.post<any>(`${this.auth['apiUrl']}/auth/register`, {
      name: this.name,
      email: this.email,
      password: this.password,
    }).subscribe({
      next: () => {
        this.loading = false;
        this.success = true;
        this.loadingService.hide();
        this.snackBar.open('Cadastro realizado com sucesso! Faça login.', 'Fechar', { duration: 2500 });
        setTimeout(() => this.router.navigate(['/login']), 1500);
      },
      error: (err) => {
        this.loading = false;
        this.loadingService.hide();
        this.error = err.error?.message || 'Erro ao cadastrar usuário';
        this.snackBar.open(this.error, 'Fechar', { duration: 3000 });
      },
    });
  }
}
