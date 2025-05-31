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
