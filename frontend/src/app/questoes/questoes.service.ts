import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';

@Injectable({ providedIn: 'root' })
export class QuestoesService {
  private apiUrl = 'http://localhost:8080';

  constructor(private http: HttpClient) {}

  getQuestoes(): Observable<any[]> {
    return this.http.get<any[]>(`${this.apiUrl}/questions`);
  }

  getQuestaoById(id: number): Observable<any> {
    return this.http.get<any>(`${this.apiUrl}/questions/detail?id=${id}`);
  }

  createQuestao(data: any): Observable<any> {
    return this.http.post(`${this.apiUrl}/questions/create`, data);
  }
}
