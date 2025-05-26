import { Routes } from '@angular/router';
import { HomepageComponent } from '../features/home/homepage/homepage.component';
import { ClientLayoutComponent } from '../layout/client-layout/client-layout.component';
import { ProductDetailComponent } from '../features/product/product-list/product-detail/product-detail.component';
import { CheckoutComponent } from '../features/payment/checkout/checkout.component';

export const routes: Routes = [
  {
    path: '',
    component: ClientLayoutComponent,
    children: [
      { path: '', component: HomepageComponent },
      { path: 'product-detail/:product-id', component: ProductDetailComponent },
      { path: 'checkout', component: CheckoutComponent },
    ],
  },
];
