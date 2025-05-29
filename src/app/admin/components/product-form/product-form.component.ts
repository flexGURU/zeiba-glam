import { Component, EventEmitter, Input, Output, signal, ViewChild } from '@angular/core';
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
import { FileUpload, FileUploadModule } from 'primeng/fileupload';
import { CommonModule } from '@angular/common';
import { InputNumberModule } from 'primeng/inputnumber';
import { SelectModule } from 'primeng/select';
import { DialogModule } from 'primeng/dialog';
import { InputTextModule } from 'primeng/inputtext';
import { TextareaModule } from 'primeng/textarea';
import { MultiSelectModule } from 'primeng/multiselect';

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
    InputTextModule,
    TextareaModule,
    MultiSelectModule,
  ],
  templateUrl: './product-form.component.html',
  styleUrl: './product-form.component.css',
})
export class ProductFormComponent {
  @Input() visible = false;
  @Input() product: Product | null = null;
  @Output() visibleChange = new EventEmitter<boolean>();
  @Output() save = new EventEmitter<Product>();

  // Add ViewChild to access the FileUpload component
  @ViewChild('fileUpload') fileUpload!: FileUpload;

  colorOptions = [
    { label: 'Red', value: 'red' },
    { label: 'Blue', value: 'blue' },
    { label: 'Green', value: 'green' },
    { label: 'Black', value: 'black' },
    { label: 'White', value: 'white' },
    { label: 'Yellow', value: 'yellow' },
    { label: 'Pink', value: 'pink' },
    { label: 'Purple', value: 'purple' },
    { label: 'Orange', value: 'orange' },
    { label: 'Gray', value: 'gray' },
  ];

  sizeOptions = [
    { label: 'Extra Small (XS)', value: 'XS' },
    { label: 'Small (S)', value: 'S' },
    { label: 'Medium (M)', value: 'M' },
    { label: 'Large (L)', value: 'L' },
    { label: 'Extra Large (XL)', value: 'XL' },
    { label: 'XXL', value: 'XXL' },
    { label: 'XXXL', value: 'XXXL' },
  ];

  // For hybrid approach - common sizes as buttons
  commonSizes = ['XS', 'S', 'M', 'L', 'XL', 'XXL'];

  productForm: FormGroup;
  categories = [
    { label: 'Electronics', value: 'electronics' },
    { label: 'Clothing', value: 'clothing' },
    { label: 'Home & Garden', value: 'home-garden' },
    { label: 'Sports', value: 'sports' },
    { label: 'Books', value: 'books' },
  ];

  imagePreview: string | null = null;
  saving = false;
  selectedFiles = signal<File[]>([]);
  imagePreviews = signal<string[]>([]);

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
      description: ['', Validators.required],
      stock: [0, [Validators.min(0)]],
      material: [''],
      colors: [[]],
      sizes: [[]],
      featured: [false],
      new: [false],
      bestSeller: [false],
    });
  }

  ngOnChanges() {
    if (this.product) {
      this.productForm.patchValue({
        name: this.product.name,
        price: this.product.price,
        category: this.product.category,
        description: this.product.description,
        stock: this.product.stock || 0,
        material: this.product.material || '',
        colors: this.product.colors || [],
        sizes: this.product.sizes || [],
        featured: this.product.featured || false,
        new: this.product.new || false,
        bestSeller: this.product.bestSeller || false,
      });
      this.imagePreview = this.product.image;
    } else {
      this.resetForm();
    }
  }
  onImageSelect(event: any) {
    if (event.files && event.files.length) {
      const newFiles: File[] = [...event.files]; // convert FileList to array

      this.selectedFiles.update((files) => [...files, ...newFiles]);

      for (let file of newFiles) {
        const reader = new FileReader();
        reader.onload = () => {
          this.imagePreviews.update((previews) => [...previews, reader.result as string]);
          console.log('Image Previews:', this.imagePreviews());
        };
        reader.readAsDataURL(file);
      }
    }
  }

  removeImage(index: number) {
    this.selectedFiles.update((files) => files.filter((_, i) => i !== index));
    this.imagePreviews.update((previews) => previews.filter((_, i) => i !== index));
  }

  async onSubmit() {
    if (this.productForm.valid) {
      this.saving = true;
      try {
        const formData = this.productForm.value;

        // Upload image if new file selected
        if (this.selectedFiles) {
          const imageUrl = await this.firebaseService.uploadImage(
            this.selectedFiles(),
            formData.name.replace(/\s+/g, '_').toLowerCase()
          );
          formData.image = imageUrl;
        } else if (this.product?.image) {
          formData.image = this.product.image;
        } else {
          // If no image is provided, set a default or handle accordingly
          formData.image = '';
        }

        const productData: Product = {
          id: this.product?.id || this.generateId(),
          name: formData.name,
          price: formData.price,
          category: formData.category,
          image: formData.image,
          description: formData.description,
          stock: formData.stock,
          material: formData.material || undefined,
          colors: formData.colors?.length > 0 ? formData.colors : undefined,
          sizes: formData.sizes?.length > 0 ? formData.sizes : undefined,
          featured: formData.featured || undefined,
          new: formData.new || undefined,
          bestSeller: formData.bestSeller || undefined,
        };

        if (this.product) {
          this.productService.updateProduct(productData).subscribe({
            next: () => {
              this.messageService.add({
                severity: 'success',
                summary: 'Success',
                detail: 'Product updated successfully',
              });
              this.save.emit(productData);
              this.resetForm();
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
                detail: 'Product added successfully',
              });
              this.save.emit(productData);
              this.resetForm();
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
    this.resetForm();
  }

  private resetForm() {
    this.productForm.reset({
      name: '',
      price: 0,
      category: '',
      description: '',
      stock: 0,
      material: '',
      colors: [],
      sizes: [],
      featured: false,
      new: false,
      bestSeller: false,
    });
    this.imagePreview = null;
    this.selectedFiles.set([]);

    // Clear the FileUpload component's internal state
    if (this.fileUpload) {
      this.fileUpload.clear();
    }
  }

  private generateId(): string {
    return Date.now().toString(36) + Math.random().toString(36).substr(2);
  }
}
