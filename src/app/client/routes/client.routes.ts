import { Routes } from '@angular/router';
import { HomepageComponent } from '../features/home/homepage/homepage.component';
import { ClientLayoutComponent } from '../layout/client-layout/client-layout.component';
import { ProductDetailComponent } from '../features/product/product-list/product-detail/product-detail.component';
import { CheckoutComponent } from '../features/payment/checkout/checkout.component';
import { ProductListComponent } from '../features/product/product-list/product-list.component';

export const routes: Routes = [
  {
    path: '',
    component: ClientLayoutComponent,
    children: [
      { path: '', component: HomepageComponent },
      { path: 'product-list', component: ProductListComponent },
      { path: 'product-detail/:product-id', component: ProductDetailComponent },
      { path: 'checkout', component: CheckoutComponent },
      {
        path: 'cart',
        loadComponent: () => import('../features/cart/cart.component').then((m) => m.CartComponent),
      },
    ],
  },
];
