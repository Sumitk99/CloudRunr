import { Injectable } from '@angular/core';
import { BehaviorSubject, Observable } from 'rxjs';
import { Router } from '@angular/router';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { environment } from '../../environments/environment';

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

  constructor(
    private router: Router,
    private http: HttpClient
  ) {
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
      const response = await this.http.post<User>(
        `${environment.backendUrl}${environment.apiEndpoints.login}`,
        credentials,
        {
          headers: new HttpHeaders({
            'Content-Type': 'application/json'
          })
        }
      ).toPromise();

      if (response) {
        this.setUser(response);
        return response;
      } else {
        throw new Error('No response from server');
      }
    } catch (error) {
      console.error('Login error:', error);
      throw error;
    }
  }

  async signup(userData: SignupRequest): Promise<User> {
    try {
      const response = await this.http.post<User>(
        `${environment.backendUrl}${environment.apiEndpoints.signup}`,
        userData,
        {
          headers: new HttpHeaders({
            'Content-Type': 'application/json'
          })
        }
      ).toPromise();

      if (response) {
        this.setUser(response);
        return response;
      } else {
        throw new Error('No response from server');
      }
    } catch (error) {
      console.error('Signup error:', error);
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

  // Helper method to create auth headers for future API calls
  getAuthHeaders(): HttpHeaders {
    const token = this.getToken();
    return new HttpHeaders({
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${token}`
    });
  }
}
