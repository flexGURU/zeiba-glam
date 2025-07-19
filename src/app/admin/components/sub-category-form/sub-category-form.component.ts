import { Component, EventEmitter, inject, Input, Output } from '@angular/core';
import { FormGroup, FormBuilder, Validators, ReactiveFormsModule } from '@angular/forms';
import { MessageService } from 'primeng/api';
import { ProductCategory, ProductSubCategory } from '../../../core/interfaces/interfaces';
import { CategoryService } from '../../../core/services/category.service';
import { ButtonModule } from 'primeng/button';
import { CommonModule } from '@angular/common';
import { ToastModule } from 'primeng/toast';
import { DialogModule } from 'primeng/dialog';
import { SelectModule } from 'primeng/select';
import { SubCategoryService } from '../../../core/services/sub-category.service';

@Component({
  selector: 'app-sub-category-form',
  imports: [
    ButtonModule,
    SelectModule,
    CommonModule,
    ToastModule,
    DialogModule,
    ReactiveFormsModule,
  ],
  templateUrl: './sub-category-form.component.html',
  styleUrl: './sub-category-form.component.css',
})
export class SubCategoryFormComponent {
  subCategoryForm!: FormGroup;
  saving = false;
  showDialog = false;
  @Input() categoryIdPassed!: number

  @Input() visible = false;
  @Output() visibleChange = new EventEmitter<boolean>();

  @Input() subCategory: ProductSubCategory | null = null;
  @Output() save = new EventEmitter<ProductSubCategory>();
  categories: ProductCategory[] = [];

  private categoryService = inject(CategoryService);
  private subCategoryService = inject(SubCategoryService);

  constructor(
    private fb: FormBuilder,
    private messageService: MessageService
  ) {
    this.subCategoryForm = this.fb.group({
      name: ['', Validators.required],
      description: ['', Validators.required],
    });
  }

  ngOnInit() {
    this.loadCategories();
  }
  loadCategories() {
    this.categoryService.getAllCategories().subscribe((resp) => {
      this.categories = resp ?? [];
    });
  }

  get f() {
    return this.subCategoryForm.controls;
  }

  ngOnChanges() {
    if (this.subCategory) {
      this.subCategoryForm.patchValue({
        name: this.subCategory.name,
        description: this.subCategory.description,
      });
    } else {
      this.subCategoryForm.reset();
    }
  }
  onSubmit() {
    if (this.subCategoryForm.invalid) return;

    this.saving = true;
    const formData = this.subCategoryForm.value;
    console.log("formdata", formData);
    const createSubCategory: ProductSubCategory = {
      name: formData.name,
      description: formData.description,
      category_id: this.categoryIdPassed}
    

    if (this.subCategory && this.subCategory.id) {
      const updatedSubCategory: ProductSubCategory = {
        ...this.subCategory,
        ...formData,
      };

      this.subCategoryService.updateSubCategory(this.subCategory.id, updatedSubCategory).subscribe({
        next: () => {
          this.afterSave('Sub-Category updated successfully', updatedSubCategory);
        },
        error: () => this.handleError('Failed to update sub-category'),
      });
    } else {
      this.subCategoryService.createSubCategory(createSubCategory).subscribe({
        next: (createdSubCategory) => {
          this.afterSave('Sub-Category updated successfully', createdSubCategory);
        },
        error: () => this.handleError('Failed to add sub-category'),
      });
    }
  }

  private afterSave(message: string, subCategory: ProductSubCategory) {
    this.saving = false;
    this.messageService.add({ severity: 'success', summary: 'Success', detail: message });
    this.visible = false;
    this.visibleChange.emit(false);
    this.save.emit(subCategory);
    this.subCategoryForm.reset();
  }

  private handleError(detail: string) {
    this.saving = false;
    this.messageService.add({ severity: 'error', summary: 'Error', detail });
  }

  onHide() {
    this.visible = false;
    this.visibleChange.emit(false);
    this.subCategoryForm.reset();
  }
}
