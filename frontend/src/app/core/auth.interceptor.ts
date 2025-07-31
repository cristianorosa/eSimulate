import { HttpEvent, HttpHandlerFn, HttpRequest } from "@angular/common/http";
import { Observable } from "rxjs";
import { AuthService } from "./auth.service";
import { inject } from "@angular/core";

export const AuthInterceptor = (
  req: HttpRequest<any>,
  next: HttpHandlerFn
): Observable<HttpEvent<any>> => {
  const auth = inject(AuthService);
  const token = auth.getToken();

  console.log('AuthInterceptor: URL =', req.url);
  console.log('AuthInterceptor: Token =', token ? 'presente' : 'ausente');

  // Verificar se o token está válido antes de adicionar à requisição
  if (token && req.url.includes('/api/')) {
    console.log('AuthInterceptor: Requisição para API com token');
    
    if (!auth.isTokenValid()) {
      console.log('AuthInterceptor: Token inválido, fazendo logout');
      auth.logout();
      return next(req);
    }
    
    console.log('AuthInterceptor: Adicionando token à requisição');
    const cloned = req.clone({
      setHeaders: { Authorization: `Bearer ${token}` },
    });
    return next(cloned);
  }
  
  console.log('AuthInterceptor: Requisição sem token ou não é API');
  return next(req);
};
