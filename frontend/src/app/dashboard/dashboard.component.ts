import { Component, OnInit } from '@angular/core';
import { DashboardService } from './dashboard.service';
import { AuthService } from '../core/auth.service';

@Component({
  selector: 'app-dashboard',
  templateUrl: './dashboard.component.html',
  styleUrls: ['./dashboard.component.scss'],
})
export class DashboardComponent implements OnInit {
  performance: any = null;
  history: any[] = [];
  loading = true;

  constructor(private dashboard: DashboardService, private auth: AuthService) {}

  ngOnInit() {
    const user = this.auth.userSubject.value;
    if (user && user.id) {
      this.dashboard.getPerformance(user.id).subscribe({
        next: (perf) => (this.performance = perf),
        error: () => (this.performance = null),
      });
      this.dashboard.getHistory(user.id).subscribe({
        next: (hist: any[]) => (this.history = hist),
        error: () => (this.history = []),
        complete: () => (this.loading = false),
      });
    } else {
      this.loading = false;
    }
  }
}
