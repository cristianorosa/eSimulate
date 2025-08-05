import { Injectable } from '@angular/core';
import { AuthService } from './auth.service';
import { Router } from '@angular/router';
import { interval, Subscription } from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class AuthMonitorService {
  private authCheckSubscription?: Subscription;
  private readonly CHECK_INTERVAL = 60000; // Verificar a cada 60 segundos (era 30)
  private lastAuthState = false;
  private isInitialized = false;

  constructor(
    private authService: AuthService,
    private router: Router
  ) {
    // Inicializar após um pequeno delay para garantir que o AuthService está pronto
    setTimeout(() => {
      this.initialize();
    }, 100);
  }

  private initialize() {
    if (this.isInitialized) return;
    
    this.isInitialized = true;
    this.lastAuthState = this.authService.isAuthenticated();
    this.startAuthMonitoring();
  }

  private startAuthMonitoring() {
    if (this.authCheckSubscription) {
      this.authCheckSubscription.unsubscribe();
    }
    
    this.authCheckSubscription = interval(this.CHECK_INTERVAL).subscribe(() => {
      this.checkAuthStatus();
    });
  }

  private stopAuthMonitoring() {
    if (this.authCheckSubscription) {
      this.authCheckSubscription.unsubscribe();
      this.authCheckSubscription = undefined;
    }
  }

  private checkAuthStatus() {
    const currentAuthState = this.authService.isAuthenticated();
    
    // Só fazer algo se o estado mudou
    if (currentAuthState !== this.lastAuthState) {
      this.lastAuthState = currentAuthState;
      
      // Se não está autenticado mas há token válido, tentar recarregar
      if (!currentAuthState && this.authService.isTokenValid()) {
        this.authService.reloadUser();
      }
      
      // Só fazer logout se não há token válido E não está autenticado
      if (!currentAuthState && !this.authService.isTokenValid()) {
        this.authService.logout();
      }
    }
  }

  // Método para verificar autenticação manualmente
  checkAuthManually() {
    this.checkAuthStatus();
  }

  // Método para parar o monitoramento (chamado quando o app é destruído)
  destroy() {
    this.stopAuthMonitoring();
  }
} 