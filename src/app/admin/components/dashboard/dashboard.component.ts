import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ButtonModule } from 'primeng/button';
import { ProductCatalogComponent } from '../product-catalog/product-catalog.component';
import { MessageService } from 'primeng/api';
import { ProductStatsComponent } from '../product-stats/product-stats.component';

@Component({
  selector: 'app-dashboard',
  imports: [CommonModule, ProductCatalogComponent, ButtonModule, ProductStatsComponent],
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

 
}
