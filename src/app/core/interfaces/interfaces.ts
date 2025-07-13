export interface Product {
  id?: number;
  name: string;
  price: number;
  category: string;
  image: string[];
  description: string;
  colors: string[];
  sizes: string[];
  stock: number;
  material?: string;
}

export interface RawProductPayload {
  id?: number;
  name: string;
  description: string;
  price: number;
  category: string[];
  image_url: string[];
  size: string[];
  color: string[];
  stock_quantity: number;
  material?: string;
}

export interface PaginatedResponse<T> {
  data: T[];
  totalRecords: number;
  currentPage: number;
  totalPages: number;
  pageSize: number;
}

export interface PaginationParams {
  page: number;
  limit: number;
  search?: string;
  sortBy?: string;
  sortOrder?: 'asc' | 'desc';
}

export interface ProductCategory {
  id?: number;
  name: string;
  description?: string;

}
