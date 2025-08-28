import { Injectable } from '@angular/core';
import { BehaviorSubject, Observable } from 'rxjs';
import { Router } from '@angular/router';

export interface User {
  name: string;
  email: string;
  user_id: string;
  token: string;
  refresh_token: string;
  github_id?: string;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface SignupRequest {
  name: string;
  email: string;
  password: string;
}

@Injectable({
  providedIn: 'root'
})
export class AuthService {
  private currentUserSubject = new BehaviorSubject<User | null>(null);
  public currentUser$ = this.currentUserSubject.asObservable();

  private readonly USER_KEY = 'cloudrunr_user';
  private readonly TOKEN_KEY = 'cloudrunr_token';

  constructor(private router: Router) {
    this.loadUserFromStorage();
  }

  private loadUserFromStorage(): void {
    const userStr = localStorage.getItem(this.USER_KEY);
    if (userStr) {
      try {
        const user = JSON.parse(userStr);
        this.currentUserSubject.next(user);
      } catch (error) {
        console.error('Error parsing user from storage:', error);
        this.clearAuth();
      }
    }
  }

  isLoggedIn(): boolean {
    const user = this.currentUserSubject.value;
    return !!(user && user.token);
  }

  getCurrentUser(): User | null {
    return this.currentUserSubject.value;
  }

  getToken(): string | null {
    const user = this.currentUserSubject.value;
    return user?.token || null;
  }

  async login(credentials: LoginRequest): Promise<User> {
    try {
      // TODO: Replace with actual API call
      const response = await this.mockLoginAPI(credentials);
      this.setUser(response);
      return response;
    } catch (error) {
      throw error;
    }
  }

  async signup(userData: SignupRequest): Promise<User> {
    try {
      // TODO: Replace with actual API call
      const response = await this.mockSignupAPI(userData);
      this.setUser(response);
      return response;
    } catch (error) {
      throw error;
    }
  }

  logout(): void {
    this.clearAuth();
    this.router.navigate(['/']);
  }

  private setUser(user: User): void {
    localStorage.setItem(this.USER_KEY, JSON.stringify(user));
    localStorage.setItem(this.TOKEN_KEY, user.token);
    this.currentUserSubject.next(user);
  }

  private clearAuth(): void {
    localStorage.removeItem(this.USER_KEY);
    localStorage.removeItem(this.TOKEN_KEY);
    this.currentUserSubject.next(null);
  }

  // Mock API calls - Replace with actual backend calls
  private async mockLoginAPI(credentials: LoginRequest): Promise<User> {
    // Simulate API delay
    await new Promise(resolve => setTimeout(resolve, 1000));
    
    // Mock response
    return {
      name: "Demo User",
      email: credentials.email,
      user_id: "31tvhGjSCAl7TRG0eMAetSA33Ym",
      token: "mock_jwt_token_" + Date.now(),
      refresh_token: "mock_refresh_token_" + Date.now(),
      github_id: ""
    };
  }

  private async mockSignupAPI(userData: SignupRequest): Promise<User> {
    // Simulate API delay
    await new Promise(resolve => setTimeout(resolve, 1000));
    
    // Mock response
    return {
      name: userData.name,
      email: userData.email,
      user_id: "31tvhGjSCAl7TRG0eMAetSA33Ym",
      token: "mock_jwt_token_" + Date.now(),
      refresh_token: "mock_refresh_token_" + Date.now()
    };
  }
}
