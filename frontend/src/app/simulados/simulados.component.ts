import { Component, OnInit } from '@angular/core';
import { SimuladosService } from './simulados.service';
import { CommonModule } from '@angular/common';
import { MatCardModule } from '@angular/material/card';
import { MatListModule } from '@angular/material/list';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatProgressBarModule } from '@angular/material/progress-bar';
import { RouterModule } from '@angular/router';

@Component({
  selector: 'app-simulados',
  templateUrl: './simulados.component.html',
  styleUrls: ['./simulados.component.scss'],
  standalone: true,
  imports: [
    CommonModule,
    MatCardModule,
    MatListModule,
    MatButtonModule,
    MatIconModule,
    MatProgressBarModule,
    RouterModule,
  ]
})
export class SimuladosComponent implements OnInit {
  simulados: any[] = [];
  loading = true;

  constructor(private simuladosService: SimuladosService) {}

  ngOnInit() {
    this.simuladosService.getSimulados().subscribe({
      next: (data) => {
        this.simulados = data || [];
        this.loading = false;
      },
      error: () => {
        this.simulados = [];
        this.loading = false;
      },
    });
  }
}
