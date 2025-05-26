import { Component } from '@angular/core';
import { Router } from '@angular/router';

@Component({
  selector: 'app-order-status',
  imports: [],
  templateUrl: './order-status.component.html',
  styleUrl: './order-status.component.css',
})
export class OrderStatusComponent {
  constructor(private router: Router) {}

  continueShopping() {
    this.router.navigate(['/']);
  }

  trackOrder() {
    this.router.navigate(['/orders']);
  }
}
