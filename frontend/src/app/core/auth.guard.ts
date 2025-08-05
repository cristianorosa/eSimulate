import { CanActivateFn, Router } from '@angular/router';
import { inject } from '@angular/core';
import { AuthService } from './auth.service';

export const AuthGuard: CanActivateFn = (route, state) => {
  const auth = inject(AuthService);
  const router = inject(Router);
  
  // Verificar se está autenticado
  if (auth.isAuthenticated()) {
    return true;
  }
  
  // Verificar se há token mas usuário não carregado
  if (auth.getToken() && !auth.user()) {
    auth.reloadUser();
    
    if (auth.isAuthenticated()) {
      return true;
    }
  }
  
  // Usa navigateByUrl em vez de parseUrl para garantir a navegação
  router.navigateByUrl('/login');
  return false;
};
