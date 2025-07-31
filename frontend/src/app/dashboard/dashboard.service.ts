import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';

@Injectable({ providedIn: 'root' })
export class DashboardService {
  private apiUrl = '/api';

  constructor(private http: HttpClient) {}

  getPerformance(userId: number): Observable<any> {
    console.log('DashboardService: Fazendo requisição para performance com userId:', userId);
    console.log('DashboardService: URL completa:', `${this.apiUrl}/performance?user_id=${userId}`);
    return this.http.get(`${this.apiUrl}/performance?user_id=${userId}`);
  }

  getHistory(userId: number): Observable<any> {
    return this.http.get(`${this.apiUrl}/history?user_id=${userId}`);
  }
}
