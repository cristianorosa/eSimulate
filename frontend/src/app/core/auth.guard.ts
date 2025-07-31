import { CanActivateFn, Router } from '@angular/router';
import { inject } from '@angular/core';
import { AuthService } from './auth.service';

export const AuthGuard: CanActivateFn = (route, state) => {
  const auth = inject(AuthService);
  const router = inject(Router);

  console.log('AuthGuard: Verificando acesso para:', state.url);
  console.log('AuthGuard: Token presente:', !!auth.getToken());
  console.log('AuthGuard: Usuário no signal:', !!auth.user());
  
  // Verificar se está autenticado
  if (auth.isAuthenticated()) {
    console.log('AuthGuard: Acesso permitido');
    return true;
  }
  
  console.log('AuthGuard: Usuário não autenticado, redirecionando para login');
  
  // Verificar se há token mas usuário não carregado
  if (auth.getToken() && !auth.user()) {
    console.log('AuthGuard: Token presente mas usuário não carregado, tentando recarregar...');
    auth.reloadUser();
    
    if (auth.isAuthenticated()) {
      console.log('AuthGuard: Usuário recarregado com sucesso');
      return true;
    }
  }
  
  // Usa navigateByUrl em vez de parseUrl para garantir a navegação
  router.navigateByUrl('/login');
  return false;
};
