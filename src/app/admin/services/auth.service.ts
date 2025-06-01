import { Injectable } from '@angular/core';
import { LoginResponse, User } from './user.interface';
import {
  BehaviorSubject,
  catchError,
  filter,
  map,
  Observable,
  of,
  take,
  tap,
  throwError,
} from 'rxjs';
import { environment } from '../../../environments/environment.development';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Router } from '@angular/router';

@Injectable({
  providedIn: 'root',
})
export class AuthService {
  private readonly apiURL = environment.baseURL;
  private static readonly accessKey = 'JWT_ACCESS_KEY';
  private static readonly refreshKey = 'JWT_REFRESH_KEY';

  // Subject to handle concurrent refresh requests
  private refreshTokenSubject: BehaviorSubject<string | null> = new BehaviorSubject<string | null>(
    null
  );
  private isRefreshing = false;

  constructor(
    private http: HttpClient,
    private router: Router
  ) {}

  get jwt(): string {
    return sessionStorage.getItem(AuthService.accessKey) ?? '';
  }

  private set jwt(value: string) {
    sessionStorage.setItem(AuthService.accessKey, value);
  }

  get refreshToken(): string {
    // Fixed: was using accessKey instead of refreshKey
    return sessionStorage.getItem(AuthService.refreshKey) ?? '';
  }

  private set refreshToken(value: string) {
    sessionStorage.setItem(AuthService.refreshKey, value);
  }

  login(email: string, password: string): Observable<LoginResponse> {
    return this.http
      .post<LoginResponse>(
        `${this.apiURL}/auth/login`,
        { email, password },
        { withCredentials: true }
      )
      .pipe(
        tap((resp) => {
          console.log('Login response:', resp);
          this.jwt = resp.data.access_token;
          this.refreshToken = resp.data.refresh_token;
          this.router.navigate(['/admin/dashboard']);
        })
      );
  }

  isLoggedIn(): boolean {
    return !!this.jwt;
  }

  handleTokenRefresh(): Observable<string> {
    if (this.isRefreshing) {
      // If refresh is already in progress, wait for it to complete
      return this.refreshTokenSubject.pipe(
        filter((token) => token != null),
        take(1)
      );
    } else {
      this.isRefreshing = true;
      this.refreshTokenSubject.next(null);

      const refreshToken = this.refreshToken;

      if (!refreshToken) {
        this.logout();
        return throwError(() => new Error('No refresh token available'));
      }

      return this.refreshAccessToken().pipe(
        tap((resp) => {
          this.isRefreshing = false;
          this.jwt = resp.data.access_token;
          this.refreshTokenSubject.next(resp.data.access_token);
        }),
        map((resp) => resp.data.access_token),
        catchError((error) => {
          this.isRefreshing = false;
          this.refreshTokenSubject.next(null);
          this.logout();
          return throwError(() => error);
        })
      );
    }
  }

  private refreshAccessToken(): Observable<any> {
    const headers = new HttpHeaders({
      Authorization: `Bearer ${this.refreshToken}`,
    });

    return this.http.get<any>(`${this.apiURL}/auth/refresh`, { headers }).pipe(
      tap((resp) => {
        console.log('Token refresh response:', resp);
      })
    );
  }

  logout(): void {
    sessionStorage.removeItem(AuthService.accessKey);
    sessionStorage.removeItem(AuthService.refreshKey);
    this.refreshTokenSubject.next(null);
    this.router.navigate(['admin/login']);
  }
}
