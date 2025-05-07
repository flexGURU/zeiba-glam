import { Component } from '@angular/core';
import { RouterLink } from '@angular/router';
import {
  ProductCategory,
  ProductService,
} from '../../../core/services/product.service';
import { CommonModule } from '@angular/common';

@Component({
  selector: 'app-product-category',
  imports: [RouterLink, CommonModule],
  templateUrl: './product-category.component.html',
  styleUrl: './product-category.component.css',
})
export class ProductCategoryComponent {
  categories: ProductCategory[] = [];

  constructor(private productService: ProductService) {}

  ngOnInit(): void {
    this.loadCategories();
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
