import { Component, inject } from '@angular/core';
import { CartItem, CartService } from '../../services/cart.service';
import { CommonModule } from '@angular/common';
import { ButtonModule } from 'primeng/button';
import { TagModule } from 'primeng/tag';
import { ChipModule } from 'primeng/chip';
import { DividerModule } from 'primeng/divider';
import { PanelModule } from 'primeng/panel';
import { CardModule } from 'primeng/card';
import { Router, RouterLink } from '@angular/router';
import { MessageService } from 'primeng/api';

@Component({
  selector: 'app-cart',
  templateUrl: './cart.component.html',
  imports: [
    CardModule,
    CommonModule,
    PanelModule,
    ButtonModule,
    TagModule,
    ChipModule,
    DividerModule,
    RouterLink,
  ],
  providers: [MessageService],
})
export class CartComponent {
  cartItems: CartItem[] = [];
  totalAmount: number = 0;
  private cartService = inject(CartService);
  private messageService = inject(MessageService);

  constructor(private router: Router) {}
  ngOnInit(): void {
    this.loadCartItems();
    this.calculateTotal();

    this.cartService.cartItems$.subscribe((items) => {
      this.cartItems = items;
    });
  }

  loadCartItems() {
    this.cartItems = this.cartService.getCartItems();
  }

  clearCart() {
    this.cartService.clearCart();
    this.loadCartItems();
  }

  calculateTotal(): void {
    this.totalAmount = this.cartItems.reduce((acc, item) => acc + item.total, 0);
  }

  proceedToCheckout(): void {
    if (this.cartItems.length === 0) {
      this.messageService.add({
        severity: 'warn',
        summary: 'Cart is empty!',
        detail: 'Add items before checking out.',
      });
      return;
    }
    this.router.navigate(['/checkout']);
  }
}
