import { CommonModule } from '@angular/common';
import { Component, EventEmitter, Input, Output } from '@angular/core';
import { DialogModule } from 'primeng/dialog';
import { Product } from '../../../core/services/product.service';
import { InputNumberModule } from 'primeng/inputnumber';
import { FormsModule } from '@angular/forms';
import { RouterLink } from '@angular/router';

@Component({
  selector: 'app-product-preview',
  imports: [
    DialogModule,
    CommonModule,
    InputNumberModule,
    FormsModule,
    RouterLink,
  ],
  templateUrl: './product-preview.component.html',
  styleUrl: './product-preview.component.css',
})
export class ProductPreviewComponent {
  @Input() quickViewVisible: boolean = false;
  @Input() selectedProduct: Product | null = null;
  @Output() quickViewVisibleChange = new EventEmitter<boolean>();

  selectedColor: string = '';
  selectedSize: string = '';
  selectedQuantity: number = 1;
}
