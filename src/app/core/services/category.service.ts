import { Injectable } from '@angular/core';
import { environment } from '../../../environments/environment.development';
import { HttpClient, HttpParams } from '@angular/common/http';
import { catchError, map, Observable } from 'rxjs';
import { PaginatedResponse, PaginationParams, ProductCategory } from '../interfaces/interfaces';

@Injectable({
  providedIn: 'root',
})
export class CategoryService {
  private readonly apiURL = environment.baseURL + '/categories';

  constructor(private http: HttpClient) {}

  getPaginatedCategories(params: PaginationParams): Observable<PaginatedResponse<ProductCategory>> {
    let httpParams = new HttpParams()
      .set('page', params.page.toString())
      .set('limit', params.limit.toString());

    // Add optional parameters if they exist
    if (params.search) {
      httpParams = httpParams.set('search', params.search);
    }

    if (params.sortBy) {
      httpParams = httpParams.set('sortBy', params.sortBy);
    }

    if (params.sortOrder) {
      httpParams = httpParams.set('sortOrder', params.sortOrder);
    }

    return this.http.get<PaginatedResponse<ProductCategory>>(`${this.apiURL}/products`, {
      params: httpParams,
    });
  }

  getAllCategories(): Observable<ProductCategory[]> {
    return this.http.get<{ data: ProductCategory[] }>(this.apiURL).pipe(
      map((response) => {
        return response.data;
      })
    );
  }

  

  getCategoryById(id: number): Observable<ProductCategory> {
    return this.http.get<{ data: ProductCategory }>(`${this.apiURL}/${id}`).pipe(
      map((response) => {
        return response.data;
      })
    );
  }

  createCategory(category: ProductCategory): Observable<ProductCategory> {
    return this.http.post<ProductCategory>(this.apiURL, category);
  }

  updateCategory(id: number, category: ProductCategory): Observable<ProductCategory> {
    return this.http.patch<ProductCategory>(`${this.apiURL}/${id}`, category);
  }

  deleteCategory(id: number): Observable<any> {
    return this.http.delete<any>(`${this.apiURL}/${id}`);
  }
}
