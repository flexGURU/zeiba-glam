import { Injectable } from '@angular/core';
import { environment } from '../../../environments/environment.development';
import { HttpClient } from '@angular/common/http';
import { map, Observable } from 'rxjs';
import { ProductCategory } from '../interfaces/interfaces';

@Injectable({
  providedIn: 'root',
})
export class CategoryService {
  private readonly apiURL = environment.baseURL + '/categories';

  constructor(private http: HttpClient) {}

  getAllCategories(): Observable<ProductCategory[]> {
    return this.http.get<{data: ProductCategory[]}>(this.apiURL).pipe(map((response)=> {
      return response.data
    }));
  }

  getCategoryById(id: number): Observable<ProductCategory> {
    return this.http.get<ProductCategory>(`${this.apiURL}/${id}`);
  }

  createCategory(category: ProductCategory): Observable<ProductCategory> {
    return this.http.post<ProductCategory>(this.apiURL, category);
  }

  updateCategory(id: number, category: ProductCategory): Observable<ProductCategory> {
    return this.http.put<ProductCategory>(`${this.apiURL}/${id}`, category);
  }

  deleteCategory(id: number): Observable<void> {
    return this.http.delete<void>(`${this.apiURL}/${id}`);
  }
}
