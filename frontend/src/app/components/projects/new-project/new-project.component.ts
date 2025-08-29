import { Component, OnInit, AfterViewInit, ElementRef, ViewChild, ViewChildren, QueryList } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormBuilder, FormGroup, Validators, ReactiveFormsModule } from '@angular/forms';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatSelectModule } from '@angular/material/select';
import { MatButtonModule } from '@angular/material/button';
import { MatCardModule } from '@angular/material/card';
import { MatIconModule } from '@angular/material/icon';
import { MatChipsModule } from '@angular/material/chips';
import { MatStepperModule } from '@angular/material/stepper';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatSnackBar, MatSnackBarModule } from '@angular/material/snack-bar';
import { Router } from '@angular/router';
import { gsap } from 'gsap';
import { ProjectService, NewProjectRequest } from '../../../services/project.service';

interface Framework {
  value: string;
  label: string;
  icon: string;
  description: string;
  buildFolder: string;
  buildCommand: string;
}

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
    MatCardModule,
    MatIconModule,
    MatChipsModule,
    MatStepperModule,
    MatProgressSpinnerModule,
    MatSnackBarModule
  ],
  templateUrl: './new-project.component.html',
  styleUrl: './new-project.component.scss'
})
export class NewProjectComponent implements OnInit, AfterViewInit {
  @ViewChild('pageContainer') pageContainer!: ElementRef;
  @ViewChild('headerSection') headerSection!: ElementRef;
  @ViewChild('formCard') formCard!: ElementRef;
  @ViewChild('frameworkSelection') frameworkSelection!: ElementRef;
  @ViewChild('configSection') configSection!: ElementRef;
  @ViewChild('submitSection') submitSection!: ElementRef;
  @ViewChildren('formField') formFields!: QueryList<ElementRef>;

  form: FormGroup;
  isLoading = false;
  selectedFramework: Framework | null = null;
  
  frameworks: Framework[] = [
    { 
      value: 'REACT', 
      label: 'React', 
      icon: 'assets/frameswork_logos/react.svg',
      description: 'A JavaScript library for building user interfaces',
      buildFolder: 'build',
      buildCommand: 'npm run build'
    },
    { 
      value: 'ANGULAR', 
      label: 'Angular', 
      icon: 'assets/frameswork_logos/angular.svg',
      description: 'Platform for building mobile and desktop web applications',
      buildFolder: 'dist',
      buildCommand: 'ng build'
    }
  ];

  constructor(
    private fb: FormBuilder,
    private projects: ProjectService,
    private snack: MatSnackBar,
    private router: Router
  ) {
    this.form = this.fb.group({
      git_url: ['', [Validators.required, Validators.pattern(/^https:\/\/github\.com\/.+\/.+$/)]],
      framework: ['', [Validators.required]],
      dist_folder: ['', [Validators.required]],
      project_id: ['', [Validators.required, Validators.pattern(/^[a-z0-9-]+$/)]],
      name: ['', [Validators.required, Validators.minLength(2)]],
      run_command: ['']
    });
  }

  ngOnInit(): void {
    // Watch for framework selection changes
    this.form.get('framework')?.valueChanges.subscribe(value => {
      const framework = this.frameworks.find(f => f.value === value);
      if (framework) {
        this.selectedFramework = framework;
        this.form.patchValue({
          dist_folder: framework.buildFolder,
          run_command: framework.buildCommand
        });
        this.animateFrameworkSelection();
      }
    });

    // Watch for name changes to auto-generate project ID
    this.form.get('name')?.valueChanges.subscribe(name => {
      if (name && !this.form.get('project_id')?.value) {
        const projectId = name.toLowerCase()
          .replace(/[^a-z0-9\s-]/g, '')
          .replace(/\s+/g, '-')
          .substring(0, 20);
        this.form.patchValue({ project_id: projectId });
      }
    });
  }

  ngAfterViewInit(): void {
    this.animatePageEntrance();
  }

