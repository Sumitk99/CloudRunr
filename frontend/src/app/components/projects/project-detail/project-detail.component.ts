import { Component, OnInit, AfterViewInit, ElementRef, ViewChild } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ActivatedRoute, Router } from '@angular/router';
import { MatCardModule } from '@angular/material/card';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatChipsModule } from '@angular/material/chips';
import { MatDividerModule } from '@angular/material/divider';
import { MatSnackBar, MatSnackBarModule } from '@angular/material/snack-bar';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { gsap } from 'gsap';
import { ProjectService, ProjectItem } from '../../../services/project.service';

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

@Component({
  selector: 'app-project-detail',
  standalone: true,
  imports: [
    CommonModule,
    MatCardModule,
    MatButtonModule,
    MatIconModule,
    MatChipsModule,
    MatDividerModule,
    MatProgressSpinnerModule,
    MatSnackBarModule
  ],
  templateUrl: './project-detail.component.html',
  styleUrl: './project-detail.component.scss'
})
export class ProjectDetailComponent implements OnInit, AfterViewInit {
  @ViewChild('detailContainer') detailContainer!: ElementRef;
  @ViewChild('projectCard') projectCard!: ElementRef;
  @ViewChild('deploymentSection') deploymentSection!: ElementRef;
  @ViewChild('actionButtons') actionButtons!: ElementRef;

  projectId!: string;
  projectDetail: ProjectDetail | null = null;
  isLoading = false;
  isDeploying = false;

  constructor(
    private route: ActivatedRoute,
    private router: Router,
    private projectService: ProjectService,
    private snackBar: MatSnackBar
  ) {}

  ngOnInit(): void {
    this.projectId = this.route.snapshot.paramMap.get('projectId') || '';
    if (this.projectId) {
      this.loadProjectDetail();
    }
  }

  ngAfterViewInit(): void {
    this.animatePageEntrance();
  }

  private animatePageEntrance(): void {
    const tl = gsap.timeline();
    
    tl.fromTo(this.detailContainer.nativeElement,
      { opacity: 0, y: 50 },
      { opacity: 1, y: 0, duration: 1, ease: "power2.out" }
    )
    .fromTo(this.projectCard.nativeElement,
      { opacity: 0, scale: 0.95, y: 30 },
      { opacity: 1, scale: 1, y: 0, duration: 0.8, ease: "back.out(1.7)" },
      "-=0.6"
    )
    .fromTo(this.deploymentSection.nativeElement,
      { opacity: 0, y: 30 },
      { opacity: 1, y: 0, duration: 0.6, ease: "power2.out" },
      "-=0.4"
    )
    .fromTo(this.actionButtons.nativeElement,
      { opacity: 0, y: 20 },
      { opacity: 1, y: 0, duration: 0.5, ease: "power2.out" },
      "-=0.3"
    );
  }

  async loadProjectDetail(): Promise<void> {
    this.isLoading = true;
    try {
      this.projectDetail = await this.projectService.getProjectDetail(this.projectId);
    } catch (error) {
      console.error('Error loading project detail:', error);
      this.snackBar.open('Failed to load project details', 'Close', { duration: 3000 });
    } finally {
      this.isLoading = false;
    }
  }

  async deployProject(): Promise<void> {
    if (!this.projectDetail) return;

    this.isDeploying = true;
    try {
      const response = await this.projectService.deployProject(this.projectId);
      
      // Animate success
      gsap.to(this.actionButtons.nativeElement, {
        scale: 1.05,
        duration: 0.2,
        yoyo: true,
        repeat: 1,
        ease: "power2.inOut"
      });

      this.snackBar.open('Deployment started successfully!', 'Close', { duration: 3000 });
      
      // Navigate to deployment logs page
      setTimeout(() => {
        this.router.navigate(['/projects', this.projectId, response.deployment_id], {
          state: {
            url: response.url,
            status: response.status
          }
        });
      }, 1000);
      
    } catch (error) {
      console.error('Error deploying project:', error);
      this.snackBar.open('Failed to start deployment', 'Close', { duration: 3000 });
    } finally {
      this.isDeploying = false;
    }
  }

  frameworkLogo(framework: string): string {
    const f = (framework || '').toLowerCase();
    if (f.includes('react')) return 'assets/frameswork_logos/react.svg';
    if (f.includes('angular')) return 'assets/frameswork_logos/angular.svg';
    return 'assets/logo.png';
  }

  openRepo(url: string): void {
    window.open(url, '_blank');
  }

  openSubdomain(): void {
    if (this.projectDetail?.subdomain) {
      window.open(`http://${this.projectDetail.subdomain}`, '_blank');
    }
  }

  goBack(): void {
    this.router.navigate(['/projects']);
  }
}
