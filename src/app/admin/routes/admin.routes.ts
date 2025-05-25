import { Routes } from '@angular/router';
import { AdminLayoutComponent } from '../layout/admin-layout/admin-layout.component';
import { DashboardComponent } from '../components/dashboard/dashboard.component';
import { ProductCatalogComponent } from '../components/product-catalog/product-catalog.component';
import { LoginComponent } from '../components/login/login.component';

export const routes: Routes = [
  {
    path: '',
    component: AdminLayoutComponent,
    children: [
      { path: 'dashboard', component: DashboardComponent },
      { path: 'catalog', component: ProductCatalogComponent },
      { path: 'login', component: LoginComponent },
    ],
  },
];
