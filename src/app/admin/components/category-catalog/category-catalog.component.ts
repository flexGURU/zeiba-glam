import { Component, inject } from '@angular/core';
import { ConfirmationService, MessageService } from 'primeng/api';
import { ButtonModule } from 'primeng/button';
import { ConfirmDialogModule } from 'primeng/confirmdialog';
import { DialogModule } from 'primeng/dialog';
import { SortIcon, TableModule } from 'primeng/table';
import { TagModule } from 'primeng/tag';
import { CategoryFormComponent } from '../category-form/category-form.component';
import { CategoryService } from '../../../core/services/category.service';
import { ProductCategory } from '../../../core/interfaces/interfaces';
import { CommonModule } from '@angular/common';
import { InputText } from 'primeng/inputtext';

@Component({
  selector: 'app-category-catalog',
  imports: [ButtonModule, InputText, CommonModule, ConfirmDialogModule, CategoryFormComponent, TagModule, TableModule],
  templateUrl: './category-catalog.component.html',
  styleUrl: './category-catalog.component.css',
  providers: [ConfirmationService, MessageService],
})
export class CategoryCatalogComponent {
  categories: ProductCategory[] = [];
  loading: boolean = false;
  showDialog = false;
  selectedCategory: ProductCategory | null = null;

  private categoryService = inject(CategoryService);
  confirmationService = inject(ConfirmationService);
  private messageService = inject(MessageService);
  

  ngOnInit() {
    this.loadCategories();
  }

  loadCategories() {
    this.loading = true;
    this.categoryService.getAllCategories().subscribe((resp) => {
      this.categories = resp ?? [];
      this.loading = false;
      console.log('categories', this.categories);
    });
  }

  onProductSave(s: any) {}
  showAddDialog() {
    this.selectedCategory = null;
    this.showDialog = true;
  }

  onPageChange(event: any) {}
  onGlobalFilter(event: any) {}
  getStockSeverity(stock: any) {
    return '';
  }
  editCategory(category: ProductCategory) {
    this.showDialog = true;
    this.selectedCategory = category;
  }



    deleteCategory(category: ProductCategory) {
      this.confirmationService.confirm({
        message: 'Are you sure you want to delete this product?',
        header: 'Confirm',
        icon: 'pi pi-exclamation-triangle',
        accept: () => {
          if (category.id) {
            this.categoryService.deleteCategory(category.id).subscribe({
              next: () => {
                this.loadCategories();
                this.messageService.add({
                  severity: 'success',
                  summary: 'Success',
                  detail: 'Product deleted',
                });
              },
              error: () => {
                this.messageService.add({
                  severity: 'error',
                  summary: 'Error',
                  detail: 'Failed to delete product',
                });
              },
            });
          }
        },
      });
    }
}
