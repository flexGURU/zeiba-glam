import { CommonModule } from '@angular/common';
import { Component, EventEmitter, Input, Output } from '@angular/core';
import { FormBuilder, FormGroup, ReactiveFormsModule, Validators } from '@angular/forms';
import { ButtonModule } from 'primeng/button';
import { Category } from '../../../core/services/category.service';
import { DialogModule } from 'primeng/dialog';

@Component({
  selector: 'app-category-form',
  imports: [ReactiveFormsModule, DialogModule, ButtonModule, CommonModule],
  templateUrl: './category-form.component.html',
  styleUrl: './category-form.component.css',
})
export class CategoryFormComponent {
  categoryForm!: FormGroup;
  saving = false;

  @Input() visible = false;
  @Output() visibleChange = new EventEmitter<boolean>();

  @Input() category: string | null = null;
  @Output() save = new EventEmitter<Category>();

  constructor(private fb: FormBuilder) {
    this.categoryForm = this.fb.group({
      name: ['', Validators.required],
    });
  }

  get f() {
    return this.categoryForm.controls;
  }

  onSubmit() {
    if (this.categoryForm.valid) {
      const newCategory = this.categoryForm.value;
      console.log('Category submitted:', newCategory);
      // Here you'd send `newCategory` to your backend
      // e.g., this.categoryService.addCategory(newCategory).subscribe(...)
    }
  }

  onHide() {
    this.visible = false;
    this.visibleChange.emit(false);
    this.categoryForm.reset();
  }
}
