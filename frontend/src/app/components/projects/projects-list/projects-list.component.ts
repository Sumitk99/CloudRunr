import { Component, OnInit, AfterViewInit, ElementRef, ViewChild, QueryList, ViewChildren } from '@angular/core';
import { CommonModule } from '@angular/common';
import { Router } from '@angular/router';
import { MatCardModule } from '@angular/material/card';
import { MatIconModule } from '@angular/material/icon';
import { MatButtonModule } from '@angular/material/button';
import { gsap } from 'gsap';
import { ProjectService, ProjectItem } from '../../../services/project.service';

@Component({
  selector: 'app-projects-list',
  standalone: true,
  imports: [CommonModule, MatCardModule, MatIconModule, MatButtonModule],
  templateUrl: './projects-list.component.html',
  styleUrl: './projects-list.component.scss'
})
export class ProjectsListComponent implements OnInit, AfterViewInit {
  @ViewChild('projectsHeader') projectsHeader!: ElementRef;
  @ViewChild('projectsGrid') projectsGrid!: ElementRef;
  @ViewChild('emptyState') emptyState!: ElementRef;
  @ViewChildren('projectCard') projectCards!: QueryList<ElementRef>;

  isLoading = false;
  projects: ProjectItem[] = [];
  hasAnimated = false;

  constructor(private projectsSvc: ProjectService, private router: Router) {}

  ngOnInit(): void {
    this.load();
  }

  ngAfterViewInit(): void {
    // Initial page load animation
    this.animatePageEntrance();
  }

  private animatePageEntrance(): void {
    // Animate header entrance
    if (this.projectsHeader) {
      gsap.fromTo(this.projectsHeader.nativeElement,
        { opacity: 0, y: -30 },
        { opacity: 1, y: 0, duration: 0.8, ease: "power2.out" }
      );
    }

    // If projects are already loaded, animate them
    if (this.projects.length > 0 && !this.hasAnimated) {
      setTimeout(() => this.animateProjectCards(), 200);
    }

    // Animate empty state if no projects
    if (this.projects.length === 0 && !this.isLoading && this.emptyState) {
      setTimeout(() => {
        gsap.fromTo(this.emptyState.nativeElement,
          { opacity: 0, scale: 0.9 },
          { opacity: 1, scale: 1, duration: 0.6, ease: "back.out(1.7)" }
        );
      }, 400);
    }
  }

  private animateProjectCards(): void {
    if (this.hasAnimated) return;
    
    this.hasAnimated = true;
    const cards = this.projectCards.toArray();
    
    if (cards.length > 0) {
      // Animate grid container first
      gsap.fromTo(this.projectsGrid.nativeElement,
        { opacity: 0 },
        { opacity: 1, duration: 0.3 }
      );

      // Animate each card with staggered timing
      cards.forEach((card, index) => {
        gsap.fromTo(card.nativeElement,
          { 
            opacity: 0, 
            y: 50, 
            scale: 0.9,
            rotateX: -15
          },
          { 
            opacity: 1, 
            y: 0, 
            scale: 1,
            rotateX: 0,
            duration: 0.8, 
            delay: index * 0.1,
            ease: "back.out(1.7)"
          }
        );
      });
    }
  }

  async load(): Promise<void> {
    this.isLoading = true;
    this.hasAnimated = false;
    
    // Reset animations
    if (this.projectsGrid) {
      gsap.set(this.projectsGrid.nativeElement, { opacity: 0 });
    }

    try {
      const res = await this.projectsSvc.fetchProjects();
      this.projects = res.projects || [];
      
      // Trigger animations after projects are loaded
      setTimeout(() => {
        if (this.projects.length > 0) {
          this.animateProjectCards();
        } else if (this.emptyState) {
          // Animate empty state
          gsap.fromTo(this.emptyState.nativeElement,
            { opacity: 0, scale: 0.9 },
            { opacity: 1, scale: 1, duration: 0.6, ease: "back.out(1.7)" }
          );
        }
      }, 100);
    } finally {
      this.isLoading = false;
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

  openProject(id: string): void {
    if (id === 'new') {
      this.router.navigate(['/projects/new']);
    } else {
      this.router.navigate(['/projects', id]);
    }
  }

  trackByProjectId(index: number, project: ProjectItem): string {
    return project.project_id;
  }
}
