import { Component, inject } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { ToastModule } from 'primeng/toast';
import { BreadcrumbModule } from 'primeng/breadcrumb';
import { CommonModule } from '@angular/common';
import { CartService } from '../../../../../core/services/cart.service';
import { MessageService } from 'primeng/api';
import { TooltipModule } from 'primeng/tooltip';
import { InputNumberModule } from 'primeng/inputnumber';
import { FormsModule } from '@angular/forms';
import { TabsModule } from 'primeng/tabs';
import { TabViewModule } from 'primeng/tabview';
import { ButtonModule } from 'primeng/button';
import { ProductRelatedComponent } from '../../product-related/product-related.component';
import { Product } from '../../../../../core/interfaces/interfaces';
import { ProductService } from '../../../../../core/services/product.service';

@Component({
  selector: 'app-product-detail',
  imports: [
    ToastModule,
    BreadcrumbModule,
    CommonModule,
    InputNumberModule,
    TooltipModule,
    FormsModule,
    TabsModule,
    TabViewModule,
    ButtonModule,
    ProductRelatedComponent,
  ],
  templateUrl: './product-detail.component.html',
  styleUrl: './product-detail.component.css',
  providers: [MessageService],
})
export class ProductDetailComponent {
  productId!: string;
  product!: Product;
  selectedColor: string = '';
  selectedSize: string = '';
  quantity: number = 1;
  relatedProducts: Product[] = [];

  productService = inject(ProductService);
  cartService = inject(CartService);
  messageService = inject(MessageService);

  constructor(
    private route: ActivatedRoute,
    private router: Router
  ) {}

  ngOnInit(): void {
    this.route.params.subscribe((params) => {
      this.productId = params['product-id'];
      console.log('product', typeof this.productId);
      this.getProductById(this.productId);
    });

    this.loadRelatedProducts();
  }

  getProductById = (id: string) => {
    this.productService.getProductById(id).subscribe((response) => {
      this.product = response;
      console.log('got', this.product);
    });
  };

  isColorSelected(color: string): boolean {
    return this.selectedColor === color;
  }

  getColorName(hexColor: string): string {
    const colorMap: { [key: string]: string } = {
      '#000000': 'Black',
      '#ffffff': 'White',
      '#6b7280': 'Gray',
      '#ef4444': 'Red',
      // Add more color mappings as needed
    };

    return colorMap[hexColor] || hexColor;
  }

  private loadRelatedProducts(): void {
    // In a real app, this would be fetched from a service
    this.relatedProducts = [];
  }

  buyNow(): void {
    this.addToCart();
    this.router.navigate(['/checkout']);
  }

  addToCart(): void {
    if (!this.selectedSize && this.product.sizes && this.product.sizes.length > 0) {
      this.messageService.add({
        severity: 'error',
        summary: 'Error',
        detail: 'Please select a size',
      });
      return;
    }

    if (!this.selectedColor && this.product.colors && this.product.colors.length > 0) {
      this.messageService.add({
        severity: 'error',
        summary: 'Error',
        detail: 'Please select a color',
      });
      return;
    }

    const item = {
      product: this.product,
      quantity: this.quantity,
      color: this.selectedColor,
      size: this.selectedSize,
      total: this.product.price * this.quantity,
    };

    this.cartService.addToCart(item);

    this.messageService.add({
      severity: 'success',
      summary: 'Success',
      detail: 'Product added to cart!',
    });
  }

  isSizeSelected(size: string): boolean {
    return this.selectedSize === size;
  }
}
