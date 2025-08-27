import { Component, OnInit, ElementRef, ViewChild, AfterViewInit } from '@angular/core';
import { Router } from '@angular/router';
import { CommonModule } from '@angular/common';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatToolbarModule } from '@angular/material/toolbar';
import { gsap } from 'gsap';

@Component({
  selector: 'app-header',
  standalone: true,
  imports: [CommonModule, MatButtonModule, MatIconModule, MatToolbarModule],
  templateUrl: './header.component.html',
  styleUrl: './header.component.scss'
})
export class HeaderComponent implements OnInit, AfterViewInit {
  @ViewChild('headerContainer') headerContainer!: ElementRef;
  @ViewChild('logoSection') logoSection!: ElementRef;
  @ViewChild('authSection') authSection!: ElementRef;
  
  isLoggedIn = false; // This will be replaced with actual auth service later

  constructor(private router: Router) {}

  ngOnInit(): void {
    // Initialize any component logic
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
      this.router.navigate(['/login']);
    }, 200);
  }

  onProfileClick(event: Event): void {
    // Add click animation
    const target = event.target as HTMLElement;
    gsap.to(target, { scale: 0.95, duration: 0.1, yoyo: true, repeat: 1 });
    // Will implement profile page later
    console.log('Profile clicked');
  }

  // Hover animations
  onButtonHover(event: any): void {
    gsap.to(event.target, { scale: 1.05, duration: 0.2, ease: "power2.out" });
  }

  onButtonLeave(event: any): void {
    gsap.to(event.target, { scale: 1, duration: 0.2, ease: "power2.out" });
  }
}
