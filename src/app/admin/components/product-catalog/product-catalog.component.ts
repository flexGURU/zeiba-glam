import { Component } from '@angular/core';
import { ConfirmationService, MessageService } from 'primeng/api';
import { Product } from '../../../core/interfaces/interfaces';
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
  products: Product[] = [];
  loading = false;
  showDialog = false;
  selectedProduct: Product | null = null;
  quickViewVisible: boolean = false;

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
    this.productService.getAllProducts().subscribe({
      next: (products) => {
        this.products = products;
        this.loading = false;
      },
      error: () => {
        this.loading = false;
        this.messageService.add({
          severity: 'error',
          summary: 'Error',
          detail: 'Failed to load products',
        });
      },
    });
  }

  showAddDialog() {
    this.selectedProduct = null;
    this.showDialog = true;
  }

  editProduct(product: Product) {
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
