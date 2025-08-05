import { Injectable, signal } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Router } from '@angular/router';
import { Observable } from 'rxjs';
import { tap } from 'rxjs/operators';

@Injectable({ providedIn: 'root' })
export class AuthService {
  private apiUrl = '/api';
  private tokenKey = 'esimulate_token';
  user = signal<any | null>(null);
  private authState = signal<boolean>(false);
  private lastAuthCheck = 0;
  private readonly AUTH_CACHE_DURATION = 5000; // 5 segundos de cache

  constructor(private http: HttpClient, private router: Router) {
    this.initializeUser();
  }

  private initializeUser() {
    const user = this.getUserFromToken();
    if (user) {
      this.user.set(user);
      this.authState.set(true);
    } else {
      // Se não há usuário mas há token, limpar token inválido
      const token = localStorage.getItem(this.tokenKey);
      if (token) {
        localStorage.removeItem(this.tokenKey);
        localStorage.removeItem('esimulate_user');
      }
      this.authState.set(false);
    }
  }

  login(email: string, password: string): Observable<any> {
    return this.http.post<any>(`${this.apiUrl}/auth/login`, { email, password }).pipe(
      tap(res => {
        if (res.token) {
          localStorage.setItem(this.tokenKey, res.token);
          // Criar objeto usuário com dados da resposta
          const userData = {
            id: res.id || 1,
            name: res.name || email,
            email: email,
            role_id: res.role_id || 1,
            token: res.token
          };
          // Salvar dados do usuário no localStorage
          localStorage.setItem('esimulate_user', JSON.stringify(userData));
          this.user.set(userData);
          this.authState.set(true);
        }
      })
    );
  }

  logout() {
    localStorage.removeItem(this.tokenKey);
    localStorage.removeItem('esimulate_user');
    this.user.set(null);
    this.authState.set(false);
    this.router.navigate(['/login']);
  }

  isAuthenticated(): boolean {
    const now = Date.now();
    
    // Usar cache se a verificação foi feita recentemente
    if (now - this.lastAuthCheck < this.AUTH_CACHE_DURATION) {
      return this.authState();
    }
    
    this.lastAuthCheck = now;
    const token = localStorage.getItem(this.tokenKey);
    const user = this.user();
    
    // Se não há token, não está autenticado
    if (!token) {
      this.authState.set(false);
      return false;
    }
    
    // Se não há usuário no signal, tenta recuperar
    if (!user) {
      const recoveredUser = this.getUserFromToken();
      if (recoveredUser) {
        this.user.set(recoveredUser);
        this.authState.set(true);
        return true;
      } else {
        this.authState.set(false);
        return false;
      }
    }
    
    // Verificar se o token não expirou
    if (!this.isTokenValid()) {
      this.logout();
      return false;
    }
    
    this.authState.set(true);
    return true;
  }

  getToken(): string | null {
    return localStorage.getItem(this.tokenKey);
  }

  getUserFromToken(): any {
    const token = this.getToken();
    if (token) {
      // Recuperar dados do usuário do localStorage se disponível
      const userData = localStorage.getItem('esimulate_user');
      if (userData) {
        try {
          const parsed = JSON.parse(userData);
          
          // Verificar se os dados do usuário são válidos
          if (parsed && parsed.id && parsed.email) {
            return parsed;
          }
        } catch (e) {
          console.error('Erro ao parsear dados do usuário:', e);
        }
      }
      
      // Se não há dados salvos válidos, tentar extrair do token
      try {
        const tokenData = JSON.parse(atob(token.split('.')[1]));
        const fallback = {
          id: tokenData.user_id || 1,
          name: tokenData.name || 'Usuário',
          email: tokenData.email || 'usuario@example.com',
          role_id: tokenData.role_id || 1,
          token: token
        };
        return fallback;
      } catch (error) {
        console.error('Erro ao extrair dados do token:', error);
      }
    }
    return null;
  }

  updateUser(userData: any) {
    this.user.set(userData);
  }

  // Método para verificar se o token está válido sem fazer logout
  isTokenValid(): boolean {
    const token = this.getToken();
    if (!token) {
      return false;
    }
    
    try {
      const tokenData = JSON.parse(atob(token.split('.')[1]));
      const currentTime = Math.floor(Date.now() / 1000);
      return tokenData.exp && tokenData.exp > currentTime;
    } catch (error) {
      return false;
    }
  }

  // Método para recarregar o usuário do localStorage
  reloadUser() {
    this.initializeUser();
  }

  // Método para forçar verificação de autenticação (ignora cache)
  forceAuthCheck(): boolean {
    this.lastAuthCheck = 0;
    return this.isAuthenticated();
  }
}
