import { CommonModule } from '@angular/common';
import { Component } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { ButtonModule } from 'primeng/button';
import { CategoryCatalogComponent } from '../category-catalog/category-catalog.component';

@Component({
  selector: 'app-category',
  imports: [ReactiveFormsModule, CategoryCatalogComponent, CommonModule, FormsModule, ButtonModule],
  templateUrl: './category.component.html',
  styleUrl: './category.component.css',
})
export class CategoryComponent {}
