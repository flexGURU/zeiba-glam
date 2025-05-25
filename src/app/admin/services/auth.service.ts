import { Injectable } from '@angular/core';
import { User } from './user.interface';
import { Observable, of } from 'rxjs';

@Injectable({
  providedIn: 'root',
})
export class AuthService {
  constructor() {}

  login(email: string, password: string): Observable<User | null> {
    return of(null);
  }
}
