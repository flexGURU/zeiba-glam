import { Routes } from '@angular/router';

export const routes: Routes = [
  {
    path: '',
    loadChildren: () => import('./client/routes/client.routes').then((m) => m.routes),
  },
  {
    path: 'admin',
    loadChildren: () => import('./admin/routes/admin.routes').then((m) => m.routes),
  },

  { path: '', redirectTo: '/', pathMatch: 'full' },
];
