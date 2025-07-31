import { Component } from '@angular/core';
import { Router } from '@angular/router';
import { AuthService } from '../../core/auth.service';
import { MatSnackBar } from '@angular/material/snack-bar';
import { LoadingService } from '../../core/loading.service';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { MatCardModule } from '@angular/material/card';
import { MatInputModule } from '@angular/material/input';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatSnackBarModule } from '@angular/material/snack-bar';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { RouterModule } from '@angular/router';

@Component({
  selector: 'app-login',
  templateUrl: './login.component.html',
  styleUrls: ['./login.component.scss'],
  standalone: true,
  imports: [
    CommonModule,
    FormsModule,
    MatCardModule,
    MatInputModule,
    MatFormFieldModule,
    MatButtonModule,
    MatIconModule,
    MatSnackBarModule,
    MatProgressSpinnerModule,
    RouterModule,
  ]
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
      next: (response) => {
        this.loading = false;
        this.loadingService.hide();
        this.snackBar.open('Login realizado com sucesso!', 'Fechar', { duration: 2000 });
        
        // Aguarda um pouco antes de navegar para garantir que o token foi salvo
        setTimeout(() => {
          this.router.navigate(['/dashboard']);
        }, 100);
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
