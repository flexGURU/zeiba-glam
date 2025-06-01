import { Component } from '@angular/core';
import { ProductService } from '../../../core/services/product.service';

@Component({
  selector: 'app-product-stats',
  imports: [],
  templateUrl: './product-stats.component.html',
  styleUrl: './product-stats.component.css',
})
export class ProductStatsComponent {
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
