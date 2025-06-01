import { HttpClient, HttpParams } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable, catchError, map, of, throwError } from 'rxjs';
import {
  PaginatedResponse,
  PaginationParams,
  Product,
  RawProductPayload,
} from '../interfaces/interfaces';
import { environment } from '../../../environments/environment.development';

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

  getPagintedProducts(params: PaginationParams): Observable<PaginatedResponse<RawProductPayload>> {
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

    return this.http
      .get<PaginatedResponse<RawProductPayload>>(`${this.apiURL}/products`, {
        params: httpParams,
      })
      .pipe(catchError(this.handleError));
  }

  getAllProducts(): Observable<Product[]> {
    return this.http.get<Product[]>(this.productsUrl).pipe(catchError(this.handleError));
  }

  getProductCategories(): Observable<ProductCategory[]> {
    const category = this.http.get<ProductCategory[]>(this.categoriesUrl);

    return category.pipe(catchError(this.handleError));
  }

  getProductById(id: number): Observable<Product> {
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

  updateProduct(
    productID: number | undefined,
    productData: Product
  ): Observable<RawProductPayload> {
    const prodPayload: RawProductPayload = {
      id: productData.id,
      name: productData.name,
      description: productData.description,
      price: productData.price,
      category: [productData.category],
      image_url: productData.image,
      size: productData.sizes,
      stock_quantity: productData.stock,
      color: productData.colors,
    };

    return this.http.patch<RawProductPayload>(`${this.apiURL}/products/${productID}`, prodPayload, {
      withCredentials: true,
    });
  }

  addProduct(productData: Product): Observable<RawProductPayload> {
    const prodPaylod: RawProductPayload = {
      name: productData.name,
      description: productData.description,
      price: productData.price,
      category: [productData.category],
      image_url: productData.image,
      size: productData.sizes,
      stock_quantity: productData.stock,
      color: productData.colors,
    };
    const product = this.http.post<RawProductPayload>(`${this.apiURL}/products`, prodPaylod, {
      withCredentials: true,
    });

    return product;
  }
  deleteProduct(productId: number): Observable<any> {
    const response = this.http.delete<any>(`${this.apiURL}/products/${productId}`);
    return response;
  }
}
