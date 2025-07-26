import { Component } from '@angular/core';
import { AuthService } from '../../core/auth.service';
import { LoadingService } from '../../core/loading.service';

@Component({
  selector: 'app-main-layout',
  templateUrl: './main-layout.component.html',
  styleUrls: ['./main-layout.component.scss'],
})
export class MainLayoutComponent {
  user$ = this.auth.user$;
  loading$ = this.loading.loading$;

  constructor(public auth: AuthService, public loading: LoadingService) {}

  logout() {
    this.auth.logout();
  }
}
