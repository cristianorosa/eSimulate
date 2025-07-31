import { Routes } from '@angular/router';
import { MainLayoutComponent } from './layout/main-layout/main-layout.component';
import { LoginComponent } from './auth/login/login.component';
import { RegisterComponent } from './auth/register/register.component';
import { AuthGuard } from './core/auth.guard';
import { AdminGuard } from './core/admin.guard';
import { AreasComponent } from './admin/areas/areas.component';
import { ExamsComponent } from './admin/exams/exams.component';
import { DomainsComponent } from './admin/domains/domains.component';
import { QuestionsComponent } from './admin/questions/questions.component';

export const routes: Routes = [
  { path: 'login', component: LoginComponent },
  { path: 'register', component: RegisterComponent },
  { path: '', redirectTo: 'dashboard', pathMatch: 'full' },
  {
    path: 'dashboard',
    component: MainLayoutComponent,
    canActivate: [AuthGuard],
    children: [
      { path: '', loadComponent: () => import('./dashboard/dashboard.component').then(m => m.DashboardComponent) },
      { path: 'simulados', loadComponent: () => import('./simulados/simulados.component').then(m => m.SimuladosComponent) },
      { path: 'simulados/:id/executar', loadComponent: () => import('./simulados/simulado-execucao/simulado-execucao.component').then(m => m.SimuladoExecucaoComponent) },
      { path: 'questoes', loadComponent: () => import('./questoes/questoes.component').then(m => m.QuestoesComponent) },
      { path: 'questoes/new', loadComponent: () => import('./questoes/questao-form/questao-form.component').then(m => m.QuestaoFormComponent) },
      { path: 'questoes/:id', loadComponent: () => import('./questoes/questao-detalhe/questao-detalhe.component').then(m => m.QuestaoDetalheComponent) },
      { path: 'perfil', loadComponent: () => import('./perfil/perfil.component').then(m => m.PerfilComponent) },
      { path: 'historico', loadComponent: () => import('./historico/historico.component').then(m => m.HistoricoComponent) },
    ],
  },

  // Rotas administrativas separadas para evitar conflitos
  {
    path: 'admin',
    component: MainLayoutComponent,
    canActivate: [AuthGuard, AdminGuard],
    children: [
      { 
        path: 'areas', 
        component: AreasComponent
      },
      { 
        path: 'exams', 
        component: ExamsComponent
      },
      { 
        path: 'domains', 
        component: DomainsComponent
      },
      { 
        path: 'questions', 
        component: QuestionsComponent
      },
    ],
  },
  { path: '**', redirectTo: 'login' },
]; 