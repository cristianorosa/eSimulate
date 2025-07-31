import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { MatCardModule } from '@angular/material/card';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatDividerModule } from '@angular/material/divider';
import { MatSnackBar } from '@angular/material/snack-bar';
import { AuthService } from '../core/auth.service';

@Component({
  selector: 'app-perfil',
  standalone: true,
  imports: [
    CommonModule,
    FormsModule,
    MatCardModule,
    MatFormFieldModule,
    MatInputModule,
    MatButtonModule,
    MatIconModule,
    MatDividerModule
  ],
  templateUrl: './perfil.component.html',
  styleUrls: ['./perfil.component.scss']
})
export class PerfilComponent implements OnInit {
  user: any = null;
  editMode = false;
  loading = false;
  
  // Dados do formulário
  name = '';
  email = '';
  currentPassword = '';
  newPassword = '';
  confirmPassword = '';

  constructor(
    private auth: AuthService,
    private snackBar: MatSnackBar
  ) {}

  ngOnInit() {
    this.user = this.auth.user();
    if (this.user) {
      this.name = this.user.name;
      this.email = this.user.email;
    }
  }

  toggleEditMode() {
    this.editMode = !this.editMode;
    if (!this.editMode) {
      // Reset form
      this.name = this.user.name;
      this.email = this.user.email;
      this.currentPassword = '';
      this.newPassword = '';
      this.confirmPassword = '';
    }
  }

  updateProfile() {
    if (!this.name.trim() || !this.email.trim()) {
      this.snackBar.open('Nome e email são obrigatórios', 'Fechar', { duration: 3000 });
      return;
    }

    if (this.newPassword && this.newPassword !== this.confirmPassword) {
      this.snackBar.open('As senhas não coincidem', 'Fechar', { duration: 3000 });
      return;
    }

    this.loading = true;
    
    // Simular atualização (implementar chamada real para API)
    setTimeout(() => {
      this.user.name = this.name;
      this.user.email = this.email;
      this.auth.updateUser(this.user);
      this.editMode = false;
      this.loading = false;
      this.snackBar.open('Perfil atualizado com sucesso!', 'Fechar', { duration: 3000 });
    }, 1000);
  }

  getStats() {
    return {
      simuladosRealizados: 12,
      mediaAcertos: 78,
      tempoMedio: '45 min',
      ultimaAtividade: '2 dias atrás'
    };
  }
} 