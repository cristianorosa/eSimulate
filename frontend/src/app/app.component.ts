import { Component } from '@angular/core';
import { RouterModule } from '@angular/router';
import { AuthMonitorService } from './core/auth-monitor.service';

@Component({
  selector: 'app-root',
  template: '<router-outlet></router-outlet>',
  styleUrls: ['./app.component.scss'],
  standalone: true,
  imports: [RouterModule]
})
export class AppComponent {
  title = 'frontend';
  
  constructor(private authMonitor: AuthMonitorService) {}
}
