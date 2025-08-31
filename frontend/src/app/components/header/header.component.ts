import { Component, OnInit, ElementRef, ViewChild, AfterViewInit } from '@angular/core';
import { Router } from '@angular/router';
import { CommonModule } from '@angular/common';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatToolbarModule } from '@angular/material/toolbar';
import { MatMenuModule } from '@angular/material/menu';
import { MatDividerModule } from '@angular/material/divider';
import { gsap } from 'gsap';
import { AuthService, User } from '../../services/auth.service';

@Component({
  selector: 'app-header',
  standalone: true,
  imports: [CommonModule, MatButtonModule, MatIconModule, MatToolbarModule, MatMenuModule, MatDividerModule],
  templateUrl: './header.component.html',
  styleUrl: './header.component.scss'
})
export class HeaderComponent implements OnInit, AfterViewInit {
  @ViewChild('headerContainer') headerContainer!: ElementRef;
  @ViewChild('logoSection') logoSection!: ElementRef;
  @ViewChild('authSection') authSection!: ElementRef;
  
  isLoggedIn = false;
  currentUser: User | null = null;

  constructor(
    private router: Router,
    private authService: AuthService
  ) {}

  ngOnInit(): void {
    // Subscribe to auth state changes
    this.authService.currentUser$.subscribe(user => {
      this.currentUser = user;
      this.isLoggedIn = !!user;
    });
  }

  ngAfterViewInit(): void {
    this.animateHeader();
  }

  private animateHeader(): void {
    // Animate header entrance
    gsap.fromTo(this.headerContainer.nativeElement, 
      { y: -100, opacity: 0 },
      { y: 0, opacity: 1, duration: 1, ease: "power3.out" }
    );

    // Animate logo section
    gsap.fromTo(this.logoSection.nativeElement,
      { x: -50, opacity: 0 },
      { x: 0, opacity: 1, duration: 1.2, delay: 0.3, ease: "back.out(1.7)" }
    );

    // Animate auth section
    gsap.fromTo(this.authSection.nativeElement,
      { x: 50, opacity: 0 },
      { x: 0, opacity: 1, duration: 1.2, delay: 0.5, ease: "back.out(1.7)" }
    );
  }

  onLoginClick(event: Event): void {
    // Add click animation
    const target = event.target as HTMLElement;
    gsap.to(target, { scale: 0.95, duration: 0.1, yoyo: true, repeat: 1 });
    setTimeout(() => {
      this.router.navigate(['/auth']);
    }, 200);
  }

  onProfileClick(event: Event): void {
    // Add click animation
    const target = event.target as HTMLElement;
    gsap.to(target, { scale: 0.95, duration: 0.1, yoyo: true, repeat: 1 });
  }

  onNewProjectClick(event: Event): void {
    const target = event.target as HTMLElement;
    gsap.to(target, { scale: 0.95, duration: 0.1, yoyo: true, repeat: 1 });
    this.router.navigate(['/projects/new']);
  }

  onProjectsClick(): void {
    this.router.navigate(['/projects']);
  }

  onMyProjectsClick(event: Event): void {
    const target = event.target as HTMLElement;
    gsap.to(target, { scale: 0.95, duration: 0.1, yoyo: true, repeat: 1 });
    this.router.navigate(['/projects']);
  }

  onSettingsClick(): void {
    // TODO: Implement settings page navigation
    console.log('Settings clicked - to be implemented');
  }

  onBillingClick(): void {
    // TODO: Implement billing page navigation
    console.log('Billing clicked - to be implemented');
  }

  onLogoutClick(): void {
    this.authService.logout();
  }

  onLogoClick(): void {
    this.router.navigate(['/']);
  }

  // Hover animations
  onButtonHover(event: any): void {
    gsap.to(event.target, { scale: 1.05, duration: 0.2, ease: "power2.out" });
  }

  onButtonLeave(event: any): void {
    gsap.to(event.target, { scale: 1, duration: 0.2, ease: "power2.out" });
  }
}
