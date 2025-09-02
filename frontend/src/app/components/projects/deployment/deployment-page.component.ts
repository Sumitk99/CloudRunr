import { Component, OnDestroy, OnInit, AfterViewInit, ElementRef, ViewChild } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ActivatedRoute, Router } from '@angular/router';
import { MatCardModule } from '@angular/material/card';
import { MatButtonModule } from '@angular/material/button';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatIconModule } from '@angular/material/icon';
import { MatChipsModule } from '@angular/material/chips';
import { MatDividerModule } from '@angular/material/divider';
import { gsap } from 'gsap';
import { ProjectService, LogData, DeploymentResponse } from '../../../services/project.service';

@Component({
  selector: 'app-deployment-page',
  standalone: true,
  imports: [CommonModule, MatCardModule, MatButtonModule, MatProgressSpinnerModule, MatIconModule, MatChipsModule, MatDividerModule],
  templateUrl: './deployment-page.component.html',
  styleUrl: './deployment-page.component.scss'
})
export class DeploymentPageComponent implements OnInit, OnDestroy, AfterViewInit {
  @ViewChild('deploymentContainer') deploymentContainer!: ElementRef;
  @ViewChild('deploymentHeader') deploymentHeader!: ElementRef;
  @ViewChild('logsContainer') logsContainer!: ElementRef;

  projectId!: string;
  deploymentId!: string;
  deploymentUrl = '';
  deploymentStatus = 'QUEUED';

  logs: LogData[] = [];
  offset = 0;
  isLoading = false;
  autoRefreshHandle: any;

  constructor(
    private route: ActivatedRoute, 
    private router: Router,
    private projects: ProjectService
  ) {}

  ngOnInit(): void {
    this.projectId = this.route.snapshot.paramMap.get('projectId') || '';
    this.deploymentId = this.route.snapshot.paramMap.get('deploymentId') || '';
    
    // Set initial deployment URL and status from route state if available
    const navigation = this.router.getCurrentNavigation();
    if (navigation?.extras?.state) {
      this.deploymentUrl = navigation.extras.state['url'] || '';
      this.deploymentStatus = navigation.extras.state['status'] || 'QUEUED';
    } else {
      // Fallback: construct URL from projectId
      this.deploymentUrl = `${this.projectId}.cloudrunr.micro-scale.software`;
    }

    this.fetchLogs();
    // Auto-refresh every 3 seconds
    this.autoRefreshHandle = setInterval(() => this.refreshLatest(), 3000);
  }

  ngAfterViewInit(): void {
    this.animatePageEntrance();
  }

  private animatePageEntrance(): void {
    const tl = gsap.timeline();
    
    tl.fromTo(this.deploymentContainer.nativeElement,
      { opacity: 0, y: 30 },
      { opacity: 1, y: 0, duration: 0.8, ease: "power2.out" }
    )
    .fromTo(this.deploymentHeader.nativeElement,
      { opacity: 0, scale: 0.95 },
      { opacity: 1, scale: 1, duration: 0.6, ease: "back.out(1.7)" },
      "-=0.5"
    )
    .fromTo(this.logsContainer.nativeElement,
      { opacity: 0, y: 20 },
      { opacity: 1, y: 0, duration: 0.5, ease: "power2.out" },
      "-=0.3"
    );
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

  openDeploymentUrl(): void {
    if (this.deploymentUrl) {
      window.open(`http://${this.deploymentUrl}`, '_blank');
    }
  }

  goBack(): void {
    this.router.navigate(['/projects', this.projectId]);
  }

  getStatusColor(): string {
    switch (this.deploymentStatus.toLowerCase()) {
      case 'success':
      case 'deployed':
        return 'success';
      case 'failed':
      case 'error':
        return 'error';
      case 'running':
      case 'building':
        return 'primary';
      default:
        return 'default';
    }
  }

  getStatusIcon(): string {
    switch (this.deploymentStatus.toLowerCase()) {
      case 'success':
      case 'deployed':
        return 'check_circle';
      case 'failed':
      case 'error':
        return 'error';
      case 'running':
      case 'building':
        return 'sync';
      default:
        return 'schedule';
    }
  }

  trackByLogTime(index: number, log: LogData): string {
    return log.time;
  }
}
