import { Component } from '@angular/core';
import { ProductListComponent } from '../../product/product-list/product-list.component';
import { RouterLink } from '@angular/router';
import { CommonModule } from '@angular/common';
import { CarouselModule } from 'primeng/carousel';
import { NewArrivalComponent } from '../../product/new-arrival/new-arrival.component';
import { ProductCategoryComponent } from '../../product/product-category/product-category.component';
import { LandingComponent } from '../landing/landing.component';
@Component({
  selector: 'app-homepage',
  imports: [
    ProductListComponent,
    CarouselModule,
    CommonModule,

    ProductCategoryComponent,
    LandingComponent,
  ],
  templateUrl: './homepage.component.html',
  styleUrl: './homepage.component.css',
})
export class HomepageComponent {}
