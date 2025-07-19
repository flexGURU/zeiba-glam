import { Component, EventEmitter, Input, Output, signal, ViewChild } from '@angular/core';
import {
  FormGroup,
  FormBuilder,
  Validators,
  FormsModule,
  ReactiveFormsModule,
} from '@angular/forms';
import { MessageService } from 'primeng/api';
import {
  Product,
  ProductCategory,
  ProductSubCategory,
  RawProductPayload,
} from '../../../core/interfaces/interfaces';
import { ProductService } from '../../../core/services/product.service';
import { FirebaseService } from '../../services/firebase.service';
import { CategoryService } from '../../../core/services/category.service';
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
import { ToastModule } from 'primeng/toast';
import { log } from 'console';
import { SubCategoryService } from '../../../core/services/sub-category.service';

@Component({
  selector: 'app-product-form',
  imports: [
    ToastModule,
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
  @Input() product: RawProductPayload | null = null;
  @Output() visibleChange = new EventEmitter<boolean>();
  @Output() save = new EventEmitter<Product>();

  // Add ViewChild to access the FileUpload component
  @ViewChild('fileUpload') fileUpload!: FileUpload;
  categories: ProductCategory[] = [];

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

  saving = false;
  selectedFiles = signal<File[]>([]);
  imagePreviews = signal<string[]>([]);
  productSubCategories: ProductSubCategory[] = [];
  selectedCategory: string = '';

  constructor(
    private fb: FormBuilder,
    private productService: ProductService,
    private firebaseService: FirebaseService,
    private messageService: MessageService,
    private categoryService: CategoryService,
    private subCategoryService: SubCategoryService
  ) {
    this.productForm = this.fb.group({
      name: ['', Validators.required],
      price: [0, [Validators.required, Validators.min(0)]],
      category: ['', Validators.required],
      sub_category: ['', Validators.required],
      description: ['', Validators.required],
      stock: [0, [Validators.min(0)]],
      material: [''],
      colors: [[]],
      sizes: [[]],
      images: [this.imagePreviews()],
    });
  }

  ngOnInit() {
    this.loadCategories();
  }
  onCategoryModelChange(event: any) {
    this.subCategoryService.getSubCategoriesByCategoryId(event).subscribe((subCategories) => {
      this.productSubCategories = subCategories;
    });
  }

  ngOnChanges() {
    if (this.product) {
      this.imagePreviews.set(this.product.image_url);

      const categoryId =
        this.categories.find((cat) => cat.name === this.product?.category)?.id ?? null;

      const subCategoryId =
        this.productSubCategories.find((sub) => sub.name === this.product?.sub_category)?.id ??
        null;

      this.productForm.patchValue({
        name: this.product.name,
        price: this.product.price,
        category: categoryId,
        sub_category: subCategoryId,
        description: this.product.description,
        stock: this.product.stock_quantity || 0,
        material: this.product.material || '',
        colors: this.product.color || [],
        sizes: this.product.size || [],
      });
    } else {
      this.resetForm();
    }
  }
  loadCategories() {
    this.categoryService.getAllCategories().subscribe((resp) => {
      this.categories = resp ?? [];
    });
  }
  onImageSelect(event: any) {
    if (event.files && event.files.length) {
      const newFiles: File[] = [...event.files]; // convert FileList to array

      this.selectedFiles.update((files) => [...files, ...newFiles]);

      for (let file of newFiles) {
        const reader = new FileReader();
        reader.onload = () => {
          this.imagePreviews.update((previews) => [...previews, reader.result as string]);
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
    this.saving = true;
    if (!this.productForm.valid) return;
    const productData: Product = {
      name: this.productForm.value.name,
      price: this.productForm.value.price,
      category: this.productForm.value.category,
      sub_category: this.productForm.value.sub_category,
      image: [],
      description: this.productForm.value.description,
      colors: this.productForm.value.colors,
      sizes: this.productForm.value.sizes,
      stock: this.productForm.value.stock,
    };

    try {
      const imageUrls = await this.firebaseService.uploadImages(this.selectedFiles());
      productData.image = imageUrls;

      const saveFn = this.product
        ? this.productService.updateProduct(this.product?.id, productData)
        : this.productService.addProduct(productData);

      saveFn.subscribe({
        next: () => {
          this.saving = false;
          this.messageService.add({
            severity: 'success',
            summary: 'Success',
            detail: this.product ? 'Product updated successfully' : 'Product added successfully',
          });
          this.save.emit(productData);
          this.resetForm();
        },
        error: () => {
          this.messageService.add({
            severity: 'error',
            summary: 'Error',
            detail: this.product ? 'Failed to update product' : 'Failed to add product',
          });
        },
      });
    } catch (error) {
      console.error('Upload failed', error);
      this.messageService.add({
        severity: 'error',
        summary: 'Upload Error',
        detail: 'Image upload failed',
      });
    }

    this.saving = false;
  }

  onHide() {
    this.visible = false;
    this.visibleChange.emit(false);
    this.resetForm();
  }

  private resetForm() {
    this.productForm.reset({});
    this.imagePreviews.set([]);
    this.selectedFiles.set([]);

    if (this.fileUpload) {
      this.fileUpload.clear();
    }
  }
}
