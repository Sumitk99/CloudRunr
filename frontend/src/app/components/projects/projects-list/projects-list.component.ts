import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { Router } from '@angular/router';
import { MatCardModule } from '@angular/material/card';
import { MatIconModule } from '@angular/material/icon';
import { MatButtonModule } from '@angular/material/button';
import { ProjectService, ProjectItem } from '../../../services/project.service';

@Component({
  selector: 'app-projects-list',
  standalone: true,
  imports: [CommonModule, MatCardModule, MatIconModule, MatButtonModule],
  templateUrl: './projects-list.component.html',
  styleUrl: './projects-list.component.scss'
})
export class ProjectsListComponent implements OnInit {
  isLoading = false;
  projects: ProjectItem[] = [];

  constructor(private projectsSvc: ProjectService, private router: Router) {}

  ngOnInit(): void {
    this.load();
  }

  async load(): Promise<void> {
    this.isLoading = true;
    try {
      const res = await this.projectsSvc.fetchProjects();
      this.projects = res.projects || [];
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
    this.router.navigate([`/projects/${id}`]);
  }
}
