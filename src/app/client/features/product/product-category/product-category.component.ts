import { Component, inject } from '@angular/core';
import { Router, RouterLink } from '@angular/router';
import { CommonModule } from '@angular/common';
import { CategoryService } from '../../../../core/services/category.service';
import { ProductCategory } from '../../../../core/interfaces/interfaces';
import { query } from 'express';

@Component({
  selector: 'app-product-category',
  imports: [CommonModule],
  templateUrl: './product-category.component.html',
  styleUrl: './product-category.component.css',
})
export class ProductCategoryComponent {
  categories: ProductCategory[] = [];

  constructor(private router: Router) {}

  private categoryService = inject(CategoryService);

  ngOnInit(): void {
    this.loadCategories();
  }

  loadCategories(): void {
    this.categoryService.getAllCategories().subscribe({
      next: (categories) => {
        this.categories = categories;
      },
      error: (error) => {
        console.error('Error loading categories:', error);
      },
    });
  }

  onCategoryClick(category: string): void {
    this.router.navigate(['/product-list'], {
      queryParams: { category: category },
    });
  }
}
