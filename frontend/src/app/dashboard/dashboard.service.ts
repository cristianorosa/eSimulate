import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';

@Injectable({ providedIn: 'root' })
export class DashboardService {
  private apiUrl = 'http://localhost:8080';

  constructor(private http: HttpClient) {}

  getPerformance(userId: number): Observable<any> {
    return this.http.get(`${this.apiUrl}/performance?user_id=${userId}`);
  }

  getHistory(userId: number): Observable<any> {
    return this.http.get(`${this.apiUrl}/history?user_id=${userId}`);
  }
}
