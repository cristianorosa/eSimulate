import { Injectable } from "@angular/core";
import { HttpClient } from "@angular/common/http";
import { Observable } from "rxjs";

export interface Pagination {
  page: number;
  page_size: number;
  total_items: number;
  total_pages: number;
}

export interface PaginatedResponse<T> {
  data: T[];
  pagination: Pagination;
}

export interface Area {
  id?: number;
  name: string;
  description?: string;
  created_at?: string;
}

export interface Exam {
  id?: number;
  title: string;
  description?: string;
  area_id: number;
  max_time_minutes: number;
  passing_score: number;
  questions_count?: number;
  is_active?: boolean;
  created_by?: number;
  created_at?: string;
  updated_at?: string;
}

export interface Topic {
  id?: number;
  exam_id: number;
  name: string;
  weight_percentage: number;
  order_index?: number;
  questions_count?: number;
  created_at?: string;
}

export interface Question {
  id?: number;
  exam_id?: number; // Para uso no frontend (filtros, exibição)
  topic_id: number;
  statement: string;
  problem: string;
  content_type: "text" | "code";
  explanation?: string;
  question_type: "objective" | "multiple_choice";
  difficulty_level?: number;
  created_by?: number;
  is_active?: boolean;
  created_at?: string;
  updated_at?: string;
  options?: Option[];
  // Campos adicionais que vêm do backend
  exam_title?: string;
  topic_name?: string;
}

export interface Option {
  id?: number;
  question_id?: number;
  text: string;
  is_correct: boolean;
  explanation?: string;
  order_index?: number;
}

@Injectable({
  providedIn: "root",
})
export class AdminService {
  private apiUrl = "/api";

  constructor(private http: HttpClient) {}

  // ===== ÁREAS =====
  getAreas(): Observable<Area[]> {
    return this.http.get<Area[]>(`${this.apiUrl}/areas`);
  }

  getAreasPaginated(
    page: number = 1,
    pageSize: number = 10
  ): Observable<PaginatedResponse<Area>> {
    return this.http.get<PaginatedResponse<Area>>(
      `${this.apiUrl}/areas/paginated?page=${page}&page_size=${pageSize}`
    );
  }

  getArea(id: number): Observable<Area> {
    return this.http.get<Area>(`${this.apiUrl}/areas/detail?id=${id}`);
  }

  createArea(area: Area): Observable<Area> {
    return this.http.post<Area>(`${this.apiUrl}/areas/create`, area);
  }

  updateArea(id: number, area: Area): Observable<any> {
    return this.http.put(`${this.apiUrl}/areas/update?id=${id}`, area);
  }

  deleteArea(id: number): Observable<any> {
    return this.http.delete(`${this.apiUrl}/areas/delete?id=${id}`);
  }

  // ===== EXAMES =====
  getExams(areaId?: number): Observable<Exam[]> {
    let url = `${this.apiUrl}/exams`;
    if (areaId) {
      url += `?area_id=${areaId}`;
    }
    return this.http.get<Exam[]>(url);
  }

  getExamsPaginated(
    page: number = 1,
    pageSize: number = 10,
    areaId?: number
  ): Observable<PaginatedResponse<Exam>> {
    let url = `${this.apiUrl}/exams/paginated?page=${page}&page_size=${pageSize}`;
    if (areaId) {
      url += `&area_id=${areaId}`;
    }
    return this.http.get<PaginatedResponse<Exam>>(url);
  }

  getExam(id: number): Observable<Exam> {
    return this.http.get<Exam>(`${this.apiUrl}/exams/detail?id=${id}`);
  }

  createExam(exam: Exam): Observable<Exam> {
    return this.http.post<Exam>(`${this.apiUrl}/exams/create`, exam);
  }

  updateExam(id: number, exam: Exam): Observable<any> {
    return this.http.put(`${this.apiUrl}/exams/update?id=${id}`, exam);
  }

  deleteExam(id: number): Observable<any> {
    return this.http.delete(`${this.apiUrl}/exams/delete?id=${id}`);
  }

  // ===== TÓPICOS =====
  getTopics(examId?: number): Observable<Topic[]> {
    let url = `${this.apiUrl}/topics`;
    if (examId) {
      url += `?exam_id=${examId}`;
    }
    return this.http.get<Topic[]>(url);
  }

