import { Component, EventEmitter, Input, Output } from '@angular/core';
import {
  FormGroup,
  FormBuilder,
  Validators,
  FormsModule,
  ReactiveFormsModule,
} from '@angular/forms';
import { MessageService } from 'primeng/api';
import { Product } from '../../../core/services/interfaces';
import { ProductService } from '../../../core/services/product.service';
import { FirebaseService } from '../../services/firebase.service';
import { ButtonModule } from 'primeng/button';
import { CheckboxModule } from 'primeng/checkbox';
import { FileUploadModule } from 'primeng/fileupload';
import { CommonModule } from '@angular/common';
import { InputNumberModule } from 'primeng/inputnumber';
import { SelectModule } from 'primeng/select';
import { DialogModule } from 'primeng/dialog';

@Component({
  selector: 'app-product-form',
  imports: [
    SelectModule,
    ButtonModule,
    CheckboxModule,
    FileUploadModule,
    CommonModule,
    InputNumberModule,
    DialogModule,
    FormsModule,
    ReactiveFormsModule,
  ],
  templateUrl: './product-form.component.html',
  styleUrl: './product-form.component.css',
})
export class ProductFormComponent {
  @Input() visible = false;
  @Input() product: Product | null = null;
  @Output() visibleChange = new EventEmitter<boolean>();
  @Output() save = new EventEmitter<Product>();

  productForm: FormGroup;
  categories = [
    { label: 'Electronics', value: 'electronics' },
    { label: 'Clothing', value: 'clothing' },
    { label: 'Home & Garden', value: 'home-garden' },
    { label: 'Sports', value: 'sports' },
    { label: 'Books', value: 'books' },
  ];

  imagePreview: string | null = null;
  selectedFile: File | null = null;
  saving = false;

  constructor(
    private fb: FormBuilder,
    private productService: ProductService,
    private firebaseService: FirebaseService,
    private messageService: MessageService
  ) {
    this.productForm = this.fb.group({
      name: ['', Validators.required],
      price: [0, [Validators.required, Validators.min(0)]],
      category: ['', Validators.required],
      description: [''],
      stock: [0],
      material: [''],
      featured: [false],
      new: [false],
      bestSeller: [false],
    });
  }

  ngOnChanges() {
    if (this.product) {
      this.productForm.patchValue(this.product);
      this.imagePreview = this.product.image;
    } else {
      this.productForm.reset();
      this.imagePreview = null;
      this.selectedFile = null;
    }
  }

  onImageSelect(event: any) {
    const file = event.files[0];
    if (file) {
      this.selectedFile = file;
      const reader = new FileReader();
      reader.onload = () => {
        this.imagePreview = reader.result as string;
      };
      reader.readAsDataURL(file);
    }
  }

  async onSubmit() {
    if (this.productForm.valid) {
      this.saving = true;
      try {
        const formData = this.productForm.value;

        // Upload image if new file selected
        if (this.selectedFile) {
          const imageUrl = await this.firebaseService.uploadImage(
            this.selectedFile,
            formData.name.replace(/\s+/g, '_').toLowerCase()
          );
          formData.image = imageUrl;
        } else if (this.product) {
          formData.image = this.product.image;
        }

        const productData: Product = {
          ...formData,
          id: this.product?.id || this.generateId(),
        };

        if (this.product) {
          this.productService.updateProduct(productData).subscribe({
            next: () => {
              this.messageService.add({
                severity: 'success',
                summary: 'Success',
                detail: 'Product updated',
              });
              this.save.emit(productData);
            },
            error: () => {
              this.messageService.add({
                severity: 'error',
                summary: 'Error',
                detail: 'Failed to update product',
              });
            },
          });
        } else {
          this.productService.addProduct(productData).subscribe({
            next: () => {
              this.messageService.add({
                severity: 'success',
                summary: 'Success',
                detail: 'Product added',
              });
              this.save.emit(productData);
            },
            error: () => {
              this.messageService.add({
                severity: 'error',
                summary: 'Error',
                detail: 'Failed to add product',
              });
            },
          });
        }
      } catch (error) {
        this.messageService.add({
          severity: 'error',
          summary: 'Error',
          detail: 'Failed to upload image',
        });
      }
      this.saving = false;
    }
  }

  onHide() {
    this.visible = false;
    this.visibleChange.emit(false);
  }

  private generateId(): string {
    return Date.now().toString(36) + Math.random().toString(36).substr(2);
  }
}
