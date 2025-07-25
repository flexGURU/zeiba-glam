import { CommonModule } from '@angular/common';
import { Component, EventEmitter, Input, Output } from '@angular/core';
import { DialogModule } from 'primeng/dialog';
import { InputNumberModule } from 'primeng/inputnumber';
import { FormsModule } from '@angular/forms';
import { RouterLink } from '@angular/router';
import { Product, RawProductPayload } from '../../../../core/interfaces/interfaces';
import { ButtonModule } from 'primeng/button';

@Component({
  selector: 'app-product-preview',
  imports: [DialogModule, CommonModule, InputNumberModule, FormsModule, RouterLink, ButtonModule],
  templateUrl: './product-preview.component.html',
  styleUrl: './product-preview.component.css',
})
export class ProductPreviewComponent {
  @Input() quickViewVisible: boolean = false;
  @Input() selectedProduct: RawProductPayload | null = null;
  @Output() quickViewVisibleChange = new EventEmitter<boolean>();

  selectedColor: string = '';
  selectedSize: string = '';
  selectedQuantity: number = 1;
}
