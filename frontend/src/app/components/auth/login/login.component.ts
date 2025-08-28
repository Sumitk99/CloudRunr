import { Component, OnInit, ElementRef, ViewChild, AfterViewInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormBuilder, FormGroup, Validators, ReactiveFormsModule } from '@angular/forms';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatSnackBar } from '@angular/material/snack-bar';
import { gsap } from 'gsap';
import { AuthService, LoginRequest } from '../../../services/auth.service';

@Component({
  selector: 'app-login',
  standalone: true,
  imports: [
    CommonModule,
    ReactiveFormsModule,
    MatFormFieldModule,
    MatInputModule,
    MatButtonModule,
    MatIconModule,
    MatProgressSpinnerModule
  ],
  templateUrl: './login.component.html',
  styleUrl: './login.component.scss'
})
export class LoginComponent implements OnInit, AfterViewInit {
  @ViewChild('loginForm') loginForm!: ElementRef;
  @ViewChild('emailField') emailField!: ElementRef;
  @ViewChild('passwordField') passwordField!: ElementRef;
  @ViewChild('submitButton') submitButton!: ElementRef;

  loginFormGroup: FormGroup;
  isLoading = false;
  hidePassword = true;

  constructor(
    private fb: FormBuilder,
    private authService: AuthService,
    private snackBar: MatSnackBar
  ) {
    this.loginFormGroup = this.fb.group({
      email: ['', [Validators.required, Validators.email]],
      password: ['', [Validators.required, Validators.minLength(6)]]
    });
  }

  ngOnInit(): void {
    // Initialize component logic
  }

  ngAfterViewInit(): void {
    this.animateLoginForm();
  }

  private animateLoginForm(): void {
    const tl = gsap.timeline();
    
    tl.fromTo(this.loginForm.nativeElement,
      { opacity: 0, y: 50 },
      { opacity: 1, y: 0, duration: 0.8, ease: "power2.out" }
    )
    .fromTo(this.emailField.nativeElement,
      { opacity: 0, x: -30 },
      { opacity: 1, x: 0, duration: 0.6, ease: "back.out(1.7)" },
      "-=0.4"
    )
    .fromTo(this.passwordField.nativeElement,
      { opacity: 0, x: -30 },
      { opacity: 1, x: 0, duration: 0.6, ease: "back.out(1.7)" },
      "-=0.3"
    )
    .fromTo(this.submitButton.nativeElement,
      { opacity: 0, y: 20 },
      { opacity: 1, y: 0, duration: 0.6, ease: "back.out(1.7)" },
      "-=0.2"
    );
  }

  async onSubmit(): Promise<void> {
    if (this.loginFormGroup.valid && !this.isLoading) {
      this.isLoading = true;
      
      try {
        const credentials: LoginRequest = {
          email: this.loginFormGroup.value.email,
          password: this.loginFormGroup.value.password
        };

        await this.authService.login(credentials);
        
        // Success animation
        gsap.to(this.submitButton.nativeElement, {
          scale: 1.1,
          duration: 0.2,
          yoyo: true,
          repeat: 1,
          onComplete: () => {
            this.snackBar.open('Login successful!', 'Welcome back!', {
              duration: 3000,
              panelClass: ['success-snackbar']
            });
          }
        });

      } catch (error) {
        console.error('Login error:', error);
        this.snackBar.open('Login failed. Please try again.', 'Error', {
          duration: 4000,
          panelClass: ['error-snackbar']
        });
      } finally {
        this.isLoading = false;
      }
    } else {
      this.markFormGroupTouched();
    }
  }

  private markFormGroupTouched(): void {
    Object.keys(this.loginFormGroup.controls).forEach(key => {
      const control = this.loginFormGroup.get(key);
      control?.markAsTouched();
    });
  }

  // Hover animations
  onFieldHover(event: any): void {
    gsap.to(event.target, { scale: 1.02, duration: 0.2, ease: "power2.out" });
  }

  onFieldLeave(event: any): void {
    gsap.to(event.target, { scale: 1, duration: 0.2, ease: "power2.out" });
  }

  onButtonHover(event: any): void {
    gsap.to(event.target, { scale: 1.05, duration: 0.2, ease: "power2.out" });
  }

  onButtonLeave(event: any): void {
    gsap.to(event.target, { scale: 1, duration: 0.2, ease: "power2.out" });
  }
}
