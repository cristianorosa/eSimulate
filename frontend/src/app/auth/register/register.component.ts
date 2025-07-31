import { Component } from '@angular/core';
import { Router } from '@angular/router';
import { AuthService } from '../../core/auth.service';
import { HttpClient } from '@angular/common/http';
import { MatSnackBar } from '@angular/material/snack-bar';
import { LoadingService } from '../../core/loading.service';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { MatCardModule } from '@angular/material/card';
import { MatInputModule } from '@angular/material/input';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { RouterModule } from '@angular/router';

@Component({
  selector: 'app-register',
  templateUrl: './register.component.html',
  styleUrls: ['./register.component.scss'],
  standalone: true,
  imports: [
    CommonModule,
    FormsModule,
    MatCardModule,
    MatInputModule,
    MatFormFieldModule,
    MatButtonModule,
    MatIconModule,
    MatProgressSpinnerModule,
    RouterModule,
  ]
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
