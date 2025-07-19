import { Injectable } from '@angular/core';
import { Observable, map } from 'rxjs';
import { environment } from '../../../environments/environment.development';
import { HttpClient } from '@angular/common/http';
import { ProductSubCategory } from '../interfaces/interfaces';

@Injectable({
  providedIn: 'root',
})
export class SubCategoryService {
  private readonly apiURL = environment.baseURL + '/sub-categories';

  constructor(private http: HttpClient) {}


  getAllSubCategories(): Observable<ProductSubCategory[]> {
    return this.http.get<{ data: ProductSubCategory[] }>(this.apiURL).pipe(
      map((response) => {
        return response.data;
      })
    );
  }

  getCategoryById(id: number): Observable<ProductSubCategory> {
    return this.http.get<ProductSubCategory>(`${this.apiURL}/${id}`);
  }

  createSubCategory(subCategory: ProductSubCategory): Observable<ProductSubCategory> {
    return this.http.post<ProductSubCategory>(this.apiURL, subCategory);
  }

  getSubCategoriesByCategoryId(categoryId: number): Observable<ProductSubCategory[]> {
    return this.http.get<{ data: ProductSubCategory[] }>(`${this.apiURL}/category/${categoryId}`).pipe(
      map((response) => {
        return response.data;
      })
    );
  }

  updateSubCategory(id: number, category: ProductSubCategory): Observable<ProductSubCategory> {
    return this.http.patch<ProductSubCategory>(`${this.apiURL}/${id}`, category);
  }

  deleteSubCategory(id: number): Observable<void> {
    return this.http.delete<void>(`${this.apiURL}/${id}`);
  }
}
