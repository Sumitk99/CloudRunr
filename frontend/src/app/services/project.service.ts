import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { firstValueFrom } from 'rxjs';
import { environment } from '../../environments/environment';
import { AuthService } from './auth.service';

export interface NewProjectRequest {
  git_url: string;
  framework: string;
  dist_folder: string;
  project_id: string;
  name: string;
  run_command?: string | null;
}

export interface NewProjectResponse {
  deployment_id: string;
}

export interface ProjectItem {
  project_id: string;
  name: string;
  git_url: string;
  framework: string;
}

export interface ProjectsResponse {
  projects: ProjectItem[];
}

export interface LogData {
  log_statement: string;
  time: string;
}

export interface LogsResponse {
  data: LogData[];
}

export interface ProjectDetail {
  user_id: string;
  git_url: string;
  framework: string;
  dist_folder: string;
  project_id: string;
  name: string;
  run_command: string;
  subdomain: string;
}

export interface DeploymentResponse {
  status: string;
  deployment_id: string;
  url: string;
}

@Injectable({ providedIn: 'root' })
export class ProjectService {
  constructor(private http: HttpClient, private auth: AuthService) {}

  private buildTokenHeaders(): HttpHeaders {
    const token = this.auth.getToken();
    return new HttpHeaders({
      'Content-Type': 'application/json',
      token: token ?? ''
    });
  }

  async createNewProject(payload: NewProjectRequest): Promise<NewProjectResponse> {
    const url = `${environment.backendUrl}/project`;
    const headers = this.buildTokenHeaders();
    const res$ = this.http.post<NewProjectResponse>(url, payload, { headers });
    return await firstValueFrom(res$);
  }

  async fetchProjects(): Promise<ProjectsResponse> {
    const url = `${environment.backendUrl}/projects`;
    const headers = this.buildTokenHeaders();
    const res$ = this.http.get<ProjectsResponse>(url, { headers });
    return await firstValueFrom(res$);
  }

  async fetchDeploymentLogs(deployId: string, offset: number): Promise<LogsResponse> {
    const url = `${environment.backendUrl}/logs/${deployId}/${offset}`;
    const headers = this.buildTokenHeaders();
    const res$ = this.http.get<LogsResponse>(url, { headers });
    return await firstValueFrom(res$);
  }

  async getProjectDetail(projectId: string): Promise<ProjectDetail> {
    const url = `${environment.backendUrl}/project/detail/${projectId}`;
    const headers = this.buildTokenHeaders();
    const res$ = this.http.get<ProjectDetail>(url, { headers });
    return await firstValueFrom(res$);
  }

  async deployProject(projectId: string): Promise<DeploymentResponse> {
    const url = `${environment.backendUrl}/deploy/${projectId}`;
    const headers = this.buildTokenHeaders();
    const res$ = this.http.post<DeploymentResponse>(url, {}, { headers });
    return await firstValueFrom(res$);
  }
}
