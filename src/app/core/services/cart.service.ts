import { Injectable } from '@angular/core';
import { BehaviorSubject } from 'rxjs';
export interface CartItem {
  id: number;
  name: string;
  price: number;
  image: string;
  quantity: number;
  category: string;
}

@Injectable({
  providedIn: 'root',
})
export class CartService {
  private cartItems = new BehaviorSubject<CartItem[]>([]);
  cartItems$ = this.cartItems.asObservable();

  constructor() {
    // Load cart from localStorage if available
    const savedCart = localStorage.getItem('zeiba_cart');
    if (savedCart) {
      this.cartItems.next(JSON.parse(savedCart));
    }
  }

  private saveCart(): void {
    localStorage.setItem('zeiba_cart', JSON.stringify(this.cartItems.value));
  }

  addToCart(product: any, quantity: number = 1): void {
    const currentItems = this.cartItems.value;
    const existingItemIndex = currentItems.findIndex(
      (item) => item.id === product.id
    );

    if (existingItemIndex !== -1) {
      // Product already in cart, update quantity
      const updatedItems = [...currentItems];
      updatedItems[existingItemIndex].quantity += quantity;
      this.cartItems.next(updatedItems);
    } else {
      // Add new product to cart
      const newItem: CartItem = {
        id: product.id,
        name: product.name,
        price: product.price,
        image: product.image,
        quantity: quantity,
        category: product.category,
      };
      this.cartItems.next([...currentItems, newItem]);
    }

    this.saveCart();
  }

  updateQuantity(productId: number, quantity: number): void {
    if (quantity <= 0) {
      this.removeFromCart(productId);
      return;
    }

    const currentItems = this.cartItems.value;
    const updatedItems = currentItems.map((item) =>
      item.id === productId ? { ...item, quantity: quantity } : item
    );

    this.cartItems.next(updatedItems);
    this.saveCart();
  }

  removeFromCart(productId: number): void {
    const filteredItems = this.cartItems.value.filter(
      (item) => item.id !== productId
    );
    this.cartItems.next(filteredItems);
    this.saveCart();
  }

  clearCart(): void {
    this.cartItems.next([]);
    this.saveCart();
  }

  getCartTotal(): number {
    return this.cartItems.value.reduce(
      (total, item) => total + item.price * item.quantity,
      0
    );
  }
}
