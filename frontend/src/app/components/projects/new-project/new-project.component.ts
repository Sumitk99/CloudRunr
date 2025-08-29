import { Component, OnInit, AfterViewInit, ElementRef, ViewChild } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormBuilder, FormGroup, Validators, ReactiveFormsModule } from '@angular/forms';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatSelectModule } from '@angular/material/select';
import { MatButtonModule } from '@angular/material/button';
import { MatCardModule } from '@angular/material/card';
import { MatSnackBar } from '@angular/material/snack-bar';
import { Router } from '@angular/router';
import { gsap } from 'gsap';
import { ProjectService, NewProjectRequest } from '../../../services/project.service';

@Component({
  selector: 'app-new-project',
  standalone: true,
  imports: [
    CommonModule,
    ReactiveFormsModule,
    MatFormFieldModule,
    MatInputModule,
    MatSelectModule,
    MatButtonModule,
    MatCardModule
  ],
  templateUrl: './new-project.component.html',
  styleUrl: './new-project.component.scss'
})
export class NewProjectComponent implements OnInit, AfterViewInit {
  @ViewChild('card') card!: ElementRef;
  @ViewChild('formEl') formEl!: ElementRef;

  form: FormGroup;
  isLoading = false;
  frameworks = [
    { value: 'REACT', label: 'React' },
    { value: 'ANGULAR', label: 'Angular' }
  ];

  constructor(
    private fb: FormBuilder,
    private projects: ProjectService,
    private snack: MatSnackBar,
    private router: Router
  ) {
    this.form = this.fb.group({
      git_url: ['', [Validators.required]],
      framework: ['React', [Validators.required]],
      dist_folder: ['dist', [Validators.required]],
      project_id: ['', [Validators.required]],
      name: ['', [Validators.required]],
      run_command: ['']
    });
  }

  ngOnInit(): void {}

  ngAfterViewInit(): void {
    gsap.fromTo(this.card.nativeElement, { opacity: 0, y: 30 }, { opacity: 1, y: 0, duration: 0.8, ease: 'power2.out' });
    gsap.fromTo(this.formEl.nativeElement, { opacity: 0, y: 20 }, { opacity: 1, y: 0, duration: 0.6, delay: 0.2, ease: 'power2.out' });
  }

  async onSubmit(): Promise<void> {
    if (this.form.invalid || this.isLoading) {
      this.form.markAllAsTouched();
      return;
    }
    this.isLoading = true;
    try {
      const payload: NewProjectRequest = {
        git_url: this.form.value.git_url,
        framework: this.form.value.framework,
        dist_folder: this.form.value.dist_folder,
        project_id: this.form.value.project_id,
        name: this.form.value.name,
        run_command: this.form.value.run_command || null
      };
      const res = await this.projects.createNewProject(payload);
      this.snack.open('Project created. Starting deployment...', 'OK', { duration: 3000, panelClass: ['success-snackbar'] });
      const projectId = payload.project_id;
      const deploymentId = res.deployment_id;
      this.router.navigate([`/projects/${projectId}/${deploymentId}`]);
    } catch (e) {
      console.error(e);
      this.snack.open('Failed to create project', 'Error', { duration: 4000, panelClass: ['error-snackbar'] });
    } finally {
      this.isLoading = false;
    }
  }
}
