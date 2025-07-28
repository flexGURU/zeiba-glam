import { CommonModule } from '@angular/common';
import { Component, inject, Input } from '@angular/core';
import { ButtonModule } from 'primeng/button';
import { TableModule } from 'primeng/table';
import { SubCategoryFormComponent } from '../sub-category-form/sub-category-form.component';
import { ProductCategory, ProductSubCategory } from '../../../core/interfaces/interfaces';
import { ConfirmationService, MessageService } from 'primeng/api';
import { CategoryService } from '../../../core/services/category.service';
import { ConfirmDialog } from 'primeng/confirmdialog';
import { InputText } from 'primeng/inputtext';
import { SubCategoryService } from '../../../core/services/sub-category.service';
import { ActivatedRoute, Router } from '@angular/router';

@Component({
  selector: 'app-sub-category-catalog',
  imports: [
    ButtonModule,
    InputText,
    CommonModule,
    TableModule,
    SubCategoryFormComponent,
    ConfirmDialog,
  ],
  templateUrl: './sub-category-catalog.component.html',
  styleUrl: './sub-category-catalog.component.css',
  providers: [ConfirmationService, MessageService],
})
export class SubCategoryCatalogComponent {
  selectedSubCategory: ProductSubCategory | null = null;
  loading: boolean = false;
  showSubDialog = false;
  subCategories: ProductCategory[] = [];
  categories: ProductCategory[] = [];
  categoryId: number = 0;
  categorySubCategory: string = '';

  private subCategoryService = inject(SubCategoryService);
  private categoryService = inject(CategoryService);
  confirmationService = inject(ConfirmationService);
  private messageService = inject(MessageService);
  constructor(private router: ActivatedRoute) {}

  showAddSubCategoryDialog() {
    this.selectedSubCategory = null;
    this.showSubDialog = true;
  }

  ngOnInit() {
    this.loadCategories();
    this.router.queryParams.subscribe((params) => {
      if (params['categoryId']) {
        this.categoryId = +params['categoryId'];
        this.loadSubCategoriesByCategoryId();
        this.loadCategoryById();
      } else {
        this.loadSubCategories();
      }
    });
  }

  loadSubCategoriesByCategoryId() {
    if (this.categoryId) {
      this.loading = true;
      this.subCategoryService.getSubCategoriesByCategoryId(this.categoryId).subscribe((resp) => {
        this.subCategories = resp ?? [];
        this.loading = false;
      });
    }
  }

  loadSubCategories() {
    this.loading = true;
    this.subCategoryService.getAllSubCategories().subscribe((resp) => {
      this.subCategories = resp ?? [];
      this.loading = false;
    });
  }
  loadCategories() {
    this.loading = true;
    this.categoryService.getAllCategories().subscribe((resp) => {
      this.categories = resp ?? [];
      this.loading = false;
    });
  }
  loadCategoryById() {
    if (this.categoryId) {
      this.categoryService.getCategoryById(this.categoryId).subscribe((category) => {
        this.categorySubCategory = category.name;
      });
    }
  }
  onProductSave(category: ProductCategory) {
    this.loadSubCategoriesByCategoryId();
  }
  onPageChange(event: any) {}
  onGlobalFilter(event: any) {}
  editSubCategory(category: ProductSubCategory) {
    this.showSubDialog = true;
    this.selectedSubCategory = category;
    
  }

  deleteSubCategory(subCategory: ProductSubCategory) {
        this.confirmationService.confirm({
          message: 'Are you sure you want to delete this sub-category?',
          header: 'Confirm',
          icon: 'pi pi-exclamation-triangle',
          accept: () => {
            if (subCategory.id) {
              this.subCategoryService.deleteSubCategory(subCategory.id).subscribe({
                next: () => {
                  this.loadSubCategoriesByCategoryId();
                  this.messageService.add({
                    severity: 'success',
                    summary: 'Success',
                    detail: 'Sub-Category deleted',
                  });
                },
                error: (e) => {
                  this.messageService.add({
                    severity: 'error',
                    summary: 'Failed to delete sub-category',
                    detail: `${e.error.message}`,
                  });
                },
              });
            }
          },
        });
  }
}
