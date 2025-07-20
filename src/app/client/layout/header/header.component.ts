import { CommonModule } from '@angular/common';
import { Component, inject } from '@angular/core';
import { RouterLink } from '@angular/router';
import { CartService } from '../../services/cart.service';
import { DrawerModule } from 'primeng/drawer';
import { LogoComponent } from '../../../core/components/logo/logo.component';

@Component({
  selector: 'app-header',
  imports: [CommonModule, RouterLink, DrawerModule, LogoComponent],
  templateUrl: './header.component.html',
  styleUrl: './header.component.css',
})
export class HeaderComponent {
  categories: string[] = ['Abayas', 'Dresses', 'Pants', 'Blouses', 'Scarves', 'Handbags', 'Shoes'];
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

private cartService = inject(CartService);

getTotalCartItems(): number {
return this.cartService.getCartItems().reduce((total, item) => total + item.quantity, 0);
}

}
