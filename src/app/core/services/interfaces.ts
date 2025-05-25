export interface Product {
  id: string;
  name: string;
  price: number;
  category: string;
  image: string;
  description: string;
  featured?: boolean;
  new?: boolean;
  bestSeller?: boolean;
  colors?: string[];
  sizes?: string[];
  stock?: number;
  material?: string;
}
