import { CommonModule } from '@angular/common';
import { Component, EventEmitter, inject, Input, Output } from '@angular/core';
import { FormBuilder, FormGroup, ReactiveFormsModule, Validators } from '@angular/forms';
import { ButtonModule } from 'primeng/button';
import { DialogModule } from 'primeng/dialog';
import { ProductCategory } from '../../../core/interfaces/interfaces';
import { CategoryService } from '../../../core/services/category.service';
import { MessageService } from 'primeng/api';
import { ToastModule } from 'primeng/toast';

@Component({
  selector: 'app-category-form',
  imports: [ReactiveFormsModule, ToastModule, DialogModule, ButtonModule, CommonModule],
  templateUrl: './category-form.component.html',
  styleUrl: './category-form.component.css',
  providers: [MessageService],
})
export class CategoryFormComponent {
  categoryForm!: FormGroup;
  saving = false;
  showDialog = false;

  @Input() visible = false;
  @Output() visibleChange = new EventEmitter<boolean>();

  @Input() category: ProductCategory | null = null;
  @Output() save = new EventEmitter<ProductCategory>();
  private categoryService = inject(CategoryService);

  constructor(
    private fb: FormBuilder,
    private messageService: MessageService
  ) {
    this.categoryForm = this.fb.group({
      name: ['', Validators.required],
      description: ['', Validators.required],
    });
  }

  get f() {
    return this.categoryForm.controls;
  }

  ngOnChanges() {
    if (this.category) {
      this.categoryForm.patchValue({
        name: this.category.name,
        description: this.category.description,
      });
    } else {
      this.categoryForm.reset();
    }
  }
  onSubmit() {
    if (this.categoryForm.invalid) return;

    this.saving = true;
    const formData = this.categoryForm.value;

    if (this.category && this.category.id) {
      const updatedCategory: ProductCategory = {
        ...this.category,
        ...formData,
      };

      this.categoryService.updateCategory(this.category.id, updatedCategory).subscribe({
        next: () => {
          this.afterSave('Category updated successfully', updatedCategory);
        },
        error: () => this.handleError('Failed to update category'),
      });
    } else {
      this.categoryService.createCategory(formData).subscribe({
        next: (createdCategory) => {
          this.afterSave('Category added successfully', createdCategory);
        },
        error: () => this.handleError('Failed to add category'),
      });
    }
  }

  private afterSave(message: string, category: ProductCategory) {
    this.saving = false;
    this.messageService.add({ severity: 'success', summary: 'Success', detail: message });
    this.visible = false;
    this.visibleChange.emit(false);
    this.save.emit(category);
    this.categoryForm.reset();
  }

  private handleError(detail: string) {
    this.saving = false;
    this.messageService.add({ severity: 'error', summary: 'Error', detail });
  }

  onHide() {
    this.visible = false;
    this.visibleChange.emit(false);
    this.categoryForm.reset();
  }
}
