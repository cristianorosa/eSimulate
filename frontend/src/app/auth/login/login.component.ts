import { Component } from '@angular/core';
import { Router } from '@angular/router';
import { AuthService } from '../../core/auth.service';
import { MatSnackBar } from '@angular/material/snack-bar';
import { LoadingService } from '../../core/loading.service';

@Component({
  selector: 'app-login',
  templateUrl: './login.component.html',
  styleUrls: ['./login.component.scss'],
})
export class LoginComponent {
  email = '';
  password = '';
  loading = false;
  error = '';

  constructor(
    private auth: AuthService,
    private router: Router,
    private snackBar: MatSnackBar,
    private loadingService: LoadingService
  ) {}

  onSubmit() {
    this.loading = true;
    this.error = '';
    this.loadingService.show();
    this.auth.login(this.email, this.password).subscribe({
      next: () => {
        this.loading = false;
        this.loadingService.hide();
        this.snackBar.open('Login realizado com sucesso!', 'Fechar', { duration: 2000 });
        this.router.navigate(['/dashboard']);
      },
      error: (err) => {
        this.loading = false;
        this.loadingService.hide();
        this.error = err.error?.message || 'Usuário ou senha inválidos';
        this.snackBar.open(this.error, 'Fechar', { duration: 3000 });
      },
    });
  }
}
