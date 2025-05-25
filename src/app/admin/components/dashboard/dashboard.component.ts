import { Component } from '@angular/core';
import { ProductService } from '../../../core/services/product.service';
import { CommonModule } from '@angular/common';
import { ButtonModule } from 'primeng/button';
import { ProductCatalogComponent } from '../product-catalog/product-catalog.component';
import { MessageService } from 'primeng/api';
import { LogoComponent } from '../../../core/components/logo/logo.component';

@Component({
  selector: 'app-dashboard',
  imports: [CommonModule, ProductCatalogComponent, ButtonModule, LogoComponent],
  templateUrl: './dashboard.component.html',
  styleUrl: './dashboard.component.css',
  providers: [MessageService],
})
export class DashboardComponent {
  stats = {
    totalProducts: 0,
    inStock: 0,
    lowStock: 0,
    outOfStock: 0,
  };

  constructor(private productService: ProductService) {}

  ngOnInit() {
    this.loadStats();
  }

  loadStats() {
    this.productService.getAllProducts().subscribe((products) => {
      this.stats.totalProducts = products.length;
      this.stats.inStock = products.filter((p) => (p.stock || 0) > 10).length;
      this.stats.lowStock = products.filter(
        (p) => (p.stock || 0) > 0 && (p.stock || 0) <= 10
      ).length;
      this.stats.outOfStock = products.filter((p) => (p.stock || 0) === 0).length;
    });
  }
}
