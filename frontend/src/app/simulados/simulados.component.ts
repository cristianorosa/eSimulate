import { Component, OnInit } from '@angular/core';
import { SimuladosService } from './simulados.service';

@Component({
  selector: 'app-simulados',
  templateUrl: './simulados.component.html',
  styleUrls: ['./simulados.component.scss'],
})
export class SimuladosComponent implements OnInit {
  simulados: any[] = [];
  loading = true;

  constructor(private simuladosService: SimuladosService) {}

  ngOnInit() {
    this.simuladosService.getSimulados().subscribe({
      next: (data) => {
        this.simulados = data;
        this.loading = false;
      },
      error: () => {
        this.simulados = [];
        this.loading = false;
      },
    });
  }
}
