import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable, catchError, map, of, throwError } from 'rxjs';
import { Product } from '../interfaces/interfaces';
import { environment } from '../../../environments/environment.development';
import { response } from 'express';

export interface ProductCategory {
  name: string;
  image: string;
}

@Injectable({
  providedIn: 'root',
})
export class ProductService {
  private productsUrl = 'products.json';
  private categoriesUrl = 'assets/images/categories/data/categories.json';
  private readonly apiURL = environment.baseURL;

  constructor(private http: HttpClient) {}
  getAllProducts(): Observable<Product[]> {
    return this.http.get<Product[]>(this.productsUrl).pipe(catchError(this.handleError));
  }

  getProductCategories(): Observable<ProductCategory[]> {
    const category = this.http.get<ProductCategory[]>(this.categoriesUrl);

    return category.pipe(catchError(this.handleError));
  }

  getProductById(id: string): Observable<Product> {
    return this.getAllProducts().pipe(
      map((products) => {
        const product = products.find((p) => p.id === id);
        if (!product) {
          throw new Error(`Product with id ${id} not found`);
        }
        return product;
      }),
      catchError(this.handleError)
    );
  }

  getProductsByCategory(category: string): Observable<Product[]> {
    return this.getAllProducts().pipe(
      map((products) =>
        products.filter((product) => product.category.toLowerCase() === category.toLowerCase())
      ),
      catchError(this.handleError)
    );
  }

  private handleError(error: any) {
    console.error('An error occurred:', error);
    return throwError(() => new Error('Something went wrong. Please try again later.'));
  }

  updateProduct(product: Product): Observable<Product | null> {
    return of(null);
  }
  addProduct(productData: Product): Observable<Product | null> {
    const product = this.http.post<Product>(`${this.apiURL}/products`, productData);

    product.pipe(
      map((response) => {
        console.log('data from backend', response);
      })
    );

    return product;
  }
  deleteProduct(productId: string): Observable<Product | null> {
    return of(null);
  }
}
