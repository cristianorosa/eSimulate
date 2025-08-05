import { Component, OnInit, OnDestroy } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule } from '@angular/router';
import { MatSidenavModule } from '@angular/material/sidenav';
import { MatToolbarModule } from '@angular/material/toolbar';
import { MatListModule } from '@angular/material/list';
import { MatIconModule } from '@angular/material/icon';
import { MatButtonModule } from '@angular/material/button';
import { MatMenuModule } from '@angular/material/menu';
import { MatDividerModule } from '@angular/material/divider';
import { MatProgressBarModule } from '@angular/material/progress-bar';
import { AuthService } from '../../core/auth.service';
import { LoadingService } from '../../core/loading.service';

@Component({
  selector: 'app-main-layout',
  standalone: true,
  imports: [
    CommonModule,
    RouterModule,
    MatToolbarModule,
    MatSidenavModule,
    MatListModule,
    MatIconModule,
    MatButtonModule,
    MatMenuModule,
    MatDividerModule,
    MatProgressBarModule
  ],
  templateUrl: './main-layout.component.html',
  styleUrls: ['./main-layout.component.scss']
})
export class MainLayoutComponent implements OnInit, OnDestroy {
  sidebarOpen = true;
  sidebarMini = false;
  private resizeListener: (() => void) | null = null;
  private readonly MOBILE_BREAKPOINT = 768;
  private readonly TABLET_BREAKPOINT = 1024;

  constructor(
    public auth: AuthService,
    public loadingService: LoadingService
  ) {}

  ngOnInit() {
    // Inicializar estado do sidebar baseado no tamanho da tela
    this.adjustSidebarForScreenSize();
    
    // Recarregar usuário se necessário
    if (!this.auth.user() && this.auth.isTokenValid()) {
      this.auth.reloadUser();
    }
    
    // Adicionar listener de resize
    this.resizeListener = () => {
      this.adjustSidebarForScreenSize();
    };
    
    window.addEventListener('resize', this.resizeListener);
  }

  ngOnDestroy() {
    // Remover listener de resize
    if (this.resizeListener) {
      window.removeEventListener('resize', this.resizeListener);
    }
  }

  private adjustSidebarForScreenSize() {
    const width = window.innerWidth;
    
    if (width <= this.MOBILE_BREAKPOINT) {
      // Mobile: fechar completamente
      this.sidebarOpen = false;
      this.sidebarMini = false;
    } else if (width <= this.TABLET_BREAKPOINT) {
      // Tablet: modo mini (apenas ícones)
      this.sidebarOpen = true;
      this.sidebarMini = true;
    } else {
      // Desktop: sidebar completo
      this.sidebarOpen = true;
      this.sidebarMini = false;
    }
  }

  toggleSidebar() {
    if (window.innerWidth <= this.MOBILE_BREAKPOINT) {
      // Em mobile, alternar entre fechado e aberto
      this.sidebarOpen = !this.sidebarOpen;
    } else if (window.innerWidth <= this.TABLET_BREAKPOINT) {
      // Em tablet, alternar entre mini e completo
      this.sidebarMini = !this.sidebarMini;
    } else {
      // Em desktop, alternar entre aberto e fechado
      this.sidebarOpen = !this.sidebarOpen;
    }
  }

  isAdminOrRedator(): boolean {
    const user = this.auth.user();
    
    if (!user || !user.role_id) {
      return false;
    }
    
    const hasPermission = user.role_id === 3 || user.role_id === 2;
    return hasPermission;
  }

  logout() {
    this.auth.logout();
  }
}

