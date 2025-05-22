import { CommonModule } from '@angular/common';
import { Component } from '@angular/core';
import { RouterLink } from '@angular/router';
import { CartService } from '../../core/services/cart.service';
import { DrawerModule } from 'primeng/drawer';


@Component({
  selector: 'app-header',
  imports: [CommonModule, RouterLink, DrawerModule],
  templateUrl: './header.component.html',
  styleUrl: './header.component.css',
})
export class HeaderComponent {
  categories: string[] = [
    'Abayas',
    'Dresses',
    'Pants',
    'Blouses',
    'Scarves',
    'Handbags',
    'Shoes',
  ];
  isMenuOpen = false;

  toggleMenu() {
    this.isMenuOpen = !this.isMenuOpen;
  }
  cartItemCount: number = 0;

  navbarContent = [
    {
      label: 'Shop',
      icon: 'pi pi-shopping-bag',
      routerLink: '/products',
    },

    {
      label: 'Socials',
      icon: 'pi pi-share-alt',
      routerLink: '/socials',
    },
    {
      label: 'Contact',
      icon: 'pi pi-envelope',
      routerLink: '/contact',
    },
    {
      label: 'About Us',
      icon: 'pi pi-info-circle',
      routerLink: '/about',
    },
  ];

  constructor(private cartService: CartService) {
    this.cartService.cartItems$.subscribe((items) => {
      this.cartItemCount = items.reduce(
        (count, item) => count + item.quantity,
        0
      );
    });
  }
}
