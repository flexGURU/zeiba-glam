import { Injectable } from '@angular/core';
import { BehaviorSubject } from 'rxjs';
import { RawProductPayload } from '../../core/interfaces/interfaces';

export interface CartItem {
  product: RawProductPayload;
  quantity: number;
  color: string;
  size: string;
  total: number;
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

  addToCart(item: CartItem): void {
    const currentItems = this.cartItems.value;

    // Check if item already exists (same product, color, size)
    const existingIndex = currentItems.findIndex(
      (cartItem) =>
        cartItem.product.id === item.product.id &&
        cartItem.color === item.color &&
        cartItem.size === item.size
    );

    if (existingIndex > -1) {
      // Update quantity and total if item exists
      currentItems[existingIndex].quantity += item.quantity;
      currentItems[existingIndex].total = currentItems[existingIndex].quantity * item.product.price;
    } else {
      // Add new item if it doesn't exist
      currentItems.push(item);
    }

    // Update the BehaviorSubject and save to localStorage
    this.cartItems.next(currentItems);
    this.saveCart();
  }

  // removeFromCart(productId: string, color: string, size: string): void {
  //   const currentItems = this.cartItems.value.filter(
  //     (item) => !(item.product.id === productId && item.color === color && item.size === size)
  //   );
  //   this.cartItems.next(currentItems);
  //   this.saveCart();
  // }

  // updateQuantity(productId: string, color: string, size: string, quantity: number): void {
  //   const currentItems = this.cartItems.value;
  //   const itemIndex = currentItems.findIndex(
  //     (item) => item.product.id === productId && item.color === color && item.size === size
  //   );

  //   if (itemIndex > -1) {
  //     currentItems[itemIndex].quantity = quantity;
  //     currentItems[itemIndex].total = quantity * currentItems[itemIndex].product.price;
  //     this.cartItems.next(currentItems);
  //     this.saveCart();
  //   }
  // }

  getCartItems(): CartItem[] {
    return this.cartItems.value;
  }

  getCartTotal(): number {
    return this.cartItems.value.reduce((total, item) => total + item.total, 0);
  }

  getCartCount(): number {
    return this.cartItems.value.reduce((count, item) => count + item.quantity, 0);
  }

  clearCart(): void {
    this.cartItems.next([]);
    this.saveCart();
  }
}
