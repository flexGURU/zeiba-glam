import { inject, Injectable } from '@angular/core';
import { environment } from '../../../environments/environment.development';
import { HttpClient } from '@angular/common/http';
import { Observable, catchError, map, throwError } from 'rxjs';
import { DashboardStats } from '../../core/interfaces/interfaces';

@Injectable({
  providedIn: 'root',
})
export class DashboardService {
  private readonly apiURL = environment.baseURL;

  private http = inject(HttpClient);

  getDashboardData(): Observable<DashboardStats> {
    return this.http.get<{ data: DashboardStats }>(`${this.apiURL}/helpers/dashboard-stats`).pipe(
      map((response) => {
        return response.data;
      })
    );
  }
}
