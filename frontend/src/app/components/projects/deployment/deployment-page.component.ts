import { Component, OnDestroy, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ActivatedRoute } from '@angular/router';
import { MatCardModule } from '@angular/material/card';
import { MatButtonModule } from '@angular/material/button';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { ProjectService, LogData } from '../../../services/project.service';

@Component({
  selector: 'app-deployment-page',
  standalone: true,
  imports: [CommonModule, MatCardModule, MatButtonModule, MatProgressSpinnerModule],
  templateUrl: './deployment-page.component.html',
  styleUrl: './deployment-page.component.scss'
})
export class DeploymentPageComponent implements OnInit, OnDestroy {
  projectId!: string;
  deploymentId!: string;

  logs: LogData[] = [];
  offset = 0;
  isLoading = false;
  autoRefreshHandle: any;

  constructor(private route: ActivatedRoute, private projects: ProjectService) {}

  ngOnInit(): void {
    this.projectId = this.route.snapshot.paramMap.get('projectId') || '';
    this.deploymentId = this.route.snapshot.paramMap.get('deploymentId') || '';
    this.fetchLogs();
    // Auto-refresh every 3 seconds
    this.autoRefreshHandle = setInterval(() => this.refreshLatest(), 3000);
  }

  ngOnDestroy(): void {
    if (this.autoRefreshHandle) clearInterval(this.autoRefreshHandle);
  }

  async fetchLogs(): Promise<void> {
    this.isLoading = true;
    try {
      const res = await this.projects.fetchDeploymentLogs(this.deploymentId, this.offset);
      // Server returns newest first; we want newest at top; keep as-is
      if (res.data && res.data.length) {
        this.logs = [...this.logs, ...res.data];
        this.offset += res.data.length;
      }
    } finally {
      this.isLoading = false;
    }
  }

  async refreshLatest(): Promise<void> {
    // Re-fetch from current offset to append new items if any
    try {
      const res = await this.projects.fetchDeploymentLogs(this.deploymentId, 1);
      if (res.data && res.data.length) {
        // Replace with latest snapshot (keeps most recent at top)
        this.logs = res.data;
      }
    } catch {}
  }
}
