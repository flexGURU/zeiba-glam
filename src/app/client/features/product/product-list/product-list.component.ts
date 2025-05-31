import { CommonModule } from '@angular/common';
import { Component, effect, signal } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { SelectItem, SelectModule } from 'primeng/select';
import { SliderModule } from 'primeng/slider';
import { ProgressSpinnerModule } from 'primeng/progressspinner';
import { ActivatedRoute, RouterLink } from '@angular/router';
import { PaginatorModule } from 'primeng/paginator';
import { DialogModule } from 'primeng/dialog';
import { InputNumberModule } from 'primeng/inputnumber';
import { CartService } from '../../../../core/services/cart.service';
import { ProductPreviewComponent } from '../product-preview/product-preview.component';
import { Product } from '../../../../core/interfaces/interfaces';
import { ProductService } from '../../../../core/services/product.service';
import { ButtonModule } from 'primeng/button';

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
  products: Product[] = [];

  filteredProducts: Product[] = [];
  totalFilteredProducts: number = 0;

  // Pagination
  currentPage: number = 1;
  pageSize: number = 12;

  // View Mode
  viewMode: 'grid' | 'list' = 'grid';

  // Filters
  showFilters: boolean = false;
  priceRange: number[] = [0, 1000];
  availableColors: string[] = [
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
  selectedColors = signal<string[]>([]);
  availableSizes: string[] = ['XS', 'S', 'M', 'L', 'XL', '2XL'];
  selectedSizes = signal<string[]>([]);

  // Sort Options
  sortOptions: SelectItem[] = [];
  selectedSortOption: string = 'newest';

  // Quick View
  quickViewVisible: boolean = false;
  selectedProduct: Product | null = null;
  selectedColor: string = '';
  selectedSize: string = '';
  selectedQuantity: number = 1;

  constructor(
    private route: ActivatedRoute,
    private productService: ProductService,
    private cartService: CartService
  ) {}

  ngOnInit(): void {
    this.route.paramMap.subscribe((params) => {
      const category = params.get('category');
      this.category = category || '';
      this.updateCategoryTitle();
      this.loadProducts();
    });
  }

  updateCategoryTitle(): void {
    if (!this.category) {
      this.categoryTitle = 'All Products';
    } else if (this.category === 'new-arrivals') {
      this.categoryTitle = 'New Arrivals';
    } else if (this.category === 'best-sellers') {
      this.categoryTitle = 'Best Sellers';
    } else {
      this.categoryTitle = this.category.charAt(0).toUpperCase() + this.category.slice(1);
    }
  }

  getDescription(): string {
    if (this.category === 'new-arrivals') {
      return 'Check out our latest additions to stay on trend with the newest styles.';
    } else if (this.category === 'best-sellers') {
      return 'Discover our most popular pieces loved by our customers.';
    } else if (this.category === 'abayas') {
      return 'Elegant and modern abayas for every occasion.';
    } else if (this.category === 'dresses') {
      return 'Beautiful long dresses with contemporary designs.';
    } else if (this.category === 'pants') {
      return 'Comfortable and stylish pants to complete your outfit.';
    } else if (this.category === 'blouses') {
      return 'Versatile blouses perfect for any day of the week.';
    } else if (this.category === 'scarves') {
      return 'Add the perfect finishing touch with our scarves collection.';
    } else if (this.category === 'handbags') {
      return 'Stylish handbags to complement your ensemble.';
    } else if (this.category === 'shoes') {
      return 'Complete your look with our elegant footwear selection.';
    } else {
      return "Discover our complete collection of women's fashion.";
    }
  }

  loadProducts(): void {
    this.loading = true;

    this.productService.getAllProducts().subscribe(this.handleProductsResponse.bind(this));
  }

  handleProductsResponse(products: Product[]): void {
    this.products = products;
    this.applyFilters();
    this.loading = false;
  }

  applyFilters(): void {
    let result = [...this.products];

    // Apply price filter
    result = result.filter((p) => p.price >= this.priceRange[0] && p.price <= this.priceRange[1]);

    // Apply color filter if any selected
    if (this.selectedColors().length > 0) {
      result = result.filter((p) => {
        return p.colors?.some((color) => this.selectedColors().includes(color));
      });
    }

    // Apply size filter if any selected
    if (this.selectedSizes().length > 0) {
      result = result.filter((p) => {
        return p.sizes?.some((size) => this.selectedSizes().includes(size));
      });
    }
    console.log('result', result);

    // Apply sorting
    this.applySorting(result);

    this.totalFilteredProducts = result.length;

    // Apply pagination
    const startIndex = (this.currentPage - 1) * this.pageSize;
    this.filteredProducts = result.slice(startIndex, startIndex + this.pageSize);
  }

  applySorting(products: Product[]): void {
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
    const currentColors = this.selectedColors();

    if (!currentColors.includes(color)) {
      this.selectedColors.update((colors) => [...colors, color]);
    } else {
      this.selectedColors.update((colors) => colors.filter((c) => c !== color));
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
    this.selectedColors.set([]);
    this.selectedSizes.set([]);
    this.currentPage = 1;
    this.applyFilters();
  }

  hasActiveFilters(): boolean {
    return (
      this.priceRange[0] > 0 ||
      this.priceRange[1] < 1000 ||
      this.selectedColors.length > 0 ||
      this.selectedSizes.length > 0
    );
  }

  addToCart(product: Product, quantity: number = 1): void {
    this.cartService.addToCart(product, quantity);

    // Use PrimeNG Toast service here if you want to show a notification
    console.log(`Added ${quantity} of ${product.name} to cart`);

    this.quickViewVisible = false;
  }

  quickView(product: Product): void {
    this.selectedProduct = product;
    this.selectedQuantity = 1;

    // Set default selections if available
    if (product.colors && product.colors.length > 0) {
      this.selectedColor = product.colors[0];
    }

    if (product.sizes && product.sizes.length > 0) {
      this.selectedSize = product.sizes[0];
    }

    this.quickViewVisible = true;
  }
}
