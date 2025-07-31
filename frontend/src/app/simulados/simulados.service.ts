import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';

@Injectable({ providedIn: 'root' })
export class SimuladosService {
  private apiUrl = '/api';

  constructor(private http: HttpClient) {}

  getSimulados(): Observable<any[]> {
    return this.http.get<any[]>(`${this.apiUrl}/quizzes`);
  }

  getSimuladoById(id: number): Observable<any> {
    return this.http.get<any>(`${this.apiUrl}/quizzes/detail?id=${id}`);
  }

  createSimulado(data: any): Observable<any> {
    return this.http.post(`${this.apiUrl}/quizzes/create`, data);
  }

  submitRespostas(quizId: number, respostas: any[]): Observable<any> {
    return this.http.post(`${this.apiUrl}/quizzes/answer`, { quiz_id: quizId, respostas });
  }
}
