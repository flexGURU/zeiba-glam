import { Routes } from '@angular/router';
import { AdminLayoutComponent } from '../layout/admin-layout/admin-layout.component';
import { DashboardComponent } from '../components/dashboard/dashboard.component';
import { ProductCatalogComponent } from '../components/product-catalog/product-catalog.component';
import { LoginComponent } from '../components/login/login.component';
import { authGuard } from '../guard/auth.guard';
import { CategoryComponent } from '../components/category/category.component';
import { SubCategoryComponent } from '../components/sub-category/sub-category.component';

export const routes: Routes = [
  { path: 'login', component: LoginComponent },
  {
    path: '',
    component: AdminLayoutComponent,
    children: [
      { path: '', redirectTo: 'dashboard', pathMatch: 'full' },
      { path: 'dashboard', component: DashboardComponent, canActivate: [authGuard] },
      { path: 'catalog', component: ProductCatalogComponent, canActivate: [authGuard] },
      { path: 'category', component: CategoryComponent, canActivate: [authGuard] },
      { path: 'sub-category', component: SubCategoryComponent, canActivate: [authGuard] },
    ],
  },
];
