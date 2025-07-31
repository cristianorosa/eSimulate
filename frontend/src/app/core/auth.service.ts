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

  constructor(private http: HttpClient, private router: Router) {
    console.log('AuthService: Inicializando...');
    this.initializeUser();
  }

  private initializeUser() {
    const user = this.getUserFromToken();
    if (user) {
      this.user.set(user);
    } else {
      // Se não há usuário mas há token, limpar token inválido
      const token = localStorage.getItem(this.tokenKey);
      if (token) {
        localStorage.removeItem(this.tokenKey);
        localStorage.removeItem('esimulate_user');
      }
    }
  }

  login(email: string, password: string): Observable<any> {
    console.log('AuthService: Tentando login para:', email);
    return this.http.post<any>(`${this.apiUrl}/auth/login`, { email, password }).pipe(
      tap(res => {
        console.log('AuthService: Resposta do login:', res);
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
          console.log('AuthService: Login bem-sucedido para:', email);
        } else {
          console.error('AuthService: Resposta sem token:', res);
        }
      })
    );
  }

  logout() {
    localStorage.removeItem(this.tokenKey);
    localStorage.removeItem('esimulate_user');
    this.user.set(null);
    this.router.navigate(['/login']);
  }

  isAuthenticated(): boolean {
    const token = localStorage.getItem(this.tokenKey);
    const user = this.user();
    
    // Se não há token, não está autenticado
    if (!token) {
      console.log('AuthService: Sem token de autenticação');
      return false;
    }
    
    // Se não há usuário no signal, tenta recuperar
    if (!user) {
      console.log('AuthService: Recuperando usuário do localStorage...');
      this.initializeUser();
      const recoveredUser = this.user();
      return !!recoveredUser;
    }
    
    // Verificar se o token não expirou
    if (!this.isTokenValid()) {
      console.log('AuthService: Token expirado, fazendo logout');
      this.logout();
      return false;
    }
    
    console.log('AuthService: Usuário autenticado:', user.email);
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
      console.error('Erro ao verificar token:', error);
      return false;
    }
  }

  // Método para recarregar o usuário do localStorage
  reloadUser() {
    this.initializeUser();
  }
}
