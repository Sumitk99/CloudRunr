import { Component, OnInit, ElementRef, ViewChild, AfterViewInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { MatTabsModule } from '@angular/material/tabs';
import { MatCardModule } from '@angular/material/card';
import { gsap } from 'gsap';
import { LoginComponent } from './login/login.component';
import { SignupComponent } from './signup/signup.component';

@Component({
  selector: 'app-auth',
  standalone: true,
  imports: [CommonModule, MatTabsModule, MatCardModule, LoginComponent, SignupComponent],
  templateUrl: './auth.component.html',
  styleUrl: './auth.component.scss'
})
export class AuthComponent implements OnInit, AfterViewInit {
  @ViewChild('authContainer') authContainer!: ElementRef;
  @ViewChild('authCard') authCard!: ElementRef;
  @ViewChild('tabsContainer') tabsContainer!: ElementRef;

  selectedTabIndex = 0;

  constructor() {}

  ngOnInit(): void {
    // Initialize component logic
  }

  ngAfterViewInit(): void {
    this.animateAuthPage();
  }

  private animateAuthPage(): void {
    const tl = gsap.timeline();
    
    tl.fromTo(this.authContainer.nativeElement,
      { opacity: 0, scale: 0.95 },
      { opacity: 1, scale: 1, duration: 1, ease: "power2.out" }
    )
    .fromTo(this.authCard.nativeElement,
      { opacity: 0, y: 50 },
      { opacity: 1, y: 0, duration: 0.8, ease: "back.out(1.7)" },
      "-=0.5"
    )
    .fromTo(this.tabsContainer.nativeElement,
      { opacity: 0, y: 30 },
      { opacity: 1, y: 0, duration: 0.6, ease: "power2.out" },
      "-=0.3"
    );
  }

  onTabChange(index: number): void {
    this.selectedTabIndex = index;
    
    // Add tab change animation
    gsap.to(this.tabsContainer.nativeElement, {
      scale: 1.02,
      duration: 0.2,
      yoyo: true,
      repeat: 1,
      ease: "power2.out"
    });
  }

  // Hover animations
  onCardHover(event: any): void {
    gsap.to(event.target, { y: -5, scale: 1.02, duration: 0.3, ease: "power2.out" });
  }

  onCardLeave(event: any): void {
    gsap.to(event.target, { y: 0, scale: 1, duration: 0.3, ease: "power2.out" });
  }
}
