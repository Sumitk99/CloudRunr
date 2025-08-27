import { Component, OnInit, ElementRef, ViewChild, AfterViewInit } from '@angular/core';
import { Router } from '@angular/router';
import { CommonModule } from '@angular/common';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatCardModule } from '@angular/material/card';
import { MatDividerModule } from '@angular/material/divider';
import { gsap } from 'gsap';
import { ScrollTrigger } from 'gsap/ScrollTrigger';

gsap.registerPlugin(ScrollTrigger);

@Component({
  selector: 'app-landing-page',
  standalone: true,
  imports: [CommonModule, MatButtonModule, MatIconModule, MatCardModule, MatDividerModule],
  templateUrl: './landing-page.component.html',
  styleUrl: './landing-page.component.scss'
})
export class LandingPageComponent implements OnInit, AfterViewInit {
  @ViewChild('heroSection') heroSection!: ElementRef;
  @ViewChild('heroTitle') heroTitle!: ElementRef;
  @ViewChild('heroSubtitle') heroSubtitle!: ElementRef;
  @ViewChild('heroButtons') heroButtons!: ElementRef;
  @ViewChild('featuresSection') featuresSection!: ElementRef;
  @ViewChild('featureCards') featureCards!: ElementRef;
  @ViewChild('ctaSection') ctaSection!: ElementRef;
  @ViewChild('floatingShapes') floatingShapes!: ElementRef;

  constructor(private router: Router) {}

  ngOnInit(): void {
    // Initialize component logic
  }

  ngAfterViewInit(): void {
    this.initializeAnimations();
  }

  private initializeAnimations(): void {
    // Hero section animations
    this.animateHeroSection();
    
    // Floating shapes animations
    this.animateFloatingShapes();
    
    // Features section animations
    this.animateFeaturesSection();
    
    // CTA section animations
    this.animateCTASection();
  }

  private animateHeroSection(): void {
    const tl = gsap.timeline();
    
    tl.fromTo(this.heroTitle.nativeElement,
      { y: 100, opacity: 0 },
      { y: 0, opacity: 1, duration: 1.2, ease: "power3.out" }
    )
    .fromTo(this.heroSubtitle.nativeElement,
      { y: 50, opacity: 0 },
      { y: 0, opacity: 1, duration: 1, ease: "power2.out" },
      "-=0.5"
    )
    .fromTo(this.heroButtons.nativeElement,
      { y: 30, opacity: 0 },
      { y: 0, opacity: 1, duration: 0.8, ease: "back.out(1.7)" },
      "-=0.3"
    );
  }

  private animateFloatingShapes(): void {
    if (this.floatingShapes?.nativeElement) {
      const shapes = this.floatingShapes.nativeElement.children;
      
      // Create continuous floating animation for each shape
      gsap.to(shapes, {
        y: "random(-30, 30)",
        x: "random(-20, 20)",
        rotation: "random(-180, 180)",
        duration: "random(3, 6)",
        ease: "power1.inOut",
        repeat: -1,
        yoyo: true,
        stagger: {
          each: 0.5,
          from: "random"
        }
      });
    }
  }

  private animateFeaturesSection(): void {
    gsap.fromTo(this.featuresSection.nativeElement,
      { opacity: 0, y: 100 },
      { 
        opacity: 1, 
        y: 0, 
        duration: 1.2, 
        ease: "power2.out",
        scrollTrigger: {
          trigger: this.featuresSection.nativeElement,
          start: "top 80%",
          end: "bottom 20%",
          toggleActions: "play none none reverse"
        }
      }
    );

    // Animate feature cards individually with stagger
    const cards = this.featureCards.nativeElement.children;
    gsap.fromTo(cards,
      { opacity: 0, y: 50, scale: 0.9 },
      { 
        opacity: 1, 
        y: 0, 
        scale: 1, 
        duration: 0.8, 
        stagger: 0.2,
        ease: "back.out(1.7)",
        scrollTrigger: {
          trigger: this.featureCards.nativeElement,
          start: "top 80%",
          end: "bottom 20%",
          toggleActions: "play none none reverse"
        }
      }
    );
  }

  private animateCTASection(): void {
    gsap.fromTo(this.ctaSection.nativeElement,
      { opacity: 0, scale: 0.95 },
      { 
        opacity: 1, 
        scale: 1, 
        duration: 1, 
        ease: "power2.out",
        scrollTrigger: {
          trigger: this.ctaSection.nativeElement,
          start: "top 80%",
          end: "bottom 20%",
          toggleActions: "play none none reverse"
        }
      }
    );
  }

  onGetStartedClick(event: Event): void {
    const target = event.target as HTMLElement;
    gsap.to(target, { scale: 0.95, duration: 0.1, yoyo: true, repeat: 1 });
    // Will implement later
    console.log('Get Started clicked');
  }

  onLearnMoreClick(event: Event): void {
    const target = event.target as HTMLElement;
    gsap.to(target, { scale: 0.95, duration: 0.1, yoyo: true, repeat: 1 });
    // Will implement later
    console.log('Learn More clicked');
  }

  onDeployClick(event: Event): void {
    const target = event.target as HTMLElement;
    gsap.to(target, { scale: 0.95, duration: 0.1, yoyo: true, repeat: 1 });
    // Will implement later
    console.log('Deploy clicked');
  }

  // Hover animations
  onButtonHover(event: any): void {
    gsap.to(event.target, { scale: 1.05, duration: 0.2, ease: "power2.out" });
  }

  onButtonLeave(event: any): void {
    gsap.to(event.target, { scale: 1, duration: 0.2, ease: "power2.out" });
  }

  onCardHover(event: any): void {
    gsap.to(event.target, { y: -10, scale: 1.02, duration: 0.3, ease: "power2.out" });
  }

  onCardLeave(event: any): void {
    gsap.to(event.target, { y: 0, scale: 1, duration: 0.3, ease: "power2.out" });
  }
}
