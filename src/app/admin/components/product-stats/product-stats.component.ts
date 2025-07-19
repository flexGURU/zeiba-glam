import { Component, inject } from '@angular/core';
import { ProductService } from '../../../core/services/product.service';
import { DashboardService } from '../../services/dashboard.service';
import { DashboardStats } from '../../../core/interfaces/interfaces';

@Component({
  selector: 'app-product-stats',
  imports: [],
  templateUrl: './product-stats.component.html',
  styleUrl: './product-stats.component.css',
})
export class ProductStatsComponent {
  stats: DashboardStats | null = null;

  private dashboardService = inject(DashboardService);

  ngOnInit() {
    this.loadStats();
  }

  loadStats() {
    this.dashboardService.getDashboardData().subscribe({
      next: (data) => {
        this.stats = data;
      },
      error: (error) => {
        console.error('Error loading product stats:', error);
      },
    });
  }
}
