import { Injectable } from '@angular/core';
import { CanActivate, Router } from '@angular/router';
import { AuthService } from './auth.service';

@Injectable({
  providedIn: 'root'
})
export class AdminGuard implements CanActivate {
  constructor(
    private authService: AuthService,
    private router: Router
  ) {}

  canActivate(): boolean {
    console.log('AdminGuard: Verificando permissões administrativas...');
    
    // Primeiro verificar se está autenticado
    if (!this.authService.isAuthenticated()) {
      console.log('AdminGuard: Usuário não autenticado');
      this.router.navigateByUrl('/login');
      return false;
    }
    
    const user = this.authService.user();
    console.log('AdminGuard: Usuário atual:', user.email, 'role_id:', user.role_id);
    
    // Verificar se o usuário existe e tem role_id
    if (!user || !user.role_id) {
      console.log('AdminGuard: Usuário inválido ou sem role_id');
      this.authService.logout();
      this.router.navigateByUrl('/login');
      return false;
    }

    // Verificar se o usuário tem role de admin ou redator
    // role_id: 2 = redator, 3 = admin
    if (user.role_id === 2 || user.role_id === 3) {
      console.log('AdminGuard: Permissão concedida');
      return true;
    }

    // Se não tem permissão, redireciona para dashboard
    console.log('AdminGuard: Sem permissão administrativa');
    this.router.navigateByUrl('/dashboard');
    return false;
  }
} 