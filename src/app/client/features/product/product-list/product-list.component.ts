import { CommonModule } from '@angular/common';
import { Component, effect, inject, signal } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { SelectItem, SelectModule } from 'primeng/select';
import { SliderModule } from 'primeng/slider';
import { ProgressSpinnerModule } from 'primeng/progressspinner';
import { ActivatedRoute, RouterLink } from '@angular/router';
import { PaginatorModule } from 'primeng/paginator';
import { DialogModule } from 'primeng/dialog';
import { InputNumberModule } from 'primeng/inputnumber';
import { CartService } from '../../../services/cart.service';
import { ProductPreviewComponent } from '../product-preview/product-preview.component';
import { Product, RawProductPayload } from '../../../../core/interfaces/interfaces';
import { ProductService } from '../../../../core/services/product.service';
import { ButtonModule } from 'primeng/button';
import { CategoryService } from '../../../../core/services/category.service';

@Component({
  selector: 'app-product-list',
  imports: [
    SelectModule,
    CommonModule,
    SliderModule,
    ReactiveFormsModule,
    ProgressSpinnerModule,
    RouterLink,
    PaginatorModule,
    DialogModule,
    InputNumberModule,
    FormsModule,
    ProductPreviewComponent,
    ButtonModule,
  ],
  templateUrl: './product-list.component.html',
  styleUrl: './product-list.component.css',
})
export class ProductListComponent {
  loading: boolean = true;
  category: string = '';
  categoryTitle: string = 'All Products';
  products: RawProductPayload[] = [];

  filteredProducts: RawProductPayload[] = [];
  totalFilteredProducts: number = 0;

  // Pagination
  currentPage: number = 1;
  pageSize: number = 12;

  // View Mode
  viewMode: 'grid' | 'list' = 'grid';

  // Filters
  showFilters: boolean = false;
  priceRange: number[] = [0, 1000];
  availableColor: string[] = [
    'black',
    'white',
    'red',
    'blue',
    'green',
    'purple',
    'pink',
    'beige',
    'brown',
  ];
  selectedColor = signal<string[]>([]);
  availableSizes: string[] = ['XS', 'S', 'M', 'L', 'XL', '2XL'];
  selectedSizes = signal<string[]>([]);

  // Sort Options
  sortOptions: SelectItem[] = [];
  selectedSortOption: string = 'newest';

  // Quick View
  quickViewVisible: boolean = false;
  selectedProduct: RawProductPayload | null = null;
  selectedSize: string = '';
  selectedQuantity: number = 1;

  private productService = inject(ProductService);
  private categoryService = inject(CategoryService);

  constructor(
    private route: ActivatedRoute,
    private cartService: CartService
  ) {}

  ngOnInit(): void {
    this.route.queryParams.subscribe((params) => {
      if (params['category']) {
        this.category = params['category'];

        this.loadProductsByCategory(this.category);
      } else {
        this.loadProducts();
      }
    });
  }

  loadProducts(): void {
    this.loading = true;

    this.productService.getAllProducts().subscribe({
      next: (productResponse) => {
        this.products = productResponse;

        this.handleProductsResponse(this.products);

        this.loading = false;
      },
      error: (error) => {
        console.error('Error loading products:', error);
        this.loading = false;
      },
    });
  }

  loadProductsByCategory(category: string): void {
    this.loading = true;
    this.productService.getProductsByCategory([category]).subscribe({
      next: (productResponse) => {
        this.products = productResponse;

        this.handleProductsResponse(this.products);

        this.loading = false;
      },
      error: (error) => {
        console.error('Error loading products by category:', error);
        this.loading = false;
      },
    });
  }

  handleProductsResponse(products: RawProductPayload[]): void {
    this.products = products;
    this.applyFilters();
    this.loading = false;
  }

  applyFilters(): void {
    let result = [...this.products];

    // Apply price filter
    result = result.filter((p) => p.price >= this.priceRange[0] && p.price <= this.priceRange[1]);

    // Apply color filter if any selected
    if (this.selectedColor().length > 0) {
      result = result.filter((p) => {
        return p.color?.some((color) => this.selectedColor().includes(color));
      });
    }

    // Apply size filter if any selected
    if (this.selectedSizes().length > 0) {
      result = result.filter((p) => {
        return p.size?.some((size) => this.selectedSizes().includes(size));
      });
    }

    // Apply sorting
    this.applySorting(result);

    this.totalFilteredProducts = result.length;

    // Apply pagination
    const startIndex = (this.currentPage - 1) * this.pageSize;
    this.filteredProducts = result.slice(startIndex, startIndex + this.pageSize);
  }

  applySorting(products: RawProductPayload[]): void {
    switch (this.selectedSortOption) {
      case 'price_asc':
        products.sort((a, b) => a.price - b.price);
        break;
      case 'price_desc':
        products.sort((a, b) => b.price - a.price);
        break;
      case 'name_asc':
        products.sort((a, b) => a.name.localeCompare(b.name));
        break;
      case 'name_desc':
        products.sort((a, b) => b.name.localeCompare(a.name));
        break;
      // For 'newest' we assume the original array is already sorted by date
      default:
        break;
    }
  }

  onPageChange(event: any): void {
    this.currentPage = event.page + 1;
    this.pageSize = event.rows;
    this.applyFilters();
  }
  toggleColorFilter(color: string): void {
    const currentColor = this.selectedColor();

    if (!currentColor.includes(color)) {
      this.selectedColor.update((colors) => [...colors, color]);
    } else {
      this.selectedColor.update((colors) => colors.filter((c) => c !== color));
    }
    this.currentPage = 1;
    this.applyFilters();
  }

  toggleSizeFilter(size: string): void {
    const currentSizes = this.selectedSizes();

    if (!currentSizes.includes(size)) {
      this.selectedSizes.update((selectedSizes) => [...selectedSizes, size]);
    } else {
      this.selectedSizes.update((sizes) => sizes.filter((s) => s !== size));
    }
    this.currentPage = 1;
    this.applyFilters();
  }

  resetPriceFilter(): void {
    this.priceRange = [0, 1000];
    this.currentPage = 1;
    this.applyFilters();
  }

  clearAllFilters(): void {
    this.priceRange = [0, 1000];
    this.selectedColor.set([]);
    this.selectedSizes.set([]);
    this.currentPage = 1;
    this.applyFilters();
  }

  hasActiveFilters(): boolean {
    return (
      this.priceRange[0] > 0 ||
      this.priceRange[1] < 1000 ||
      this.selectedColor.length > 0 ||
      this.selectedSizes.length > 0
    );
  }

  addToCart(product: Product, quantity: number = 1): void {
    // this.cartService.addToCart(product, quantity);

    // Use PrimeNG Toast service here if you want to show a notification
    console.log(`Added ${quantity} of ${product.name} to cart`);

    this.quickViewVisible = false;
  }

  quickView(product: RawProductPayload): void {
    this.selectedProduct = product;
    this.selectedQuantity = 1;

    // Set default selections if available
    if (product.color && product.color.length > 0) {
      this.selectedColor.set([product.color[0]]);
    }

    if (product.size && product.size.length > 0) {
      this.selectedSize = product.size[0];
    }

    this.quickViewVisible = true;
  }
}