  getTopicsPaginated(
    page: number = 1,
    pageSize: number = 10,
    examId?: number
  ): Observable<PaginatedResponse<Topic>> {
    let url = `${this.apiUrl}/topics/paginated?page=${page}&page_size=${pageSize}`;
    if (examId) {
      url += `&exam_id=${examId}`;
    }
    console.log("Chamando URL dos tópicos:", url);
    return this.http.get<PaginatedResponse<Topic>>(url);
  }

  getTopic(id: number): Observable<Topic> {
    return this.http.get<Topic>(`${this.apiUrl}/topics/detail?id=${id}`);
  }

  updateTopic(id: number, topic: Topic): Observable<any> {
    return this.http.put(`${this.apiUrl}/topics/update?id=${id}`, topic);
  }

  deleteTopic(id: number): Observable<any> {
    return this.http.delete(`${this.apiUrl}/topics/delete?id=${id}`);
  }

  // ===== QUESTÕES =====
  getQuestions(examId?: number, topicId?: number): Observable<Question[]> {
    let url = `${this.apiUrl}/questions`;
    const params = [];
    if (examId) params.push(`exam_id=${examId}`);
    if (topicId) params.push(`topic_id=${topicId}`);
    if (params.length > 0) {
      url += `?${params.join("&")}`;
    }
    return this.http.get<Question[]>(url);
  }

  getQuestionsPaginated(
    page: number = 1,
    pageSize: number = 10,
    examId?: number,
    topicId?: number
  ): Observable<PaginatedResponse<Question>> {
    let url = `${this.apiUrl}/questions/paginated?page=${page}&page_size=${pageSize}`;
    const params = [];
    if (examId) params.push(`exam_id=${examId}`);
    if (topicId) params.push(`topic_id=${topicId}`);
    if (params.length > 0) {
      url += `&${params.join("&")}`;
    }
    console.log("Chamando URL das questões:", url);
    return this.http.get<PaginatedResponse<Question>>(url);
  }

  getQuestion(id: number): Observable<Question> {
    return this.http.get<Question>(`${this.apiUrl}/questions/detail?id=${id}`);
  }

  createQuestion(question: Question): Observable<Question> {
    return this.http.post<Question>(
      `${this.apiUrl}/questions/create`,
      question
    );
  }

  updateQuestion(id: number, question: Question): Observable<any> {
    return this.http.put(`${this.apiUrl}/questions/update?id=${id}`, question);
  }

  deleteQuestion(id: number): Observable<any> {
    return this.http.delete(`${this.apiUrl}/questions/delete?id=${id}`);
  }

  // ===== OPÇÕES =====
  getOptions(questionId: number): Observable<Option[]> {
    return this.http.get<Option[]>(
      `${this.apiUrl}/options?question_id=${questionId}`
    );
  }

  createOption(option: Option): Observable<Option> {
    return this.http.post<Option>(`${this.apiUrl}/options/create`, option);
  }

  updateOption(id: number, option: Option): Observable<any> {
    return this.http.put(`${this.apiUrl}/options/update?id=${id}`, option);
  }

  deleteOption(id: number): Observable<any> {
    return this.http.delete(`${this.apiUrl}/options/delete?id=${id}`);
  }

  // Importação de questões
  importQuestions(data: any): Observable<any> {
    return this.http.post(`${this.apiUrl}/questions/import`, data);
  }

  // ===== TAGS =====

  getTags(): Observable<any[]> {
    return this.http.get<any[]>(`${this.apiUrl}/tags`);
  }

  searchTags(query: string, limit: number = 20): Observable<any> {
    return this.http.get(`${this.apiUrl}/tags/search?q=${encodeURIComponent(query)}&limit=${limit}`);
  }

  createTag(tag: { name: string }): Observable<any> {
    return this.http.post(`${this.apiUrl}/tags/create`, tag);
  }

  findOrCreateTag(tag: { name: string }): Observable<any> {
    return this.http.post(`${this.apiUrl}/tags/find-or-create`, tag);
  }

  associateQuestionTags(questionId: number, tagIds: number[]): Observable<any> {
    return this.http.put(`${this.apiUrl}/questions/${questionId}/tags`, { tag_ids: tagIds });
  }

  getQuestionTags(questionId: number): Observable<any> {
    return this.http.get(`${this.apiUrl}/questions/${questionId}/tags`);
  }

  createTopic(topic: { name: string; description?: string }): Observable<Topic> {
    return this.http.post<Topic>(`${this.apiUrl}/topics/create`, topic);
  }
}
