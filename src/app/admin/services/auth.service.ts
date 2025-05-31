import { Injectable } from '@angular/core';
import { LoginResponse, User } from './user.interface';
import { map, Observable, of, tap } from 'rxjs';
import { environment } from '../../../environments/environment.development';
import { HttpClient, HttpHeaders } from '@angular/common/http';

@Injectable({
  providedIn: 'root',
})
export class AuthService {
  private readonly apiURL = environment.baseURL;
  private static readonly accessKey = 'JWT_ACCESS_KEY';

  constructor(private http: HttpClient) {}

  get jwt(): string {
    return sessionStorage.getItem(AuthService.accessKey) ?? '';
  }

  private set jwt(value: string) {
    sessionStorage.setItem(AuthService.accessKey, value);
  }

  login(email: string, password: string): Observable<LoginResponse> {
    const login = this.http
      .post<LoginResponse>(`${this.apiURL}/auth/login`, { email, password })
      .pipe(
        tap((resp) => {
          console.log('ssss', resp);
          this.jwt = resp.data.access_token;
        })
      );

    return login;
  }

  isLoggedIn(): boolean {
    return !!this.jwt;
  }
}