  private animatePageEntrance(): void {
    const tl = gsap.timeline();
    
    // Animate page container
    tl.fromTo(this.pageContainer.nativeElement,
      { opacity: 0, y: 30 },
      { opacity: 1, y: 0, duration: 1, ease: "power2.out" }
    )
    
    // Animate header
    .fromTo(this.headerSection.nativeElement,
      { opacity: 0, y: -20, scale: 0.95 },
      { opacity: 1, y: 0, scale: 1, duration: 0.8, ease: "back.out(1.7)" },
      "-=0.6"
    )
    
    // Animate form card
    .fromTo(this.formCard.nativeElement,
      { opacity: 0, y: 40, rotateX: 10 },
      { opacity: 1, y: 0, rotateX: 0, duration: 0.8, ease: "power2.out" },
      "-=0.4"
    )
    
    // Animate form fields with stagger
    .fromTo(this.formFields.toArray().map(field => field.nativeElement),
      { opacity: 0, x: -20 },
      { opacity: 1, x: 0, duration: 0.6, stagger: 0.1, ease: "power2.out" },
      "-=0.3"
    );
  }

  private animateFrameworkSelection(): void {
    if (this.configSection) {
      gsap.fromTo(this.configSection.nativeElement,
        { opacity: 0.5, scale: 0.98 },
        { opacity: 1, scale: 1, duration: 0.4, ease: "power2.out" }
      );
    }
  }

  selectFramework(framework: Framework): void {
    this.form.patchValue({ framework: framework.value });
    
    // Animate selection
    gsap.to(`.framework-card.selected`, {
      scale: 1.02,
      duration: 0.2,
      ease: "power2.out"
    });
  }

  onGitUrlChange(): void {
    const gitUrl = this.form.get('git_url')?.value;
    if (gitUrl && !this.form.get('name')?.value) {
      // Extract repo name from git URL
      const repoName = gitUrl.split('/').pop()?.replace('.git', '') || '';
      if (repoName) {
        this.form.patchValue({ name: repoName });
      }
    }
  }

  async onSubmit(): Promise<void> {
    if (this.form.invalid || this.isLoading) {
      this.form.markAllAsTouched();
      this.animateFormErrors();
      return;
    }

    this.isLoading = true;
    this.animateSubmitLoading();

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
      
      // Success animation
      this.animateSuccess();
      
      this.snack.open('Project created successfully! Starting deployment...', 'OK', { 
        duration: 3000, 
        panelClass: ['success-snackbar'] 
      });

      const projectId = payload.project_id;
      const deploymentId = res.deployment_id;
      
      // Navigate after animation
      setTimeout(() => {
        this.router.navigate([`/projects/${projectId}/${deploymentId}`]);
      }, 1500);
      
    } catch (e) {
      console.error(e);
      this.animateError();
      this.snack.open('Failed to create project. Please check your details and try again.', 'Error', { 
        duration: 4000, 
        panelClass: ['error-snackbar'] 
      });
      this.isLoading = false;
    }
  }

  private animateFormErrors(): void {
    const invalidFields = Object.keys(this.form.controls)
      .filter(key => this.form.get(key)?.invalid)
      .map(key => document.querySelector(`[formControlName="${key}"]`))
      .filter(el => el);

    invalidFields.forEach(field => {
      gsap.to(field, {
        x: 10,
        duration: 0.1,
        yoyo: true,
        repeat: 3,
        ease: "power2.inOut"
      });
    });
  }

  private animateSubmitLoading(): void {
    gsap.to(this.submitSection.nativeElement, {
      scale: 0.98,
      duration: 0.2,
      ease: "power2.inOut"
    });
  }

  private animateSuccess(): void {
    gsap.to(this.formCard.nativeElement, {
      scale: 1.02,
      duration: 0.3,
      yoyo: true,
      repeat: 1,
      ease: "power2.inOut"
    });
  }

  private animateError(): void {
    gsap.to(this.formCard.nativeElement, {
      x: 10,
      duration: 0.1,
      yoyo: true,
      repeat: 3,
      ease: "power2.inOut"
    });
  }

  goBack(): void {
    this.router.navigate(['/projects']);
  }
}
