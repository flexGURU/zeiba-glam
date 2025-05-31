export interface Product {
  id?: string;
  name: string;
  price: number;
  category: string;
  image: string[];
  description: string;
  colors: string[];
  sizes?: string[];
  stock?: number;
  material?: string;
}

export interface RawProductPayload {
  name: string;
  description: string;
  price: number;
  category: string[];
  image_url: string[];
  size: string[];
  color: string[];
  stock_quantity: number;
}
