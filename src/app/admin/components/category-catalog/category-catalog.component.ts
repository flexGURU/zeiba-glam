import { Component } from '@angular/core';
import { ConfirmationService } from 'primeng/api';
import { ButtonModule } from 'primeng/button';
import { ConfirmDialogModule } from 'primeng/confirmdialog';
import { DialogModule } from 'primeng/dialog';
import { SortIcon, TableModule } from 'primeng/table';
import { TagModule } from 'primeng/tag';
import { CategoryFormComponent } from '../category-form/category-form.component';

@Component({
  selector: 'app-category-catalog',
  imports: [ButtonModule, ConfirmDialogModule, CategoryFormComponent, TagModule, TableModule],
  templateUrl: './category-catalog.component.html',
  styleUrl: './category-catalog.component.css',
  providers: [ConfirmationService],
})
export class CategoryCatalogComponent {
  categories: string[] = [];
  loading: boolean = false;
  showDialog = false;
  selectedCategory: string | null = null;
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
  editCategory(category: string) {}
  deleteCategory(category: string) {}
}
