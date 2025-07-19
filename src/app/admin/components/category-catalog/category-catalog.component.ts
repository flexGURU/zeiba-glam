import { Component, inject } from '@angular/core';
import { ConfirmationService, MessageService } from 'primeng/api';
import { ButtonModule } from 'primeng/button';
import { ConfirmDialogModule } from 'primeng/confirmdialog';
import { DialogModule } from 'primeng/dialog';
import { SortIcon, TableModule } from 'primeng/table';
import { TagModule } from 'primeng/tag';
import { CategoryFormComponent } from '../category-form/category-form.component';
import { CategoryService } from '../../../core/services/category.service';
import { ProductCategory, ProductSubCategory } from '../../../core/interfaces/interfaces';
import { CommonModule } from '@angular/common';
import { InputText } from 'primeng/inputtext';
import { SubCategoryFormComponent } from '../sub-category-form/sub-category-form.component';
import { Router } from '@angular/router';
import { ToastModule } from 'primeng/toast';

@Component({
  selector: 'app-category-catalog',
  imports: [
    ButtonModule,
    InputText,
    CommonModule,
    ConfirmDialogModule,
    CategoryFormComponent,
    TagModule,
    TableModule,
    ToastModule,
  ],
  templateUrl: './category-catalog.component.html',
  styleUrl: './category-catalog.component.css',
  providers: [ConfirmationService, MessageService],
})
export class CategoryCatalogComponent {
  categories: ProductCategory[] = [];
  loading: boolean = false;
  showDialog = false;
  selectedCategory: ProductCategory | null = null;
  searchTerm = '';
  currentPage = 1;
  private searchTimeout: any;

  private categoryService = inject(CategoryService);
  confirmationService = inject(ConfirmationService);
  private messageService = inject(MessageService);

  constructor(private router: Router) {}

  ngOnInit() {
    this.loadCategories();
  }

  onProductSave(category: ProductCategory) {
    this.loadCategories();
  }
  showAddDialog() {
    this.selectedCategory = null;
    this.showDialog = true;
  }

  loadCategories() {
    this.loading = true;
    this.categoryService.getAllCategories().subscribe((resp) => {
      this.categories = resp ?? [];
      this.loading = false;
    });
  }

  onRowSelected(event: any) {
    this.selectedCategory = event.data;
    this.router.navigate(['/admin/sub-category'], {
      queryParams: { categoryId: this.selectedCategory?.id },
    });
  }
  onRowUnselect(event: any) {
    this.selectedCategory = null;
  }

  onPageChange(event: any) {}
  onGlobalFilter(event: any) {}
  private debounceSearch() {
    clearTimeout(this.searchTimeout);
    this.searchTimeout = setTimeout(() => {
      this.loadCategories();
    }, 500); // 500ms delay
  }

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
            next: (response) => {
              console.log('response', response);

              this.loadCategories();
              this.messageService.add({
                severity: 'success',
                summary: 'Success',
                detail: 'Product deleted',
              });
            },
            error: (e) => {
              console.log('error', e);

              this.messageService.add({
                severity: 'error',
                summary: 'Faailed to delete category',
                detail: `${e.error.message}`,
              });
            },
          });
        }
      },
    });
  }

  getGlobalFilterValue(dt: any): string {
    const globalFilter = dt.filters['global'];
    if (!globalFilter) return '';
    if (Array.isArray(globalFilter)) return '';
    return globalFilter.value || '';
  }
}
