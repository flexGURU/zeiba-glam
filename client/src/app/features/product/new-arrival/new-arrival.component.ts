import { Component } from '@angular/core';
import { RouterLink } from '@angular/router';
import { CarouselModule } from 'primeng/carousel';
import {
  Product,
  ProductService,
} from '../../../core/services/product.service';

@Component({
  selector: 'app-new-arrival',
  imports: [CarouselModule, RouterLink],
  templateUrl: './new-arrival.component.html',
  styleUrl: './new-arrival.component.css',
})
export class NewArrivalComponent {
  newArrivals: Product[] = [];
  constructor(private productService: ProductService) {}

  carouselResponsiveOptions = [
    {
      breakpoint: '1024px',
      numVisible: 3,
      numScroll: 1,
    },
    {
      breakpoint: '768px',
      numVisible: 2,
      numScroll: 1,
    },
    {
      breakpoint: '560px',
      numVisible: 1,
      numScroll: 1,
    },
  ];

  ngOnInit(): void {
    console.log('new arrivals');

    this.loadNewArrivals();
    console.log('new arrivals');
  }

  loadNewArrivals(): void {
    this.productService.getNewArrivals().subscribe({
      next: (products) => {
        this.newArrivals = products.slice(0, 8);
      },
      error: (error) => {
        console.error('Error loading new arrivals:', error);
      },
    });
  }
}
