import { Component } from '@angular/core';
import { ConfirmationService, MessageService } from 'primeng/api';
import { PaginationParams, Product, RawProductPayload } from '../../../core/interfaces/interfaces';
import { ProductService } from '../../../core/services/product.service';
import { ButtonModule } from 'primeng/button';
import { TableModule } from 'primeng/table';
import { TagModule } from 'primeng/tag';
import { ProductFormComponent } from '../product-form/product-form.component';
import { ConfirmDialogModule } from 'primeng/confirmdialog';
import { InputText } from 'primeng/inputtext';

@Component({
  selector: 'app-product-catalog',
  imports: [
    ButtonModule,
    InputText,
    TableModule,
    TagModule,
    ProductFormComponent,
    ConfirmDialogModule,
  ],
  templateUrl: './product-catalog.component.html',
  styleUrl: './product-catalog.component.css',
  providers: [ConfirmationService, MessageService],
})
export class ProductCatalogComponent {
  products: RawProductPayload[] = [];
  loading = false;
  showDialog = false;
  selectedProduct: RawProductPayload | null = null;
  quickViewVisible: boolean = false;

  totalRecords = 0;
  currentPage = 1;
  pageSize = 10;

  searchTerm = '';
  sortField = '';
  sortOrder: 'asc' | 'desc' = 'asc';

  constructor(
    private productService: ProductService,
    private confirmationService: ConfirmationService,
    private messageService: MessageService
  ) {}

  ngOnInit() {
    this.loadProducts();
  }

  loadProducts() {
    this.loading = true;

    const params: PaginationParams = {
      page: this.currentPage,
      limit: this.pageSize,
      search: this.searchTerm || undefined,
      sortBy: this.sortField || undefined,
      sortOrder: this.sortOrder,
    };

    this.productService.getPagintedProducts(params).subscribe({
      next: (response) => {
        this.products = response.data;
        this.totalRecords = response.totalRecords;
        this.loading = false;
      },
      error: (error) => {
        this.loading = false;
        this.messageService.add({
          severity: 'error',
          summary: 'Error',
          detail: 'Failed to load products',
        });
      },
    });
  }

  // Handle pagination events from PrimeNG table
  onPageChange(event: any) {
    this.currentPage = event.page + 1; // PrimeNG uses 0-based indexing
    this.pageSize = event.rows;
    this.loadProducts();
  }

  // Handle sorting events
  onSort(event: any) {
    this.sortField = event.field;
    this.sortOrder = event.order === 1 ? 'asc' : 'desc';
    this.currentPage = 1; // Reset to first page when sorting
    this.loadProducts();
  }

  // Handle search/filter
  onGlobalFilter(event: any) {
    this.searchTerm = event.target.value;
    this.currentPage = 1; // Reset to first page when searching
    // Add debounce to avoid too many API calls
    this.debounceSearch();
  }

  private searchTimeout: any;
  private debounceSearch() {
    clearTimeout(this.searchTimeout);
    this.searchTimeout = setTimeout(() => {
      this.loadProducts();
    }, 500); // 500ms delay
  }

  showAddDialog() {
    this.selectedProduct = null;
    this.showDialog = true;
  }

  editProduct(product: RawProductPayload) {
    this.selectedProduct = { ...product };

    this.showDialog = true;
  }

  deleteProduct(product: Product) {
    this.confirmationService.confirm({
      message: 'Are you sure you want to delete this product?',
      header: 'Confirm',
      icon: 'pi pi-exclamation-triangle',
      accept: () => {
        if (product.id) {
          this.productService.deleteProduct(product.id).subscribe({
            next: () => {
              this.loadProducts();
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

  onProductSave(product: Product) {
    this.loadProducts();
    this.showDialog = false;
  }

  getStockSeverity(stock: number): string {
    if (stock === 0) return 'danger';
    if (stock < 10) return 'warning';
    return 'success';
  }
}
