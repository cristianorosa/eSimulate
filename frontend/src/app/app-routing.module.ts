import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { MainLayoutComponent } from './layout/main-layout/main-layout.component';
import { LoginComponent } from './auth/login/login.component';
import { RegisterComponent } from './auth/register/register.component';
import { DashboardComponent } from './dashboard/dashboard.component';
import { SimuladosComponent } from './simulados/simulados.component';
import { SimuladoExecucaoComponent } from './simulados/simulado-execucao/simulado-execucao.component';
import { QuestoesComponent } from './questoes/questoes.component';
import { PerfilComponent } from './perfil/perfil.component';
import { AuthGuard } from './core/auth.guard';
import { QuestaoDetalheComponent } from './questoes/questao-detalhe/questao-detalhe.component';
import { QuestaoFormComponent } from './questoes/questao-form/questao-form.component';

const routes: Routes = [
  { path: 'login', component: LoginComponent },
  { path: 'register', component: RegisterComponent },
  {
    path: '',
    component: MainLayoutComponent,
    canActivate: [AuthGuard],
    children: [
      { path: '', redirectTo: 'dashboard', pathMatch: 'full' },
      { path: 'dashboard', component: DashboardComponent },
      { path: 'simulados', component: SimuladosComponent },
      { path: 'simulados/:id/executar', component: SimuladoExecucaoComponent },
      { path: 'questoes', component: QuestoesComponent },
      { path: 'questoes/new', component: QuestaoFormComponent },
      { path: 'questoes/:id', component: QuestaoDetalheComponent },
      { path: 'perfil', component: PerfilComponent },
    ],
  },
  { path: '**', redirectTo: '' },
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule],
})
export class AppRoutingModule {}
