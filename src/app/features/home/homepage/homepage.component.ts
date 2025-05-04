import { Component } from '@angular/core';
import { ProductListComponent } from '../../product/product-list/product-list.component';
import { RouterLink } from '@angular/router';
import { CommonModule } from '@angular/common';
import { CarouselModule } from 'primeng/carousel';
import {
  Product,
  ProductCategory,
  ProductService,
} from '../../../core/services/product.service';
@Component({
  selector: 'app-homepage',
  imports: [ProductListComponent, CarouselModule, RouterLink, CommonModule],
  templateUrl: './homepage.component.html',
  styleUrl: './homepage.component.css',
})
export class HomepageComponent {
  newArrivals: Product[] = [];
  categories: ProductCategory[] = [];
  constructor(private productService: ProductService) {}

  ngOnInit() {
    this.loadNewArrivals();
    this.loadCategories();
  }

  carouselResponsiveOptions = [
    {
      breakpoint: '1024px',
      numVisible: 3,
      numScroll: 1,
    },
    {
      breakpoint: '768px',
      numVisible: 2,
      numScroll: 1,
    },
    {
      breakpoint: '560px',
      numVisible: 1,
      numScroll: 1,
    },
  ];

  loadNewArrivals(): void {
    this.productService.getNewArrivals().subscribe({
      next: (products) => {
        this.newArrivals = products.slice(0, 8);
      },
      error: (error) => {
        console.error('Error loading new arrivals:', error);
      },
    });
  }

  loadCategories(): void {
    this.productService.getProductCategories().subscribe({
      next: (categories) => {
        console.log(categories);

        this.categories = categories;
      },
      error: (error) => {
        console.error('Error loading categories:', error);
      },
    });
  }
}
