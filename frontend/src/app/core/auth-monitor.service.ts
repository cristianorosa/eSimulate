import { Injectable, OnInit, OnDestroy } from '@angular/core';
import { AuthService } from './auth.service';
import { Router } from '@angular/router';
import { interval, Subscription } from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class AuthMonitorService implements OnInit, OnDestroy {
  private authCheckSubscription?: Subscription;
  private readonly CHECK_INTERVAL = 30000; // Verificar a cada 30 segundos

  constructor(
    private authService: AuthService,
    private router: Router
  ) {
    this.startAuthMonitoring();
  }

  ngOnInit() {
    this.startAuthMonitoring();
  }

  ngOnDestroy() {
    this.stopAuthMonitoring();
  }

  private startAuthMonitoring() {
    this.authCheckSubscription = interval(this.CHECK_INTERVAL).subscribe(() => {
      this.checkAuthStatus();
    });
  }

  private stopAuthMonitoring() {
    if (this.authCheckSubscription) {
      this.authCheckSubscription.unsubscribe();
    }
  }

  private checkAuthStatus() {
    const isAuthenticated = this.authService.isAuthenticated();
    
    // Se não está autenticado mas há token válido, tentar recarregar
    if (!isAuthenticated && this.authService.isTokenValid()) {
      console.log('AuthMonitor: Token válido mas usuário não autenticado, recarregando...');
      this.authService.reloadUser();
    }
    
    // Só fazer logout se não há token válido E não está autenticado
    if (!isAuthenticated && !this.authService.isTokenValid()) {
      console.log('AuthMonitor: Sem autenticação e sem token válido, fazendo logout');
      this.authService.logout();
    }
  }

  // Método para verificar autenticação manualmente
  checkAuthManually() {
    this.checkAuthStatus();
  }
} 